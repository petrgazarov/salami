package types

type Object struct {
	ParsedResource *ParsedResource
	ParsedVariable *ParsedVariable
	TargetCode     string
}

func (o *Object) IsResource() bool {
	return o.ParsedResource != nil
}

func (o *Object) IsVariable() bool {
	return o.ParsedVariable != nil
}

func (o *Object) GetSourceFilePath() string {
	if o.IsResource() {
		return o.ParsedResource.SourceFilePath
	} else if o.IsVariable() {
		return o.ParsedVariable.SourceFilePath
	}
	return ""
}

func (o *Object) GetSourceFileLine() int {
	if o.IsResource() {
		return o.ParsedResource.SourceFileLine
	} else if o.IsVariable() {
		return o.ParsedVariable.SourceFileLine
	}
	return 0
}

const (
	DiffTypeAdd     = "add"
	DiffTypeRemove  = "remove"
	DiffTypeReplace = "replace"
	DiffTypeUpdate  = "update"
)

type ChangeSetDiff struct {
	OldObject *Object
	NewObject *Object
	DiffType  string
}

type ChangeSet struct {
	Diffs []ChangeSetDiff
}

type TargetFileMeta struct {
	FilePath string
	Checksum string
}

type TargetConfig struct {
	Platform string
}

type LlmConfig struct {
	Provider string
	Model    string
	ApiKey   string
}

const TerraformPlatform = "terraform"
const LlmOpenaiProvider = "openai"
const LlmGpt4Model = "gpt4"

type TargetFile struct {
	FilePath string
	Content  string
}
