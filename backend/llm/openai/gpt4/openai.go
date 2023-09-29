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

type OpenaiGpt4 struct {
	slug   string
	model  string
	client *openai.Client
}

func NewLlm(llmConfig commonTypes.LlmConfig) backendTypes.Llm {
	return &OpenaiGpt4{
		client: getClient(llmConfig),
		model:  getModel(llmConfig),
		slug:   commonTypes.LlmOpenaiGpt4,
	}
}

func (o *OpenaiGpt4) GetSlug() string {
	return o.slug
}

func (o *OpenaiGpt4) CreateCompletion(messages []*backendTypes.LlmMessage) (string, error) {
	response, err := o.client.CreateChatCompletion(context.Background(), o.getChatCompletionRequest(messages))
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

func getClient(llmConfig commonTypes.LlmConfig) *openai.Client {
	return openai.NewClient(llmConfig.ApiKey)
}

func getModel(llmConfig commonTypes.LlmConfig) string {
	switch llmConfig.Model {
	case commonTypes.LlmGpt4Model:
		return openai.GPT4
	}
	return ""
}

func (o *OpenaiGpt4) getChatCompletionRequest(
	messages []*backendTypes.LlmMessage,
) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:        o.model,
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
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "Valid and runnable code",
					},
				},
			},
		},
	}
}
