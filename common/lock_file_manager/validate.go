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
	for _, object := range lockFile.objects {
		if err := object.parsed.validate(); err != nil {
			return err
		}
	}
	return nil
}

type lockFile struct {
	version         string           `toml:"version" validate:"required,semver"`
	sourceDir       string           `toml:"source_dir" validate:"required"`
	destinationDir  string           `toml:"destination_dir" validate:"required"`
	targetFilesMeta []targetFileMeta `toml:"target_files_meta"`
	objects         []object         `toml:"objects"`
}

type targetFileMeta struct {
	filePath string `toml:"file_path" validate:"required"`
	checksum string `toml:"checksum" validate:"required"`
}

type object struct {
	sourceFilePath string                 `toml:"source_file_path" validate:"required"`
	parsed         parsedObject           `toml:"-"`
	codeSegments   []codeSegment          `toml:"code_segments" validate:"required"`
	rawParsed      map[string]interface{} `toml:"parsed"`
}

type parsedObject interface {
	validate() error
	getObjectType() string
}

type parsedVariable struct {
	objectType   string `toml:"object_type" validate:"required,eq=Variable"`
	name         string `toml:"name" validate:"required"`
	description  string `toml:"description"`
	variableType string `toml:"type" validate:"required oneof=string number boolean"`
	defaultValue string `toml:"default"`
}

func (v parsedVariable) validate() error {
	validate := validator.New()
	return validate.Struct(v)
}

func (v parsedVariable) getObjectType() string {
	return v.objectType
}

type parsedResource struct {
	objectType          string            `toml:"object_type" validate:"required,eq=Resource"`
	resourceType        string            `toml:"resource_type" validate:"required"`
	logicalName         string            `toml:"logical_name" validate:"required"`
	naturalLanguage     string            `toml:"natural_language"`
	uses                []string          `toml:"uses"`
	exports             map[string]string `toml:"exports"`
	referencedVariables []string          `toml:"referenced_variables"`
}

func (r parsedResource) validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r parsedResource) getObjectType() string {
	return r.objectType
}

type codeSegment struct {
	segmentType    string `toml:"segment_type" validate:"required oneof=Variable Resource"`
	targetFilePath string `toml:"target_file_path" validate:"required"`
	content        string `toml:"content" validate:"required"`
}

func (l *object) UnmarshalTOML(data []byte) error {
	type temp object
	if err := toml.Unmarshal(data, (*temp)(l)); err != nil {
		return err
	}

	switch l.parsed.getObjectType() {
	case "Variable":
		var v parsedVariable
		if err := mapstructure.Decode(l.rawParsed, &v); err != nil {
			return err
		}
		l.parsed = v
	case "Resource":
		var r parsedResource
		if err := mapstructure.Decode(l.rawParsed, &r); err != nil {
			return err
		}
		l.parsed = r
	default:
		return fmt.Errorf("unknown object type: %s", l.parsed.getObjectType())
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

func isEmptyLockFile(lf *lockFile) bool {
	return lf.version == "" &&
		lf.sourceDir == "" &&
		lf.destinationDir == "" &&
		len(lf.targetFilesMeta) == 0 &&
		len(lf.objects) == 0
}
