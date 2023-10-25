package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"salami/backend/prompts/terraform/openai_gpt4"
	backendTypes "salami/backend/types"
	"salami/common/change_set"
	"salami/common/constants"
	"salami/common/logger"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
)

type TerraformValidateOutput struct {
	Diagnostics []TerraformValidationError `json:"diagnostics"`
}

type TerraformValidationError struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Range    *struct {
		Start struct {
			Line int `json:"line"`
		} `json:"start"`
	} `json:"range"`
	Snippet struct {
		Code string `json:"code"`
	} `json:"snippet"`
}

func (t *Terraform) ValidateCode(
	newObjects []*commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
	changeSetRepository *change_set.ChangeSetRepository,
	llm backendTypes.Llm,
	retryCount int,
) error {
	if retryCount > constants.MaxFixValidationErrorRetries {
		return nil
	}
	if len(changeSetRepository.Diffs) == 0 {
		return nil
	}
	validationResults, err := generateValidationResults(newObjects)
	if err != nil {
		return err
	}
	if len(validationResults) == 0 {
		return nil
	}

	if err := processValidationResults(validationResults, symbolTable, changeSetRepository, llm); err != nil {
		return err
	}

	if err := t.ValidateCode(newObjects, symbolTable, changeSetRepository, llm, retryCount+1); err != nil {
		return err
	}

	return nil
}

func generateValidationResults(
	newObjects []*commonTypes.Object,
) ([]*backendTypes.CodeValidationResult, error) {
	var targetCodeBlocks []string
	var endLineNumbers []int

	for _, obj := range newObjects {
		numLines := strings.Count(obj.TargetCode, "\n")
		targetCodeBlocks = append(targetCodeBlocks, obj.TargetCode)

		if len(endLineNumbers) == 0 {
			endLineNumbers = append(endLineNumbers, numLines+1)
		} else {
			lastLineNumber := endLineNumbers[len(endLineNumbers)-1]
			endLineNumbers = append(endLineNumbers, lastLineNumber+2+numLines)
		}
	}

	targetCode := strings.Join(targetCodeBlocks, "\n\n")
	terraformValidateOutput, err := runTerraformValidate(targetCode)
	if err != nil {
		return nil, err
	}

	var validationResults []*backendTypes.CodeValidationResult
	if err := parseTerraformValidateOutput(terraformValidateOutput, newObjects, endLineNumbers, &validationResults); err != nil {
		return nil, err
	}

	populateValidationResultsReferencedObjects(newObjects, &validationResults)

	return validationResults, nil
}

func processValidationResults(
	validationResults []*backendTypes.CodeValidationResult,
	symbolTable *symbol_table.SymbolTable,
	changeSetRepository *change_set.ChangeSetRepository,
	llm backendTypes.Llm,
) error {
	var g errgroup.Group
	semaphoreChannel := make(chan struct{}, llm.GetMaxConcurrentExecutions())

	for _, validationResult := range validationResults {
		validationResult := validationResult

		g.Go(func() error {
			diff := changeSetRepository.GetDiffForObject(validationResult.ValidatedObject)
			if diff == nil {
				return nil
			}

			semaphoreChannel <- struct{}{}
			defer func() { <-semaphoreChannel }()

			messages, err := getValidateCodeLlmMessages(symbolTable, diff, validationResult, llm)
			if err != nil {
				return err
			}

			completion, err := llm.CreateCompletion(messages)
			if err != nil {
				return err
			}
			diff.NewObject.SetTargetCode(completion)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func runTerraformValidate(targetCode string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "terraform.*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "main.tf")
	err = os.WriteFile(tempFile, []byte(targetCode), 0644)
	if err != nil {
		return nil, err
	}

	initCmd := exec.Command("terraform", "init", "-backend=false")
	initCmd.Dir = tempDir
	if err := initCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run terraform init: %w", err)
	}

	validateCmd := exec.Command("terraform", "validate", "-json")
	validateCmd.Dir = tempDir
	output, _ := validateCmd.CombinedOutput()

	return output, nil
}

func parseTerraformValidateOutput(
	terraformValidateOutput []byte,
	newObjects []*commonTypes.Object,
	endLineNumbers []int,
	validationResults *[]*backendTypes.CodeValidationResult,
) error {
	var tfValidateOutput TerraformValidateOutput
	err := json.Unmarshal(terraformValidateOutput, &tfValidateOutput)
	if err != nil {
		return err
	}

	var tfErrors []TerraformValidationError
	for _, diagnostic := range tfValidateOutput.Diagnostics {
		if diagnostic.Severity == "error" && diagnostic.Range != nil {
			tfErrors = append(tfErrors, diagnostic)
		}
	}
	sort.Slice(tfErrors, func(i, j int) bool {
		return tfErrors[i].Range.Start.Line < tfErrors[j].Range.Start.Line
	})

	logger.Debug("tfErrors: " + fmt.Sprint(tfErrors))

	j := 0
	for _, tfError := range tfErrors {
		for tfError.Range.Start.Line > endLineNumbers[j] {
			j++
		}
		validatedObject := newObjects[j]
		errorMessage := fmt.Sprintf(
			"Summary: %s\nDetail: %s\nCode line: %s",
			tfError.Summary,
			tfError.Detail,
			tfError.Snippet.Code,
		)

		if len(*validationResults) > 0 && (*validationResults)[len(*validationResults)-1].ValidatedObject == validatedObject {
			errorMessages := (*validationResults)[len(*validationResults)-1].ErrorMessages
			(*validationResults)[len(*validationResults)-1].ErrorMessages = append(errorMessages, errorMessage)
		} else {
			*validationResults = append(*validationResults, &backendTypes.CodeValidationResult{
				ValidatedObject: validatedObject,
				ErrorMessages:   []string{errorMessage},
			})
		}
	}

	return nil
}

func populateValidationResultsReferencedObjects(
	newObjects []*commonTypes.Object,
	validationResults *[]*backendTypes.CodeValidationResult,
) {
	resourceMap := make(map[commonTypes.LogicalName]*commonTypes.Object)
	variableMap := make(map[string]*commonTypes.Object)

	for _, object := range newObjects {
		if object.IsResource() {
			resourceMap[object.ParsedResource.LogicalName] = object
		} else if object.IsVariable() {
			variableMap[object.ParsedVariable.Name] = object
		}
	}

	for _, validationResult := range *validationResults {
		if !validationResult.ValidatedObject.IsResource() {
			continue
		}
		for _, referencedLogicalName := range validationResult.ValidatedObject.ParsedResource.ReferencedResources {
			if obj, ok := resourceMap[referencedLogicalName]; ok {
				validationResult.ReferencedObjects = append(validationResult.ReferencedObjects, obj)
			}
		}
		for _, referencedVariableName := range validationResult.ValidatedObject.ParsedResource.ReferencedVariables {
			if obj, ok := variableMap[referencedVariableName]; ok {
				validationResult.ReferencedObjects = append(validationResult.ReferencedObjects, obj)
			}
		}
	}
}

func getValidateCodeLlmMessages(
	symbolTable *symbol_table.SymbolTable,
	diff *commonTypes.ChangeSetDiff,
	validationResult *backendTypes.CodeValidationResult,
	llm backendTypes.Llm,
) ([]interface{}, error) {
	var messages []interface{}

	switch llm.GetSlug() {
	case commonTypes.LlmOpenaiGpt4:
		llmMessages, err := openai_gpt4.GetValidateCodeMessages(symbolTable, diff, validationResult)
		if err != nil {
			return nil, err
		}

		for _, v := range llmMessages {
			messages = append(messages, v)
		}
	}
	return messages, nil
}
