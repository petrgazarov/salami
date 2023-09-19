package lock_file_manager

import (
	"fmt"
	"regexp"

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
			default:
				return err
			}
		}
	}
	for _, object := range lockFile.Objects {
		if err := object.Parsed.validate(); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return err
			}

			for _, err := range err.(validator.ValidationErrors) {
				namespace := err.Namespace()
				switch err.Tag() {
				case "required":
					return getMissingFieldError(namespace)
				default:
					return err
				}
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
	SourceFilePath string        `toml:"source_file_path" validate:"required"`
	Parsed         ParsedObject  `toml:"parsed" validate:"required"`
	CodeSegments   []CodeSegment `toml:"code_segments" validate:"required,dive"`
}

type ParsedObject interface {
	validate() error
	getObjectType() string
}

type ParsedVariable struct {
	ObjectType   string `toml:"object_type" validate:"required,eq=Variable"`
	Name         string `toml:"name" validate:"required"`
	Description  string `toml:"description"`
	VariableType string `toml:"type" validate:"required,oneof=string number boolean"`
	DefaultValue string `toml:"default"`
}

func (v ParsedVariable) validate() error {
	validate := validator.New()
	return validate.Struct(v)
}

func (v ParsedVariable) getObjectType() string {
	return v.ObjectType
}

type ParsedResource struct {
	ObjectType          string            `toml:"object_type" validate:"required,eq=Resource"`
	ResourceType        string            `toml:"resource_type" validate:"required"`
	LogicalName         string            `toml:"logical_name" validate:"required"`
	NaturalLanguage     string            `toml:"natural_language"`
	Uses                []string          `toml:"uses"`
	Exports             map[string]string `toml:"exports"`
	ReferencedVariables []string          `toml:"referenced_variables"`
}

func (r ParsedResource) validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r ParsedResource) getObjectType() string {
	return r.ObjectType
}

type CodeSegment struct {
	SegmentType    string `toml:"segment_type" validate:"required,oneof=Variable Resource Export"`
	TargetFilePath string `toml:"target_file_path" validate:"required"`
	Content        string `toml:"content" validate:"required"`
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
