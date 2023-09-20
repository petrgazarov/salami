package config

import (
	"fmt"
	"os"
	"path/filepath"
	"salami/common/types"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func ValidateConfig() error {
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &loadedConfig); err != nil {
		return &ConfigError{Message: "could not parse config file. Ensure it is valid yaml format"}
	}
	validate := newValidator()

	if err = validate.Struct(loadedConfig); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fieldValue := err.Value()
			namespace := err.Namespace()
			switch err.Tag() {
			case "valid_target":
				return getInvalidFieldError(namespace, nil)
			case "valid_llm":
				return getInvalidFieldError(namespace, nil)
			case "dir_exists":
				return &ConfigError{Message: fmt.Sprintf("'%s' directory could not be resolved", fieldValue)}
			case "target_dir_valid":
				return &ConfigError{
					Message: "target directory cannot be outside of the current working directory",
				}
			case "required":
				return getMissingFieldError(namespace)
			default:
				return err
			}
		}
	}
	return nil
}

type CompilerConfig struct {
	Target    CompilerTargetConfig `yaml:"target" validate:"valid_target"`
	Llm       CompilerLlmConfig    `yaml:"llm" validate:"valid_llm"`
	SourceDir string               `yaml:"source_dir" validate:"required,dir_exists"`
	TargetDir string               `yaml:"target_dir" validate:"required,target_dir_valid"`
}

type ConfigType struct {
	Compiler CompilerConfig `yaml:"compiler" validate:"required"`
}

type CompilerTargetConfig struct {
	Platform string `yaml:"platform" validate:"required"`
}

type CompilerLlmConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

func validateTarget(fl validator.FieldLevel) bool {
	target, ok := fl.Field().Interface().(CompilerTargetConfig)
	if !ok {
		return true
	}
	return target.Platform == types.TerraformPlatform
}

func validateLlm(fl validator.FieldLevel) bool {
	llmConfig, ok := fl.Field().Interface().(CompilerLlmConfig)
	if !ok {
		return true
	}
	provider := llmConfig.Provider
	model := llmConfig.Model
	return provider == types.LlmOpenaiProvider && model == types.LlmGpt4Model
}

func validateDirExists(fl validator.FieldLevel) bool {
	dir := fl.Field().String()
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

// Target directory cannot be outside of the current working directory
func validateTargetDir(fl validator.FieldLevel) bool {
	targetDir := fl.Field().String()
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return false
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return false
	}

	rel, err := filepath.Rel(rootDir, absTargetDir)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}

func newValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("valid_target", validateTarget)
	validate.RegisterValidation("valid_llm", validateLlm)
	validate.RegisterValidation("dir_exists", validateDirExists)
	validate.RegisterValidation("target_dir_valid", validateTargetDir)
	return validate
}
