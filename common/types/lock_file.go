package types

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type LockFile struct {
	Version        string           `toml:"version" validate:"required"`
	SourceDir      string           `toml:"source_dir" validate:"required"`
	DestinationDir string           `toml:"destination_dir" validate:"required"`
	TargetFiles    []TargetFile     `toml:"target_files"`
	Objects        []LockFileObject `toml:"objects"`
}

type TargetFile struct {
	FilePath string `toml:"file_path" validate:"required"`
	Checksum string `toml:"checksum" validate:"required"`
}

type LockFileObject struct {
	SourceFilePath string                 `toml:"source_file_path" validate:"required"`
	Parsed         LockFileParsedObject   `toml:"-"`
	CodeSegments   []LockFileCodeSegment  `toml:"code_segments" validate:"required"`
	RawParsed      map[string]interface{} `toml:"parsed"`
}

type LockFileParsedObject interface {
	Validate() error
	GetObjectType() string
}

type LockFileVariable struct {
	ObjectType  string `toml:"object_type" validate:"required,eq=Variable"`
	Description string `toml:"description"`
	Type        string `toml:"type" validate:"required oneof=string number boolean"`
	Default     string `toml:"default"`
}

func (v LockFileVariable) Validate() error {
	validate := validator.New()
	return validate.Struct(v)
}

func (v LockFileVariable) GetObjectType() string {
	return v.ObjectType
}

type LockFileResource struct {
	ObjectType          string            `toml:"object_type" validate:"required,eq=Resource"`
	ResourceType        string            `toml:"resource_type" validate:"required"`
	LogicalName         string            `toml:"logical_name" validate:"required"`
	NaturalLanguage     string            `toml:"natural_language"`
	Uses                []string          `toml:"uses"`
	Exports             map[string]string `toml:"exports"`
	ReferencedVariables []string          `toml:"referenced_variables"`
}

func (r LockFileResource) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r LockFileResource) GetObjectType() string {
	return r.ObjectType
}

type LockFileCodeSegment struct {
	SegmentType string `toml:"segment_type" validate:"required oneof=Variable Resource"`
	FilePath    string `toml:"file_path" validate:"required"`
	Content     string `toml:"content" validate:"required"`
}

func (l *LockFileObject) UnmarshalTOML(data []byte) error {
	type temp LockFileObject
	if err := toml.Unmarshal(data, (*temp)(l)); err != nil {
		return err
	}

	switch l.Parsed.GetObjectType() {
	case "Variable":
		var v LockFileVariable
		if err := mapstructure.Decode(l.RawParsed, &v); err != nil {
			return err
		}
		l.Parsed = v
	case "Resource":
		var r LockFileResource
		if err := mapstructure.Decode(l.RawParsed, &r); err != nil {
			return err
		}
		l.Parsed = r
	default:
		return fmt.Errorf("unknown object type: %s", l.Parsed.GetObjectType())
	}

	return nil
}
