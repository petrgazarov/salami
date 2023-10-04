package llm

import (
	openaiGpt4 "salami/backend/llm/openai/gpt4"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

func ResolveLlm(llmConfig commonTypes.LlmConfig) backendTypes.Llm {
	var llm backendTypes.Llm
	if llmConfig.Provider == commonTypes.LlmOpenaiProvider && llmConfig.Model == commonTypes.LlmGpt4Model {
		llm = openaiGpt4.NewLlm(llmConfig)
	}
	return llm
}
