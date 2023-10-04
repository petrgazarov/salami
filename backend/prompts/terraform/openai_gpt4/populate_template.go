package openai_gpt4

import (
	"salami/backend/prompts/terraform/openai_gpt4/templates"
	backendTypes "salami/backend/types"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
)

type CreateUpdateResourceDataStruct struct {
	ResourceType        string
	LogicalName         string
	NaturalLanguage     string
	ReferencedResources []*commonTypes.ParsedResource
	ReferencedVariables []*commonTypes.ParsedVariable
}

type CreateUpdateVariableDataStruct struct {
	Name            string
	VariableType    string
	Default         string
	NaturalLanguage string
}

type FixResourceValidationIssueDataStruct struct {
	ErrorMessage      string
	ReferencedObjects []*struct{ TargetCode string }
	LogicalName       string
}

type FixVariableValidationIssueDataStruct struct {
	ErrorMessage string
}

func populateCreateUpdateTemplate(
	templatePath string,
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) (string, error) {
	var dataStruct interface{}

	switch templatePath {
	case templates.CreateResourceTemplatePath, templates.UpdateResourceTemplatePath:
		referencedResources, referencedVariables := getCreateUpdateTemplateReferencedObjects(object, symbolTable)

		dataStruct = CreateUpdateResourceDataStruct{
			ResourceType:        string(object.ParsedResource.ResourceType),
			LogicalName:         string(object.ParsedResource.LogicalName),
			NaturalLanguage:     object.ParsedResource.NaturalLanguage,
			ReferencedResources: referencedResources,
			ReferencedVariables: referencedVariables,
		}
	case templates.CreateVariableTemplatePath, templates.UpdateVariableTemplatePath:
		dataStruct = CreateUpdateVariableDataStruct{
			Name:            string(object.ParsedVariable.Name),
			VariableType:    string(object.ParsedVariable.Type),
			Default:         object.ParsedVariable.Default,
			NaturalLanguage: object.ParsedVariable.NaturalLanguage,
		}
	}

	result, err := templates.PopulateTemplate(templatePath, dataStruct)
	if err != nil {
		return "", err
	}

	return result, nil
}

func populateValidationTemplate(
	templatePath string,
	validationResult *backendTypes.CodeValidationResult,
) (string, error) {
	var dataStruct interface{}

	referencedObjects := make([]*struct{ TargetCode string }, len(validationResult.ReferencedObjects))
	for i, referencedObject := range validationResult.ReferencedObjects {
		referencedObjects[i].TargetCode = referencedObject.TargetCode
	}

	switch templatePath {
	case templates.FixResourceValidationIssueTemplatePath:
		dataStruct = FixResourceValidationIssueDataStruct{
			ErrorMessage:      validationResult.ErrorMessage,
			ReferencedObjects: referencedObjects,
			LogicalName:       string(validationResult.ValidatedObject.ParsedResource.LogicalName),
		}
	case templates.FixVariableValidationIssueTemplatePath:
		dataStruct = FixVariableValidationIssueDataStruct{
			ErrorMessage: validationResult.ErrorMessage,
		}
	}

	result, err := templates.PopulateTemplate(templatePath, dataStruct)
	if err != nil {
		return "", err
	}

	return result, nil
}

func getCreateUpdateTemplateReferencedObjects(
	object *commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) ([]*commonTypes.ParsedResource, []*commonTypes.ParsedVariable) {
	referencedResources := make([]*commonTypes.ParsedResource, len(object.ParsedResource.ReferencedResources))
	referencedVariables := make([]*commonTypes.ParsedVariable, len(object.ParsedResource.ReferencedVariables))

	for _, logicalName := range object.ParsedResource.ReferencedResources {
		referencedResource, _ := symbolTable.LookupResource(logicalName)
		referencedResources = append(referencedResources, referencedResource)
	}

	for _, variableName := range object.ParsedResource.ReferencedVariables {
		referencedVariable, _ := symbolTable.LookupVariable(variableName)
		referencedVariables = append(referencedVariables, referencedVariable)
	}

	return referencedResources, referencedVariables
}
