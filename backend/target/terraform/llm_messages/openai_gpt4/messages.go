package openai_gpt4

import (
	"bytes"
	"embed"
	"io/fs"
	"path/filepath"
	backendTypes "salami/backend/types"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"strings"
	"text/template"
)

func GetMessages(
	changeSetDiff *commonTypes.ChangeSetDiff,
	symbolTable *symbol_table.SymbolTable,
) ([]*backendTypes.LlmMessage, error) {
	templatesDirectory := getTemplatesDirectory(changeSetDiff)

	systemMessageContent, err := readTemplateFile(templatesDirectory, "system.txt")
	if err != nil {
		return nil, err
	}
	var userMessageContent string
	if changeSetDiff.NewObject.IsResource() {
		userMessageContent, err = populateResourceTemplate(templatesDirectory, changeSetDiff, symbolTable)
	} else if changeSetDiff.NewObject.IsVariable() {
		userMessageContent, err = populateVariableTemplate(templatesDirectory, changeSetDiff)
	}
	if err != nil {
		return nil, err
	}

	return []*backendTypes.LlmMessage{
		{
			Role:    backendTypes.LlmMessageRole("system"),
			Content: systemMessageContent,
		},
		{
			Role:    backendTypes.LlmMessageRole("user"),
			Content: userMessageContent,
		},
	}, nil
}

func getTemplatesDirectory(changeSetDiff *commonTypes.ChangeSetDiff) string {
	if changeSetDiff.DiffType == commonTypes.DiffTypeAdd {
		return "new"
	}

	newObject := changeSetDiff.NewObject
	oldObject := changeSetDiff.OldObject
	if newObject.IsResource() {
		shouldGenerateNewObject := oldObject.ParsedResource.ResourceType != newObject.ParsedResource.ResourceType

		if shouldGenerateNewObject {
			return "new"
		}
	}

	return "update"
}

func populateResourceTemplate(
	templatesDirectory string,
	changeSetDiff *commonTypes.ChangeSetDiff,
	symbolTable *symbol_table.SymbolTable,
) (string, error) {
	templateString, err := readTemplateFile(templatesDirectory, "resource.txt")
	if err != nil {
		return "", err
	}
	newObject := changeSetDiff.NewObject

	data := struct {
		ResourceType             string
		LogicalName              string
		NaturalLanguage          string
		ReferencedObjectsSection string
	}{
		ResourceType:             string(newObject.ParsedResource.ResourceType),
		LogicalName:              string(newObject.ParsedResource.LogicalName),
		NaturalLanguage:          newObject.ParsedResource.NaturalLanguage,
		ReferencedObjectsSection: getReferencedObjectsSection(changeSetDiff, symbolTable),
	}

	result, err := populateTemplate(templateString, data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func populateVariableTemplate(
	templateDirectory string,
	changeSetDiff *commonTypes.ChangeSetDiff,
) (string, error) {
	templateString, err := readTemplateFile(templateDirectory, "variable.txt")
	if err != nil {
		return "", err
	}
	newObject := changeSetDiff.NewObject

	data := struct {
		Name                   string
		VariableType           string
		VariableDetailsSection string
	}{
		Name:                   string(newObject.ParsedVariable.Name),
		VariableType:           string(newObject.ParsedVariable.Type),
		VariableDetailsSection: getVariableDetailsSection(changeSetDiff),
	}

	result, err := populateTemplate(templateString, data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func populateTemplate(templateString string, dataStruct interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, dataStruct)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getReferencedObjectsSection(
	changeSetDiff *commonTypes.ChangeSetDiff,
	symbolTable *symbol_table.SymbolTable,
) string {
	referencedResourcesString := ""
	referencedResources := changeSetDiff.NewObject.ParsedResource.ReferencedResources
	if len(referencedResources) > 0 {
		referencedResourcesString = "The following resources are referenced:\n"
		for _, referencedLogicalName := range referencedResources {
			referencedResource, _ := symbolTable.LookupResource(referencedLogicalName)
			referencedResourcesString += string(referencedLogicalName) + ": "
			referencedResourcesString += string(referencedResource.ResourceType) + "\n"
		}
	}

	referencedVariablesString := ""
	referencedVariables := changeSetDiff.NewObject.ParsedResource.ReferencedVariables
	if len(referencedVariables) > 0 {
		if referencedResourcesString != "" {
			referencedVariablesString += "\n"
		}
		referencedVariablesString += "The following variables are referenced:\n"
		for _, referencedVariableName := range referencedVariables {
			referencedVariable, _ := symbolTable.LookupVariable(referencedVariableName)
			referencedVariablesString += string(referencedVariableName) + ": "
			referencedVariablesString += string(referencedVariable.Type) + "\n"
		}
	}

	return referencedResourcesString + referencedVariablesString
}

func getVariableDetailsSection(
	changeSetDiff *commonTypes.ChangeSetDiff,
) string {
	variableDetailsString := ""
	variable := changeSetDiff.NewObject.ParsedVariable
	if variable.Default != "" {
		variableDetailsString += "Default value: " + variable.Default + "\n"
	}
	if variable.NaturalLanguage != "" {
		variableDetailsString += variable.NaturalLanguage
	}

	return variableDetailsString
}

//go:embed new update
var templatesFS embed.FS

func readTemplateFile(templatesDirectory string, fileName string) (string, error) {
	templateFilePath := filepath.Join(
		templatesDirectory,
		fileName,
	)

	data, err := fs.ReadFile(templatesFS, templateFilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
