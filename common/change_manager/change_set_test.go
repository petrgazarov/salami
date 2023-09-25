package change_manager_test

import (
	"encoding/json"
	"io"
	"os"
	"salami/common/change_manager"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/common/utils"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateChangeSet(t *testing.T) {
	t.Run("should return an empty change set when there are no changes", func(t *testing.T) {
		previousObjects := getObjects("testdata/change_set_test/previous_objects.json")
		previousResourcesMap, previousVariablesMap := utils.GetObjectMaps(previousObjects)
		newResources, newVariables := utils.ObjectsToParsedObjects(previousObjects)
		symbolTable, err := symbol_table.NewSymbolTable(newResources, newVariables)
		require.NoError(t, err)
		changeSet := change_manager.GenerateChangeSet(previousResourcesMap, previousVariablesMap, symbolTable)
		require.Equal(t, changeSet, &types.ChangeSet{Diffs: []types.ChangeSetDiff{}})
	})

	t.Run("should return a change set with additions, deletions, and changes when they exist", func(t *testing.T) {
		previousObjects := getObjects("testdata/change_set_test/previous_objects.json")
		previousResourcesMap, previousVariablesMap := utils.GetObjectMaps(previousObjects)

		newObjects := getObjects("testdata/change_set_test/new_objects.json")
		newResources, newVariables := utils.ObjectsToParsedObjects(newObjects)
		symbolTable, err := symbol_table.NewSymbolTable(newResources, newVariables)

		require.NoError(t, err)
		changeSet := change_manager.GenerateChangeSet(previousResourcesMap, previousVariablesMap, symbolTable)
		changeSetDiffs := sortChangeSetDiffs(changeSet.Diffs)
		require.Equal(t, 4, len(changeSetDiffs))
		expectedDiffs := getChangeSetDiffs("testdata/change_set_test/change_set_diffs.json")
		for i, actualDiff := range changeSetDiffs {
			require.Equal(t, expectedDiffs[i], actualDiff)
		}
	})
}

func getObjects(filePaths string) []*types.Object {
	jsonFile, err := os.Open(filePaths)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var objects []*types.Object
	json.Unmarshal(byteValue, &objects)
	return objects
}

func getChangeSetDiffs(filePath string) []types.ChangeSetDiff {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var changeSetDiffs []types.ChangeSetDiff
	json.Unmarshal(byteValue, &changeSetDiffs)
	return changeSetDiffs
}

func sortChangeSetDiffs(diffs []types.ChangeSetDiff) []types.ChangeSetDiff {
	getNameAndType := func(obj *types.Object) (string, bool) {
		if obj == nil {
			return "", false
		}
		isVar := obj.IsVariable()
		name := ""
		if isVar {
			name = obj.ParsedVariable.Name
		} else {
			name = string(obj.ParsedResource.LogicalName)
		}
		return name, isVar
	}

	sort.Slice(diffs, func(i, j int) bool {
		iOldName, iOldIsVar := getNameAndType(diffs[i].OldObject)
		iNewName, iNewIsVar := getNameAndType(diffs[i].NewObject)
		jOldName, jOldIsVar := getNameAndType(diffs[j].OldObject)
		jNewName, jNewIsVar := getNameAndType(diffs[j].NewObject)

		if iOldIsVar != jOldIsVar {
			return iOldIsVar
		}

		if iNewIsVar != jNewIsVar {
			return iNewIsVar
		}

		if iOldName != jOldName {
			return iOldName < jOldName
		}

		return iNewName < jNewName
	})

	return diffs
}
