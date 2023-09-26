package change_set

import (
	"reflect"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/common/utils/object_utils"
)

func NewChangeSet(
	previousObjects []*types.Object,
	symbolTable *symbol_table.SymbolTable,
) *types.ChangeSet {
	previousResourcesMap, previousVariablesMap := object_utils.GetObjectMaps(previousObjects)
	seenResources := make(map[types.LogicalName]bool)
	seenVariables := make(map[string]bool)
	changeSet := &types.ChangeSet{
		Diffs: make([]types.ChangeSetDiff, 0),
	}

	recordResourceChangesOrAdditions(symbolTable, changeSet, previousResourcesMap, seenResources)
	recordVariableChangesOrAdditions(symbolTable, changeSet, previousVariablesMap, seenVariables)
	recordResourceDeletions(changeSet, previousResourcesMap, seenResources)
	recordVariableDeletions(changeSet, previousVariablesMap, seenVariables)

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
			ParsedResource: resource,
		}
		if _, ok := previousResources[logicalName]; !ok {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: nil,
				NewObject: object,
			})
		} else if !reflect.DeepEqual(previousResources[logicalName].ParsedResource, object.ParsedResource) {
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
			ParsedVariable: variable,
		}
		if _, ok := previousVariables[name]; !ok {
			changeSet.Diffs = append(changeSet.Diffs, types.ChangeSetDiff{
				OldObject: nil,
				NewObject: object,
			})
		} else if !reflect.DeepEqual(previousVariables[name].ParsedVariable, object.ParsedVariable) {
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
