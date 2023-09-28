package terraform

import (
	"salami/backend/target/terraform/llm_messages/openai_gpt4"
	backendTypes "salami/backend/types"
	"salami/common/constants"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Terraform struct{}

func NewTarget() backendTypes.Target {
	return &Terraform{}
}

func (t *Terraform) GenerateCode(
	changeSet *commonTypes.ChangeSet,
	symbolTable *symbol_table.SymbolTable,
	llm backendTypes.Llm,
) []error {
	var g errgroup.Group

	for _, diff := range changeSet.Diffs {
		if !diff.ShouldGenerateCode() {
			continue
		}
		diff := diff
		g.Go(func() error {
			completion, err := llm.CreateCompletion(t.getMessages(diff, symbolTable, llm))
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

func (t *Terraform) GetFilesFromObjects(objects []*commonTypes.Object) []*commonTypes.TargetFile {
	targetFiles := []*commonTypes.TargetFile{}

	var currentTargetFile *commonTypes.TargetFile
	var lastSourceFilePath string
	for _, object := range objects {
		if object.GetSourceFilePath() != lastSourceFilePath {
			currentTargetFile = &commonTypes.TargetFile{
				FilePath: getTargetFilePath(object.GetSourceFilePath()),
				Content:  object.TargetCode,
			}
			targetFiles = append(targetFiles, currentTargetFile)
			lastSourceFilePath = object.GetSourceFilePath()
		} else {
			currentTargetFile.Content += ("\n\n" + object.TargetCode)
		}
	}

	return targetFiles
}

func (t *Terraform) getMessages(
	diff *commonTypes.ChangeSetDiff,
	symbolTable *symbol_table.SymbolTable,
	llm backendTypes.Llm,
) []*backendTypes.LlmMessage {
	switch llm.GetSlug() {
	case commonTypes.LlmOpenaiGpt4:
		return openai_gpt4.GetMessages(diff, symbolTable)
	}
	return nil
}

func getTargetFilePath(sourceFilePath string) string {
	if strings.HasSuffix(sourceFilePath, constants.SalamiFileExtension) {
		return strings.TrimSuffix(sourceFilePath, constants.SalamiFileExtension) + constants.TerraformFileExtension
	}
	return sourceFilePath
}
