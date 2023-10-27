package gpt4

import (
	"context"
	"encoding/json"
	"math"
	backendTypes "salami/backend/types"
	"salami/common/constants"
	"salami/common/errors"
	"salami/common/logger"
	commonTypes "salami/common/types"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const LlmMessageUserRole = "user"
const LlmMessageFunctionRole = "function"
const LlmMessageSystemRole = "system"
const LlmMessageAssistantRole = "assistant"

const functionCallName = "save_code"

type OpenaiGpt4 struct {
	slug                    string
	model                   string
	client                  *openai.Client
	maxConcurrentExecutions int
}

type LlmMessage struct {
	Role         string
	Content      string
	Name         string
	FunctionCall *openai.FunctionCall
}

func NewLlm(llmConfig commonTypes.LlmConfig) backendTypes.Llm {
	return &OpenaiGpt4{
		client:                  getClient(llmConfig),
		model:                   openai.GPT4,
		slug:                    commonTypes.LlmOpenaiGpt4,
		maxConcurrentExecutions: llmConfig.MaxConcurrentExecutions,
	}
}

func GetAssistantMessage(code string) (*LlmMessage, error) {
	codeBytes, err := json.Marshal(code)
	if err != nil {
		return nil, err
	}
	codeJson := string(codeBytes)
	return &LlmMessage{
		Role:         LlmMessageAssistantRole,
		FunctionCall: &openai.FunctionCall{Name: functionCallName, Arguments: `{"code": ` + codeJson + `}`},
	}, nil
}

func GetFunctionMessage() *LlmMessage {
	return &LlmMessage{
		Role:    LlmMessageFunctionRole,
		Name:    functionCallName,
		Content: "true",
	}
}

func (o *OpenaiGpt4) GetSlug() string {
	return o.slug
}

func (o *OpenaiGpt4) GetMaxConcurrentExecutions() int {
	return o.maxConcurrentExecutions
}

func (o *OpenaiGpt4) CreateCompletion(messages []interface{}) (string, error) {
	llmMessages := make([]*LlmMessage, len(messages))
	for i, message := range messages {
		llmMessage := message.(*LlmMessage)
		llmMessages[i] = llmMessage
	}

	logMessages(llmMessages)

	response, err := o.callApi(llmMessages)
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

	return strings.TrimSpace(code), nil
}

func (o *OpenaiGpt4) callApi(llmMessages []*LlmMessage) (openai.ChatCompletionResponse, error) {
	var response openai.ChatCompletionResponse
	var err error

	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(constants.LlmTimeoutDurationSeconds)*time.Second,
		)
		defer cancel()

		response, err = o.client.CreateChatCompletion(
			ctx,
			o.getChatCompletionRequest(llmMessages),
		)
		if err == nil {
			return response, err
		}
	}

	return response, err
}

func getClient(llmConfig commonTypes.LlmConfig) *openai.Client {
	return openai.NewClient(llmConfig.ApiKey)
}

func (o *OpenaiGpt4) getChatCompletionRequest(
	messages []*LlmMessage,
) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:        o.model,
		Messages:     getMessages(messages),
		Temperature:  math.SmallestNonzeroFloat32,
		Functions:    getFunctions(),
		FunctionCall: map[string]string{"name": functionCallName},
	}
}

func getMessages(messages []*LlmMessage) []openai.ChatCompletionMessage {
	var openaiMessages []openai.ChatCompletionMessage
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:         string(message.Role),
			Name:         message.Name,
			FunctionCall: message.FunctionCall,
			Content:      message.Content,
		})
	}
	return openaiMessages
}

func getFunctions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        functionCallName,
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

func logMessages(llmMessages []*LlmMessage) {
	str := "\n"

	for _, llmMessage := range llmMessages {
		if llmMessage.Content == "" {
			str += llmMessage.FunctionCall.Arguments + "\n\n"
		} else {
			str += llmMessage.Content + "\n\n"
		}
	}
	logger.Debug(str)
}
