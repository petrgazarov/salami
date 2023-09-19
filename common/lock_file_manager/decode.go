package lock_file_manager

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

func decodeLockFile(lockFile *LockFile) error {
	if _, err := toml.DecodeFile(lockFilePath, lockFile); err != nil {
		if err != nil && !os.IsNotExist(err) {
			return &LockFileError{Message: "could not parse lock file"}
		}
	}
	loadedLockFile = lockFile
	return nil
}

// UnmarshalTOML function is called by toml.DecodeFile
func (o *Object) UnmarshalTOML(data interface{}) error {
	var rawObject map[string]interface{}
	decoder, err := mapstructure.NewDecoder(getDecoderConfig(&rawObject))
	if err != nil {
		return err
	}
	if err := decoder.Decode(data); err != nil {
		return err
	}

	if value, ok := rawObject["source_file_path"].(string); ok {
		o.SourceFilePath = value
	}
	rawCodeSegments := rawObject["code_segments"].([]map[string]interface{})
	o.CodeSegments = make([]CodeSegment, len(rawCodeSegments))
	err = decodeCodeSegments(o, rawCodeSegments)
	if err != nil {
		return err
	}
	rawParsedObject := rawObject["parsed"].(map[string]interface{})
	err = decodeParsedObject(o, rawParsedObject)
	if err != nil {
		return err
	}

	return nil
}

func decodeCodeSegments(o *Object, rawCodeSegments []map[string]interface{}) error {
	for i, rawCodeSegment := range rawCodeSegments {
		var cs CodeSegment
		decoder, err := mapstructure.NewDecoder(getDecoderConfig(&cs))
		if err != nil {
			return err
		}
		if err := decoder.Decode(rawCodeSegment); err != nil {
			return err
		}
		o.CodeSegments[i] = cs
	}
	return nil
}

func decodeParsedObject(o *Object, rawParsedObject map[string]interface{}) error {
	switch rawParsedObject["object_type"] {
	case "Variable":
		var v ParsedVariable
		decoder, err := mapstructure.NewDecoder(getDecoderConfig(&v))
		if err != nil {
			return err
		}
		if err := decoder.Decode(rawParsedObject); err != nil {
			return err
		}
		o.Parsed = v
	case "Resource":
		var r ParsedResource
		decoder, err := mapstructure.NewDecoder(getDecoderConfig(&r))
		if err != nil {
			return err
		}
		if err := decoder.Decode(rawParsedObject); err != nil {
			return err
		}
		o.Parsed = r
	default:
		o.Parsed = ParsedResource{}
	}
	return nil
}

func getDecoderConfig(result interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   result,
		TagName:  "toml",
	}
}
