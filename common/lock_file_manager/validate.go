package lock_file_manager

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ValidateLockFile() error {
	lockFile := getLockFile()
	if isEmptyLockFile(lockFile) {
		return nil
	}
	validate := newValidator()
	if err := validate.Struct(lockFile); err != nil {
		return err
	}
	for _, object := range lockFile.Objects {
		if err := object.Parsed.validate(); err != nil {
			return err
		}
	}
	return nil
}

type LockFile struct {
	Version         string           `toml:"version" validate:"required,semver"`
	SourceDir       string           `toml:"source_dir" validate:"required"`
	DestinationDir  string           `toml:"destination_dir" validate:"required"`
	TargetFilesMeta []TargetFileMeta `toml:"target_files_meta"`
	Objects         []Object         `toml:"objects"`
}

type TargetFileMeta struct {
	FilePath string `toml:"file_path" validate:"required"`
	Checksum string `toml:"checksum" validate:"required"`
}

type Object struct {
	SourceFilePath string                 `toml:"source_file_path" validate:"required"`
	Parsed         ParsedObject           `toml:"-"`
	CodeSegments   []CodeSegment          `toml:"code_segments" validate:"required"`
	RawParsed      map[string]interface{} `toml:"parsed"`
}

type ParsedObject interface {
	validate() error
	getObjectType() string
}

type ParsedVariable struct {
	ObjectType  string `toml:"object_type" validate:"required,eq=Variable"`
	Name        string `toml:"name" validate:"required"`
	Description string `toml:"description"`
	Type        string `toml:"type" validate:"required oneof=string number boolean"`
	Default     string `toml:"default"`
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
	SegmentType    string `toml:"segment_type" validate:"required oneof=Variable Resource"`
	TargetFilePath string `toml:"target_file_path" validate:"required"`
	Content        string `toml:"content" validate:"required"`
}

func (l *Object) UnmarshalTOML(data []byte) error {
	type temp Object
	if err := toml.Unmarshal(data, (*temp)(l)); err != nil {
		return err
	}

	switch l.Parsed.getObjectType() {
	case "Variable":
		var v ParsedVariable
		if err := mapstructure.Decode(l.RawParsed, &v); err != nil {
			return err
		}
		l.Parsed = v
	case "Resource":
		var r ParsedResource
		if err := mapstructure.Decode(l.RawParsed, &r); err != nil {
			return err
		}
		l.Parsed = r
	default:
		return fmt.Errorf("unknown object type: %s", l.Parsed.getObjectType())
	}

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
		lf.SourceDir == "" &&
		lf.DestinationDir == "" &&
		len(lf.TargetFilesMeta) == 0 &&
		len(lf.Objects) == 0
}
