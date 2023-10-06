package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"salami/backend/prompts/terraform/openai_gpt4"
	backendTypes "salami/backend/types"
	"salami/common/change_set"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"strings"

	"golang.org/x/sync/errgroup"
)

type TerraformValidationError struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Range    struct {
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
) []error {
	validationResults, err := generateValidationResults(newObjects)
	if err != nil {
		return []error{err}
	}
	if len(validationResults) == 0 {
		return nil
	}

	var g errgroup.Group
	semaphoreChannel := make(chan struct{}, llm.GetMaxConcurrentExecutions())

	for _, validationResult := range validationResults {
		validationResult := validationResult

		g.Go(func() error {
			semaphoreChannel <- struct{}{}
			defer func() { <-semaphoreChannel }()

			diff := changeSetRepository.GetDiffForObject(validationResult.ValidatedObject)
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
		return []error{err}
	}

	return nil
}

func generateValidationResults(
	newObjects []*commonTypes.Object,
) ([]*backendTypes.CodeValidationResult, error) {
	var targetCodeLines []string
	var endLineNumbers []int

	for _, obj := range newObjects {
		lines := strings.Split(obj.TargetCode, "\n")
		targetCodeLines = append(targetCodeLines, lines...)

		if len(endLineNumbers) == 0 {
			endLineNumbers = append(endLineNumbers, len(lines))
		} else {
			lastLineNumber := endLineNumbers[len(endLineNumbers)-1]
			endLineNumbers = append(endLineNumbers, lastLineNumber+2+len(lines))
		}
	}

	targetCode := strings.Join(targetCodeLines, "\n\n")
	terraformValidateOutput, err := runTerraformValidate(targetCode)
	if err != nil {
		return nil, err
	}

	var validationResults []*backendTypes.CodeValidationResult
	parseTerraformValidateOutput(terraformValidateOutput, newObjects, endLineNumbers, &validationResults)

	populateValidationResultsReferencedObjects(newObjects, &validationResults)

	return validationResults, nil
}

func runTerraformValidate(targetCode string) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "terraform.*.tf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())
	os.WriteFile(tempFile.Name(), []byte(targetCode), 0644)

	cmd := exec.Command("terraform", "validate", "-json", tempFile.Name())
	output, _ := cmd.CombinedOutput()

	return output, nil
}

func parseTerraformValidateOutput(
	terraformValidateOutput []byte,
	newObjects []*commonTypes.Object,
	endLineNumbers []int,
	validationResults *[]*backendTypes.CodeValidationResult,
) {
	var tfErrors []TerraformValidationError
	json.Unmarshal(terraformValidateOutput, &tfErrors)

	for _, tfError := range tfErrors {
		if tfError.Severity == "error" {
			var validatedObject *commonTypes.Object
			for i, endLineNumber := range endLineNumbers {
				if i > 0 && tfError.Range.Start.Line <= endLineNumber && tfError.Range.Start.Line > endLineNumbers[i-1] {
					validatedObject = newObjects[i]
					break
				} else if i == 0 && tfError.Range.Start.Line <= endLineNumber {
					validatedObject = newObjects[i]
					break
				}
			}

			if validatedObject != nil {
				*validationResults = append(*validationResults, &backendTypes.CodeValidationResult{
					ValidatedObject: validatedObject,
					ErrorMessage: fmt.Sprintf(
						"Summary: %s\nDetail: %s\nCode line: %s",
						tfError.Summary,
						tfError.Detail,
						tfError.Snippet.Code,
					),
				})
			}
		}
	}
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
