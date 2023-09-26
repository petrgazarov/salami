package llm

import (
	openaiGpt4 "salami/backend/llm/openai/gpt4"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

var llmFuncsMap = map[string]backendTypes.Llm{
	"openai_gpt4": {
		CreateCompletion: openaiGpt4.CreateCompletion,
	},
}

func ResolveLlm(llmConfig commonTypes.LlmConfig) backendTypes.Llm {
	var llm backendTypes.Llm
	if llmConfig.Provider == commonTypes.LlmOpenaiProvider && llmConfig.Model == commonTypes.LlmGpt4Model {
		llm = llmFuncsMap["openai_gpt4"]
	}
	return llm
}
