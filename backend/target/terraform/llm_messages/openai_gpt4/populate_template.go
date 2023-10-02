package openai_gpt4

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"strings"
)

func populateResourceTemplate(
	templatePath string,
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (string, error) {
	templateString, err := readTemplateFile(templatePath)
	if err != nil {
		return "", err
	}

	data := struct {
		ResourceType             string
		LogicalName              string
		NaturalLanguage          string
		ReferencedObjectsSection string
	}{
		ResourceType:             string(object.ParsedResource.ResourceType),
		LogicalName:              string(object.ParsedResource.LogicalName),
		NaturalLanguage:          object.ParsedResource.NaturalLanguage,
		ReferencedObjectsSection: getReferencedObjectsSection(object, symbolTable),
	}

	result, err := populateTemplate(templateString, data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func populateVariableTemplate(templatePath string, object *commonTypes.Object) (string, error) {
	templateString, err := readTemplateFile(templatePath)
	if err != nil {
		return "", err
	}

	data := struct {
		Name                   string
		VariableType           string
		VariableDetailsSection string
	}{
		Name:                   string(object.ParsedVariable.Name),
		VariableType:           string(object.ParsedVariable.Type),
		VariableDetailsSection: getVariableDetailsSection(object),
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

//go:embed templates
var templatesFS embed.FS

func readTemplateFile(filePath string) (string, error) {
	data, err := fs.ReadFile(templatesFS, filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getReferencedObjectsSection(object *commonTypes.Object, symbolTable *symbol_table.SymbolTable) string {
	referencedResourcesString := ""
	referencedResources := object.ParsedResource.ReferencedResources
	if len(referencedResources) > 0 {
		referencedResourcesString = "The following resources are referenced:\n"
		for _, referencedLogicalName := range referencedResources {
			referencedResource, _ := symbolTable.LookupResource(referencedLogicalName)
			referencedResourcesString += string(referencedLogicalName) + ": "
			referencedResourcesString += string(referencedResource.ResourceType) + "\n"
		}
	}

	referencedVariablesString := ""
	referencedVariables := object.ParsedResource.ReferencedVariables
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

func getVariableDetailsSection(object *commonTypes.Object) string {
	variableDetailsString := ""
	variable := object.ParsedVariable
	if variable.Default != "" {
		variableDetailsString += "Default value: " + variable.Default + "\n"
	}
	if variable.NaturalLanguage != "" {
		variableDetailsString += variable.NaturalLanguage
	}

	return variableDetailsString
}
