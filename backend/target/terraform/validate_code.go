package terraform

import (
	"salami/backend/prompts/terraform/openai_gpt4"
	backendTypes "salami/backend/types"
	"salami/common/change_set"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"

	"golang.org/x/sync/errgroup"
)

func (t *Terraform) ValidateCode(
	newObjects []*commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
	changeSetRepository *change_set.ChangeSetRepository,
	llm backendTypes.Llm,
) []error {
	validationResults, errors := runTerraformValidate(newObjects)
	if len(errors) > 0 {
		return errors
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

func runTerraformValidate(
	newObjects []*commonTypes.Object,
) ([]*backendTypes.CodeValidationResult, []error) {
	return nil, nil
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
