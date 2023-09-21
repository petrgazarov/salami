package change_manager

import "salami/common/types"

func ComputeNewObjects(
	previousResources map[types.LogicalName]*types.Object,
	previousVariables map[string]*types.Object,
	changeSet *types.ChangeSet,
) []*types.Object {
	return nil
}
