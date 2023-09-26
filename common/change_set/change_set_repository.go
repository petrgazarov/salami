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

		if oldObject != nil && oldObject.IsResource() {
			if newObject == nil {
				deletedResources[oldObject.ParsedResource.LogicalName] = oldObject
			} else {
				changedResources[newObject.ParsedResource.LogicalName] = newObject
			}
		} else if oldObject != nil && oldObject.IsVariable() {
			if newObject == nil {
				deletedVariables[oldObject.ParsedVariable.Name] = oldObject
			} else {
				changedVariables[newObject.ParsedVariable.Name] = newObject
			}
		}

		if oldObject == nil {
			addedObjects = append(addedObjects, newObject)
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
