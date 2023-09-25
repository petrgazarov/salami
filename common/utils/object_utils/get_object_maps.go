package object_utils

import "salami/common/types"

func GetObjectMaps(
	objects []*types.Object,
) (map[types.LogicalName]*types.Object, map[string]*types.Object) {
	resourcesMap := make(map[types.LogicalName]*types.Object)
	variablesMap := make(map[string]*types.Object)

	for _, object := range objects {
		if object.IsResource() {
			resourcesMap[object.ParsedResource.LogicalName] = object
		} else if object.IsVariable() {
			variablesMap[object.ParsedVariable.Name] = object
		}
	}

	return resourcesMap, variablesMap
}
