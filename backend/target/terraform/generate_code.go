package terraform

import (
	"salami/backend/prompts/terraform/openai_gpt4"
	backendTypes "salami/backend/types"
	"salami/common/change_set"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"

	"golang.org/x/sync/errgroup"
)

type Terraform struct{}

func NewTarget() backendTypes.Target {
	return &Terraform{}
}

func (t *Terraform) GenerateCode(
	symbolTable *symbol_table.SymbolTable,
	changeSetRepository *change_set.ChangeSetRepository,
	llm backendTypes.Llm,
) []error {
	var g errgroup.Group
	semaphoreChannel := make(chan struct{}, llm.GetMaxConcurrentExecutions())

	for _, diff := range changeSetRepository.Diffs {
		if !(diff.IsUpdate() || diff.IsAdd()) {
			continue
		}
		diff := diff
		g.Go(func() error {
			semaphoreChannel <- struct{}{}
			defer func() { <-semaphoreChannel }()

			messages, err := getGenerateCodeLlmMessages(symbolTable, diff, llm)
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

func getGenerateCodeLlmMessages(
	symbolTable *symbol_table.SymbolTable,
	diff *commonTypes.ChangeSetDiff,
	llm backendTypes.Llm,
) ([]interface{}, error) {
	var messages []interface{}

	switch llm.GetSlug() {
	case commonTypes.LlmOpenaiGpt4:
		llmMessages, err := openai_gpt4.GetGenerateCodeMessages(symbolTable, diff)
		if err != nil {
			return nil, err
		}

		for _, v := range llmMessages {
			messages = append(messages, v)
		}
	}
	return messages, nil
}
