package change_manager

import "salami/common/types"

func ComputeNewObjects(
	previousResources map[types.LogicalName]*types.Object,
	previousVariables map[string]*types.Object,
	changeSet *types.ChangeSet,
) []*types.Object {
	deletedResources,
		deletedVariables,
		addedResources,
		addedVariables,
		changedResources,
		changedVariables := getObjectMapsPerAction(previousResources, previousVariables, changeSet)

	objects := make([]*types.Object, 0)
	for _, object := range previousResources {
		if _, ok := deletedResources[object.ParsedResource.LogicalName]; ok {
			continue
		}
		if _, ok := changedResources[object.ParsedResource.LogicalName]; ok {
			objects = append(objects, changedResources[object.ParsedResource.LogicalName])
		} else {
			objects = append(objects, object)
		}
	}
	for _, object := range previousVariables {
		if _, ok := deletedVariables[object.ParsedVariable.Name]; ok {
			continue
		}
		if _, ok := changedVariables[object.ParsedVariable.Name]; ok {
			objects = append(objects, changedVariables[object.ParsedVariable.Name])
		} else {
			objects = append(objects, object)
		}
	}
	for _, object := range addedResources {
		objects = append(objects, object)
	}
	for _, object := range addedVariables {
		objects = append(objects, object)
	}

	return objects
}

func getObjectMapsPerAction(
	previousResources map[types.LogicalName]*types.Object,
	previousVariables map[string]*types.Object,
	changeSet *types.ChangeSet,
) (
	map[types.LogicalName]*types.Object,
	map[string]*types.Object,
	map[types.LogicalName]*types.Object,
	map[string]*types.Object,
	map[types.LogicalName]*types.Object,
	map[string]*types.Object,
) {
	deletedResources := make(map[types.LogicalName]*types.Object)
	deletedVariables := make(map[string]*types.Object)
	addedResources := make(map[types.LogicalName]*types.Object)
	addedVariables := make(map[string]*types.Object)
	changedResources := make(map[types.LogicalName]*types.Object)
	changedVariables := make(map[string]*types.Object)

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
			if newObject.IsResource() {
				addedResources[newObject.ParsedResource.LogicalName] = newObject
			} else if newObject.IsVariable() {
				addedVariables[newObject.ParsedVariable.Name] = newObject
			}
		}
	}

	return deletedResources, deletedVariables, addedResources, addedVariables, changedResources, changedVariables
}
