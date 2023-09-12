package types

type CompilerTargetConfig struct {
	Platform string
}

const TerraformPlatform = "terraform"

type CompilerLlmConfig struct {
	Provider string
	Model    string
}

const LlmOpenaiProvider = "openai"
const LlmGpt4Model = "gpt-4"
