package llm

import (
	openaiGpt4 "salami/backend/llm/openai/gpt4"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

var llmFuncsMap = map[string]backendTypes.NewLlmFunc{
	commonTypes.LlmOpenaiGpt4: openaiGpt4.NewLlm,
}

func ResolveLlm(llmConfig commonTypes.LlmConfig) backendTypes.Llm {
	var llm backendTypes.Llm
	if llmConfig.Provider == commonTypes.LlmOpenaiProvider && llmConfig.Model == commonTypes.LlmGpt4Model {
		llm = llmFuncsMap[commonTypes.LlmOpenaiGpt4](llmConfig)
	}
	return llm
}
