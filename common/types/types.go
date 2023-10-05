package types

type ResourceType string
type LogicalName string

type ParsedResource struct {
	ResourceType        ResourceType
	LogicalName         LogicalName
	NaturalLanguage     string
	ReferencedResources []LogicalName
	ReferencedVariables []string
	SourceFilePath      string
	SourceFileLine      int
}

func NewParsedResource(SourceFilePath string, SourceFileLine int) *ParsedResource {
	return &ParsedResource{
		ReferencedResources: []LogicalName{},
		ReferencedVariables: []string{},
		SourceFilePath:      SourceFilePath,
		SourceFileLine:      SourceFileLine,
	}
}

func (r *ParsedResource) AddNaturalLanguage(NaturalLanguage string) {
	if r.NaturalLanguage == "" {
		r.NaturalLanguage = NaturalLanguage
	} else {
		r.NaturalLanguage = r.NaturalLanguage + "\n" + NaturalLanguage
	}
}

type VariableType string

type ParsedVariable struct {
	Name            string
	NaturalLanguage string
	Default         string
	Type            VariableType
	SourceFilePath  string
	SourceFileLine  int
}

func NewParsedVariable(SourceFilePath string, SourceFileLine int) *ParsedVariable {
	return &ParsedVariable{
		SourceFilePath: SourceFilePath,
		SourceFileLine: SourceFileLine,
	}
}

func (v *ParsedVariable) AddNaturalLanguage(NaturalLanguage string) {
	if v.NaturalLanguage == "" {
		v.NaturalLanguage = NaturalLanguage
	} else {
		v.NaturalLanguage = v.NaturalLanguage + "\n" + NaturalLanguage
	}
}

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
