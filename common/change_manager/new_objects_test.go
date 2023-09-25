package change_manager_test

import (
	"salami/common/change_manager"
	"salami/common/types"
	"salami/common/utils/object_utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeNewObjects(t *testing.T) {
	t.Run("should return all new objects", func(t *testing.T) {
		previousObjects := getObjects("testdata/new_objects_test/previous_objects.json")
		previousResourcesMap, previousVariablesMap := object_utils.GetObjectMaps(previousObjects)
		changeSetDiffs := getChangeSetDiffs("testdata/new_objects_test/change_set_diffs.json")
		changeSet := &types.ChangeSet{
			Diffs: changeSetDiffs,
		}
		expectedNewObjects := getObjects("testdata/new_objects_test/new_objects.json")
		actualNewObjects := change_manager.ComputeNewObjects(previousResourcesMap, previousVariablesMap, changeSet)
		require.Equal(t, len(expectedNewObjects), len(actualNewObjects))
		for i, object := range actualNewObjects {
			require.Equal(t, expectedNewObjects[i], object)
		}
	})
}
