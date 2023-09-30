package lock_file_manager

import "salami/common/types"

func lockFileToCommonResource(lockFileObject Object) *types.ParsedResource {
	referencedResources := make([]types.LogicalName, len(lockFileObject.ParsedResource.ReferencedResources))
	for j, use := range lockFileObject.ParsedResource.ReferencedResources {
		referencedResources[j] = types.LogicalName(use)
	}
	return &types.ParsedResource{
		ResourceType:        types.ResourceType(lockFileObject.ParsedResource.ResourceType),
		LogicalName:         types.LogicalName(lockFileObject.ParsedResource.LogicalName),
		NaturalLanguage:     lockFileObject.ParsedResource.NaturalLanguage,
		ReferencedResources: referencedResources,
		ReferencedVariables: lockFileObject.ParsedResource.ReferencedVariables,
		SourceFilePath:      lockFileObject.ParsedResource.SourceFilePath,
		SourceFileLine:      lockFileObject.ParsedResource.SourceFileLine,
	}
}

func lockFileToCommonVariable(lockFileObject Object) *types.ParsedVariable {
	return &types.ParsedVariable{
		Name:            lockFileObject.ParsedVariable.Name,
		NaturalLanguage: lockFileObject.ParsedVariable.NaturalLanguage,
		Default:         lockFileObject.ParsedVariable.Default,
		Type:            types.VariableType(lockFileObject.ParsedVariable.VariableType),
		SourceFilePath:  lockFileObject.ParsedVariable.SourceFilePath,
		SourceFileLine:  lockFileObject.ParsedVariable.SourceFileLine,
	}
}

func commonToLockFileResource(commonResource *types.ParsedResource) *ParsedResource {
	if commonResource == nil {
		return nil
	}
	parsedResource := &ParsedResource{
		ResourceType:        string(commonResource.ResourceType),
		LogicalName:         string(commonResource.LogicalName),
		NaturalLanguage:     commonResource.NaturalLanguage,
		ReferencedResources: make([]string, len(commonResource.ReferencedResources)),
		ReferencedVariables: commonResource.ReferencedVariables,
		SourceFilePath:      commonResource.SourceFilePath,
		SourceFileLine:      commonResource.SourceFileLine,
	}
	for j, referencedResource := range commonResource.ReferencedResources {
		parsedResource.ReferencedResources[j] = string(referencedResource)
	}

	return parsedResource
}

func commonToLockFileVariable(commonVariable *types.ParsedVariable) *ParsedVariable {
	if commonVariable == nil {
		return nil
	}
	parsedVariable := &ParsedVariable{
		Name:            commonVariable.Name,
		NaturalLanguage: commonVariable.NaturalLanguage,
		VariableType:    string(commonVariable.Type),
		Default:         commonVariable.Default,
		SourceFilePath:  commonVariable.SourceFilePath,
		SourceFileLine:  commonVariable.SourceFileLine,
	}

	return parsedVariable
}
