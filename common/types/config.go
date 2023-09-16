package types

type CompilerConfig struct {
	Target    CompilerTargetConfig `yaml:"target"`
	Llm       CompilerLlmConfig    `yaml:"llm"`
	SourceDir string               `yaml:"source_dir"`
	TargetDir string               `yaml:"target_dir"`
}

type Config struct {
	Compiler CompilerConfig `yaml:"compiler"`
}

type CompilerTargetConfig struct {
	Platform string
}

type CompilerLlmConfig struct {
	Provider string
	Model    string
}

const TerraformPlatform = "terraform"
const LlmOpenaiProvider = "openai"
const LlmGpt4Model = "gpt-4"
