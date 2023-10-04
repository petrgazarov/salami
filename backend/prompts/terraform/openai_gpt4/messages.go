package openai_gpt4

import (
	"salami/backend/llm/openai/gpt4"
	"salami/backend/prompts/terraform/openai_gpt4/templates"
	backendTypes "salami/backend/types"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
)

func GetGenerateCodeMessages(
	symbolTable *symbol_table.SymbolTable,
	changeSetDiff *commonTypes.ChangeSetDiff,
) ([]*gpt4.LlmMessage, error) {
	var messages []*gpt4.LlmMessage

	systemMessage, err := getSystemMessage()
	if err != nil {
		return nil, err
	}

	var createMessage *gpt4.LlmMessage
	if changeSetDiff.IsAdd() {
		createMessage, err = getCreateMessage(changeSetDiff.NewObject, symbolTable)
	} else if changeSetDiff.IsUpdate() {
		createMessage, err = getCreateMessage(changeSetDiff.OldObject, symbolTable)
	}
	if err != nil {
		return nil, err
	}
	messages = append(messages, systemMessage, createMessage)

	if changeSetDiff.IsUpdate() {
		assistantMessage, err := gpt4.GetAssistantMessage(changeSetDiff.OldObject.TargetCode)
		if err != nil {
			return nil, err
		}
		functionMessage := gpt4.GetFunctionMessage()
		updateMessage, err := getUpdateMessage(changeSetDiff.NewObject, symbolTable)
		if err != nil {
			return nil, err
		}
		messages = append(messages, assistantMessage, functionMessage, updateMessage)
	}

	return messages, nil
}

func GetValidateCodeMessages(
	symbolTable *symbol_table.SymbolTable,
	changeSetDiff *commonTypes.ChangeSetDiff,
	validationResult *backendTypes.CodeValidationResult,
) ([]*gpt4.LlmMessage, error) {
	generateCodeMessages, err := GetGenerateCodeMessages(symbolTable, changeSetDiff)
	if err != nil {
		return nil, err
	}

	assistantMessage, err := gpt4.GetAssistantMessage(changeSetDiff.NewObject.TargetCode)
	if err != nil {
		return nil, err
	}
	functionMessage := gpt4.GetFunctionMessage()

	validationMessage, err := getValidationMessage(validationResult)
	if err != nil {
		return nil, err
	}

	return append(generateCodeMessages, assistantMessage, functionMessage, validationMessage), nil
}

func getSystemMessage() (*gpt4.LlmMessage, error) {
	systemMessageContent, err := templates.ReadTemplateFile(templates.SystemTemplatePath)

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageSystemRole,
		Content: systemMessageContent,
	}, err
}

func getCreateMessage(
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (*gpt4.LlmMessage, error) {
	var messageContent string
	var err error

	if object.IsResource() {
		messageContent, err = populateCreateUpdateTemplate(
			templates.CreateResourceTemplatePath,
			object,
			symbolTable,
		)
	} else if object.IsVariable() {
		messageContent, err = populateCreateUpdateTemplate(
			templates.CreateVariableTemplatePath,
			object,
			symbolTable,
		)
	}

	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageUserRole,
		Content: messageContent,
	}, nil
}

func getUpdateMessage(
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (*gpt4.LlmMessage, error) {
	var messageContent string
	var err error

	if object.IsResource() {
		messageContent, err = populateCreateUpdateTemplate(
			templates.UpdateResourceTemplatePath,
			object,
			symbolTable,
		)
	} else if object.IsVariable() {
		messageContent, err = populateCreateUpdateTemplate(
			templates.UpdateVariableTemplatePath,
			object,
			symbolTable,
		)
	}
	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageUserRole,
		Content: messageContent,
	}, nil
}

func getValidationMessage(
	validationResult *backendTypes.CodeValidationResult,
) (*gpt4.LlmMessage, error) {
	var messageContent string
	var err error

	if validationResult.ValidatedObject.IsResource() {
		messageContent, err = populateValidationTemplate(
			templates.FixResourceValidationIssueTemplatePath,
			validationResult,
		)
	} else if validationResult.ValidatedObject.IsVariable() {
		messageContent, err = populateValidationTemplate(
			templates.FixVariableValidationIssueTemplatePath,
			validationResult,
		)
	}

	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageUserRole,
		Content: messageContent,
	}, nil
}
