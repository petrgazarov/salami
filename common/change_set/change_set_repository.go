package change_set

import "salami/common/types"

type ChangeSetRepository struct {
	AddedObjects     []*types.Object
	deletedResources map[types.LogicalName]*types.Object
	deletedVariables map[string]*types.Object
	changedResources map[types.LogicalName]*types.Object
	changedVariables map[string]*types.Object
}

func (ch *ChangeSetRepository) GetChanged(object *types.Object) *types.Object {
	if object.IsResource() {
		return ch.changedResources[object.ParsedResource.LogicalName]
	} else if object.IsVariable() {
		return ch.changedVariables[object.ParsedVariable.Name]
	}
	return nil
}

func (ch *ChangeSetRepository) WasChanged(object *types.Object) bool {
	if object.IsResource() {
		_, ok := ch.changedResources[object.ParsedResource.LogicalName]
		return ok
	} else if object.IsVariable() {
		_, ok := ch.changedVariables[object.ParsedVariable.Name]
		return ok
	}
	return false
}

func (ch *ChangeSetRepository) WasDeleted(object *types.Object) bool {
	if object.IsResource() {
		_, ok := ch.deletedResources[object.ParsedResource.LogicalName]
		return ok
	} else if object.IsVariable() {
		_, ok := ch.deletedVariables[object.ParsedVariable.Name]
		return ok
	}
	return false
}

func NewChangeSetRepository(changeSet *types.ChangeSet) *ChangeSetRepository {
	deletedResources := make(map[types.LogicalName]*types.Object)
	deletedVariables := make(map[string]*types.Object)
	changedResources := make(map[types.LogicalName]*types.Object)
	changedVariables := make(map[string]*types.Object)
	addedObjects := make([]*types.Object, 0)

	for _, diff := range changeSet.Diffs {
		oldObject := diff.OldObject
		newObject := diff.NewObject
		diffType := diff.DiffType

		if diffType == types.DiffTypeRemove {
			if oldObject.IsResource() {
				deletedResources[oldObject.ParsedResource.LogicalName] = oldObject
			} else if oldObject.IsVariable() {
				deletedVariables[oldObject.ParsedVariable.Name] = oldObject
			}
		} else if diffType == types.DiffTypeAdd {
			addedObjects = append(addedObjects, newObject)
		} else {
			if newObject.IsResource() {
				changedResources[newObject.ParsedResource.LogicalName] = newObject
			} else if newObject.IsVariable() {
				changedVariables[newObject.ParsedVariable.Name] = newObject
			}
		}
	}

	return &ChangeSetRepository{
		AddedObjects:     addedObjects,
		deletedResources: deletedResources,
		deletedVariables: deletedVariables,
		changedResources: changedResources,
		changedVariables: changedVariables,
	}
}
