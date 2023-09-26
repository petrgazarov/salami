package gpt4

import (
	"context"
	"encoding/json"
	"math"
	backendTypes "salami/backend/types"
	"salami/common/errors"
	commonTypes "salami/common/types"

	"github.com/sashabaranov/go-openai"
)

func CreateCompletion(
	messages []*backendTypes.LlmMessage,
	llmConfig commonTypes.LlmConfig,
) (string, error) {
	client := getClient(llmConfig)
	response, err := client.CreateChatCompletion(context.Background(), getChatCompletionRequest(messages, llmConfig))
	if err != nil {
		return "", err
	}
	functionCall := response.Choices[0].Message.FunctionCall
	if functionCall == nil {
		return "", &errors.LlmError{Message: "Function call is nil"}
	}
	var parsedArguments map[string]interface{}
	err = json.Unmarshal([]byte(functionCall.Arguments), &parsedArguments)
	if err != nil {
		return "", err
	}
	code, ok := parsedArguments["code"].(string)
	if !ok {
		return "", &errors.LlmError{Message: "Code is not a string"}
	}
	return code, nil
}

var initializedClient *openai.Client

func getClient(llmConfig commonTypes.LlmConfig) *openai.Client {
	if initializedClient == nil {
		initializedClient = openai.NewClient(llmConfig.ApiKey)
	}
	return initializedClient
}

func getModel(llmConfig commonTypes.LlmConfig) string {
	switch llmConfig.Model {
	case commonTypes.LlmGpt4Model:
		return openai.GPT4
	}
	return ""
}

func getChatCompletionRequest(
	messages []*backendTypes.LlmMessage,
	llmConfig commonTypes.LlmConfig,
) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:        getModel(llmConfig),
		Messages:     getMessages(messages),
		Temperature:  math.SmallestNonzeroFloat32,
		Functions:    getFunctions(),
		FunctionCall: map[string]string{"name": "save_code"},
	}
}

func getMessages(messages []*backendTypes.LlmMessage) []openai.ChatCompletionMessage {
	var openaiMessages []openai.ChatCompletionMessage
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    string(message.Role),
			Content: message.Content,
		})
	}
	return openaiMessages
}

func getFunctions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        "save_code",
			Description: "Save the provided code",
			Parameters:  []byte(`{"type": "object", "properties": {"code": {"type": "string", "description": "Valid and runnable code"}}}`),
		},
	}
}
