package config

import (
	"salami/common/types"

	"github.com/go-playground/validator/v10"
)

func ValidateConfig() error {
	config := getConfig()
	validate := newValidator()

	err := validate.Struct(config)
	if err != nil {
		return err
	}

	return nil
}

type compilerConfig struct {
	target    compilerTargetConfig `yaml:"target" validate:"required,valid_target"`
	llm       compilerLlmConfig    `yaml:"llm" validate:"required,valid_llm"`
	sourceDir string               `yaml:"source_dir" validate:"required"`
	targetDir string               `yaml:"target_dir" validate:"required"`
}

type configType struct {
	compiler compilerConfig `yaml:"compiler" validate:"required"`
}

type compilerTargetConfig struct {
	platform string `yaml:"platform" validate:"required"`
}

type compilerLlmConfig struct {
	provider string `yaml:"provider" validate:"required"`
	model    string `yaml:"model" validate:"required"`
}

func validateTarget(fl validator.FieldLevel) bool {
	platform := fl.Field().String()
	return platform == types.TerraformPlatform
}

func validateLlm(fl validator.FieldLevel) bool {
	provider := fl.Field().Interface().(compilerLlmConfig).provider
	model := fl.Field().Interface().(compilerLlmConfig).model
	return provider == types.LlmOpenaiProvider && model == types.LlmGpt4Model
}

func newValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("valid_target", validateTarget)
	validate.RegisterValidation("valid_llm", validateLlm)
	return validate
}
