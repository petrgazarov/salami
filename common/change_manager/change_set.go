package change_manager

import (
	"reflect"
	"salami/common/symbol_table"
	"salami/common/types"
)

func GenerateChangeSet(
	previousResources map[types.LogicalName]*types.Object,
	previousVariables map[string]*types.Object,
	symbolTable *symbol_table.SymbolTable,
) *types.ChangeSet {
	seenResources := make(map[types.LogicalName]bool)
	seenVariables := make(map[string]bool)
	changeSet := &types.ChangeSet{
		Diffs: make([]types.ChangeSetDiff, 0),
	}

	recordResourceChangesOrAdditions(symbolTable, changeSet, previousResources, seenResources)
	recordVariableChangesOrAdditions(symbolTable, changeSet, previousVariables, seenVariables)
	recordResourceDeletions(changeSet, previousResources, seenResources)
	recordVariableDeletions(changeSet, previousVariables, seenVariables)

	return changeSet
}

func recordResourceChangesOrAdditions(
	symbolTable *symbol_table.SymbolTable,
	changeSet *types.ChangeSet,
	previousResources map[types.LogicalName]*types.Object,
	seenResources map[types.LogicalName]bool,
) {
	for logicalName, resource := range symbolTable.ResourceTable {
		object := &types.Object{
			SourceFilePath: resource.SourceFilePath,
			Parsed:         resource,
		}
		if _, ok := previousResources[logicalName]; !ok {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: nil,
				NewObject: object,
			})
		} else if !reflect.DeepEqual(previousResources[logicalName].Parsed, object.Parsed) {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: previousResources[logicalName],
				NewObject: object,
			})
		}
		seenResources[logicalName] = true
	}
}

func recordVariableChangesOrAdditions(
	symbolTable *symbol_table.SymbolTable,
	changeSet *types.ChangeSet,
	previousVariables map[string]*types.Object,
	seenVariables map[string]bool,
) {
	for name, variable := range symbolTable.VariableTable {
		object := &types.Object{
			SourceFilePath: variable.SourceFilePath,
			Parsed:         variable,
		}
		if _, ok := previousVariables[name]; !ok {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: nil,
				NewObject: object,
			})
		} else if !reflect.DeepEqual(previousVariables[name].Parsed, object.Parsed) {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: previousVariables[name],
				NewObject: object,
			})
		}
		seenVariables[name] = true
	}
}

func recordResourceDeletions(
	changeSet *types.ChangeSet,
	previousResources map[types.LogicalName]*types.Object,
	seenResources map[types.LogicalName]bool,
) {
	for logicalName, object := range previousResources {
		if !seenResources[logicalName] {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: object,
				NewObject: nil,
			})
		}
	}
}

func recordVariableDeletions(
	changeSet *types.ChangeSet,
	previousVariables map[string]*types.Object,
	seenVariables map[string]bool,
) {
	for name, object := range previousVariables {
		if !seenVariables[name] {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: object,
				NewObject: nil,
			})
		}
	}
}
