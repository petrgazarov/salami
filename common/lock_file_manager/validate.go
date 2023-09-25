package lock_file_manager

import (
	"fmt"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

func ValidateLockFile() error {
	lockFile := &LockFile{}
	decodeLockFile(lockFile)

	if isEmptyLockFile(lockFile) {
		return nil
	}
	validate := newValidator()
	if err := validate.Struct(lockFile); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fieldValue := err.Value()
			namespace := err.Namespace()
			switch err.Tag() {
			case "semver":
				return &LockFileError{Message: fmt.Sprintf("'%s' is not a valid semver", fieldValue)}
			case "required":
				return getMissingFieldError(namespace)
			case "oneof":
				return &LockFileError{Message: fmt.Sprintf("'%s' is not a valid value", fieldValue)}
			default:
				return err
			}
		}
	}
	return nil
}

type LockFile struct {
	Version         string           `toml:"version" validate:"required,semver"`
	TargetFilesMeta []TargetFileMeta `toml:"target_files_meta" validate:"dive"`
	Objects         []Object         `toml:"objects" validate:"dive"`
}

type TargetFileMeta struct {
	FilePath string `toml:"file_path" validate:"required"`
	Checksum string `toml:"checksum" validate:"required"`
}

type Object struct {
	SourceFilePath string          `toml:"source_file_path" validate:"required"`
	ParsedResource *ParsedResource `toml:"parsed_resource" validate:"required_without=ParsedVariable"`
	ParsedVariable *ParsedVariable `toml:"parsed_variable" validate:"required_without=ParsedResource"`
	CodeSegments   []CodeSegment   `toml:"code_segments" validate:"required,dive"`
}

func (o *Object) IsResource() bool {
	return o.ParsedResource != nil
}

func (o *Object) IsVariable() bool {
	return o.ParsedVariable != nil
}

type ParsedVariable struct {
	Name         string `toml:"name" validate:"required"`
	Description  string `toml:"description"`
	VariableType string `toml:"type" validate:"required,oneof=string number boolean"`
	DefaultValue string `toml:"default"`
}

type ParsedResource struct {
	ResourceType        string            `toml:"resource_type" validate:"required"`
	LogicalName         string            `toml:"logical_name" validate:"required"`
	NaturalLanguage     string            `toml:"natural_language"`
	Uses                []string          `toml:"uses"`
	Exports             map[string]string `toml:"exports"`
	ReferencedVariables []string          `toml:"referenced_variables"`
}

type CodeSegment struct {
	SegmentType    string `toml:"segment_type" validate:"required,oneof=Variable Resource Export"`
	Content        string `toml:"content" validate:"required"`
}

func decodeLockFile(lockFile *LockFile) error {
	if _, err := toml.DecodeFile(lockFilePath, lockFile); err != nil {
		if err != nil && !os.IsNotExist(err) {
			return &LockFileError{Message: "could not parse lock file"}
		}
	}
	loadedLockFile = lockFile
	return nil
}

func validateSemVer(fl validator.FieldLevel) bool {
	numericVersionRegex := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`
	match, _ := regexp.MatchString(numericVersionRegex, fl.Field().String())
	return match
}

func newValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("semver", validateSemVer)
	return validate
}

func isEmptyLockFile(lf *LockFile) bool {
	return lf.Version == "" &&
		len(lf.TargetFilesMeta) == 0 &&
		len(lf.Objects) == 0
}
