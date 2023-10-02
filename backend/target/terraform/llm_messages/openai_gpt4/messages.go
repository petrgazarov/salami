package openai_gpt4

import (
	"salami/backend/llm/openai/gpt4"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
)

func GetMessages(
	changeSetDiff *commonTypes.ChangeSetDiff,
	symbolTable *symbol_table.SymbolTable,
) ([]*gpt4.LlmMessage, error) {
	var messages []*gpt4.LlmMessage

	systemMessage, err := getSystemMessage()
	if err != nil {
		return nil, err
	}

	var newMessage *gpt4.LlmMessage
	if changeSetDiff.IsAdd() {
		newMessage, err = getNewMessage(changeSetDiff.NewObject, symbolTable)
	} else if changeSetDiff.IsUpdate() {
		newMessage, err = getNewMessage(changeSetDiff.OldObject, symbolTable)
	}
	if err != nil {
		return nil, err
	}
	messages = append(messages, systemMessage, newMessage)

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

func getSystemMessage() (*gpt4.LlmMessage, error) {
	systemMessageContent, err := readTemplateFile("templates/system.txt")
	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageRole("system"),
		Content: systemMessageContent,
	}, nil
}

func getNewMessage(
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (*gpt4.LlmMessage, error) {
	var userMessageContent string
	var err error
	if object.IsResource() {
		userMessageContent, err = populateResourceTemplate("templates/new/resource.txt", object, symbolTable)
	} else if object.IsVariable() {
		userMessageContent, err = populateVariableTemplate("templates/new/variable.txt", object)
	}
	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageRole("user"),
		Content: userMessageContent,
	}, nil
}

func getUpdateMessage(
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (*gpt4.LlmMessage, error) {
	var messageContent string
	var err error
	if object.IsResource() {
		messageContent, err = populateResourceTemplate("templates/update/resource.txt", object, symbolTable)
	} else if object.IsVariable() {
		messageContent, err = populateVariableTemplate("templates/update/variable.txt", object)
	}
	if err != nil {
		return nil, err
	}

	return &gpt4.LlmMessage{
		Role:    gpt4.LlmMessageRole("user"),
		Content: messageContent,
	}, nil
}
