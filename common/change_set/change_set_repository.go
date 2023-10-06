package change_set

import (
	"salami/common/types"
)

type ChangeSetRepository struct {
	Diffs            []*types.ChangeSetDiff
	AddedObjects     []*types.Object
	resourceDiffs    map[types.LogicalName]*types.ChangeSetDiff
	variableDiffs    map[string]*types.ChangeSetDiff
	deletedResources map[types.LogicalName]*types.Object
	deletedVariables map[string]*types.Object
	changedResources map[types.LogicalName]*types.Object
	changedVariables map[string]*types.Object
}

func (ch *ChangeSetRepository) GetChangedObject(object *types.Object) *types.Object {
	if object.IsResource() {
		return ch.changedResources[object.ParsedResource.LogicalName]
	} else if object.IsVariable() {
		return ch.changedVariables[object.ParsedVariable.Name]
	}

	return nil
}

func (ch *ChangeSetRepository) WasObjectChanged(object *types.Object) bool {
	if object.IsResource() {
		_, ok := ch.changedResources[object.ParsedResource.LogicalName]
		return ok
	} else if object.IsVariable() {
		_, ok := ch.changedVariables[object.ParsedVariable.Name]
		return ok
	}

	return false
}

func (ch *ChangeSetRepository) WasObjectDeleted(object *types.Object) bool {
	if object.IsResource() {
		_, ok := ch.deletedResources[object.ParsedResource.LogicalName]
		return ok
	} else if object.IsVariable() {
		_, ok := ch.deletedVariables[object.ParsedVariable.Name]
		return ok
	}

	return false
}

func (ch *ChangeSetRepository) GetDiffForObject(object *types.Object) *types.ChangeSetDiff {
	if object.IsResource() {
		return ch.resourceDiffs[object.ParsedResource.LogicalName]
	} else if object.IsVariable() {
		return ch.variableDiffs[object.ParsedVariable.Name]
	}

	return nil
}

func NewChangeSetRepository(changeSet *types.ChangeSet) *ChangeSetRepository {
	addedObjects := make([]*types.Object, 0)
	deletedResources := make(map[types.LogicalName]*types.Object)
	deletedVariables := make(map[string]*types.Object)
	changedResources := make(map[types.LogicalName]*types.Object)
	changedVariables := make(map[string]*types.Object)
	resourceDiffs := make(map[types.LogicalName]*types.ChangeSetDiff)
	variableDiffs := make(map[string]*types.ChangeSetDiff)

	for _, diff := range changeSet.Diffs {
		addDiffToRepository(
			diff,
			&addedObjects,
			deletedResources,
			deletedVariables,
			changedResources,
			changedVariables,
			resourceDiffs,
			variableDiffs,
		)
	}

	return &ChangeSetRepository{
		Diffs:            changeSet.Diffs,
		AddedObjects:     addedObjects,
		deletedResources: deletedResources,
		deletedVariables: deletedVariables,
		changedResources: changedResources,
		changedVariables: changedVariables,
		resourceDiffs:    resourceDiffs,
		variableDiffs:    variableDiffs,
	}
}

func addDiffToRepository(
	diff *types.ChangeSetDiff,
	addedObjects *[]*types.Object,
	deletedResources map[types.LogicalName]*types.Object,
	deletedVariables map[string]*types.Object,
	changedResources map[types.LogicalName]*types.Object,
	changedVariables map[string]*types.Object,
	resourceDiffs map[types.LogicalName]*types.ChangeSetDiff,
	variableDiffs map[string]*types.ChangeSetDiff,
) {
	oldObject := diff.OldObject
	newObject := diff.NewObject
	diffType := diff.DiffType

	if diffType == types.DiffTypeRemove {
		if oldObject.IsResource() {
			deletedResources[oldObject.ParsedResource.LogicalName] = oldObject
			resourceDiffs[oldObject.ParsedResource.LogicalName] = diff
		} else if oldObject.IsVariable() {
			deletedVariables[oldObject.ParsedVariable.Name] = oldObject
			variableDiffs[oldObject.ParsedVariable.Name] = diff
		}
	} else if diffType == types.DiffTypeAdd {
		*addedObjects = append(*addedObjects, newObject)

		if newObject.IsResource() {
			resourceDiffs[newObject.ParsedResource.LogicalName] = diff
		} else if newObject.IsVariable() {
			variableDiffs[newObject.ParsedVariable.Name] = diff
		}
	} else {
		if newObject.IsResource() {
			changedResources[newObject.ParsedResource.LogicalName] = newObject
			resourceDiffs[newObject.ParsedResource.LogicalName] = diff
		} else if newObject.IsVariable() {
			changedVariables[newObject.ParsedVariable.Name] = newObject
			variableDiffs[newObject.ParsedVariable.Name] = diff
		}
	}
}
