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
	Type            string
	Default         string
	NaturalLanguage string
}

type FixResourceValidationIssueDataStruct struct {
	ErrorMessages     []string
	ReferencedObjects []*struct{ TargetCode string }
	LogicalName       string
}

type FixVariableValidationIssueDataStruct struct {
	ErrorMessages []string
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
			Type:            string(object.ParsedVariable.Type),
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
		referencedObjects[i] = &struct{ TargetCode string }{}
		referencedObjects[i].TargetCode = referencedObject.TargetCode
	}

	switch templatePath {
	case templates.FixResourceValidationIssueTemplatePath:
		dataStruct = FixResourceValidationIssueDataStruct{
			ErrorMessages:     validationResult.ErrorMessages,
			ReferencedObjects: referencedObjects,
			LogicalName:       string(validationResult.ValidatedObject.ParsedResource.LogicalName),
		}
	case templates.FixVariableValidationIssueTemplatePath:
		dataStruct = FixVariableValidationIssueDataStruct{
			ErrorMessages: validationResult.ErrorMessages,
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

	for i, logicalName := range object.ParsedResource.ReferencedResources {
		referencedResource, _ := symbolTable.LookupResource(logicalName)
		referencedResources[i] = referencedResource
	}

	for i, variableName := range object.ParsedResource.ReferencedVariables {
		referencedVariable, _ := symbolTable.LookupVariable(variableName)
		referencedVariables[i] = referencedVariable
	}

	return referencedResources, referencedVariables
}
