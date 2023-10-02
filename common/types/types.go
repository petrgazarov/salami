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

func (o *Object) SetTargetCode(targetCode string) {
	o.TargetCode = targetCode
}

const (
	DiffTypeAdd    = "add"
	DiffTypeRemove = "remove"
	DiffTypeUpdate = "update"
	DiffTypeMove   = "move"
)

type ChangeSetDiff struct {
	OldObject *Object
	NewObject *Object
	DiffType  string
}

func (ch *ChangeSetDiff) IsUpdate() bool {
	return ch.DiffType == DiffTypeUpdate
}

func (ch *ChangeSetDiff) IsAdd() bool {
	return ch.DiffType == DiffTypeAdd
}

type ChangeSet struct {
	Diffs []*ChangeSetDiff
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
const LlmOpenaiGpt4 = LlmOpenaiProvider + "_" + LlmGpt4Model

type TargetFile struct {
	FilePath string
	Content  string
}
