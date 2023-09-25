package object_utils

import "salami/common/types"

func ObjectsToParsedObjects(
	objects []*types.Object,
) ([]*types.ParsedResource, []*types.ParsedVariable) {
	resources := make([]*types.ParsedResource, 0)
	variables := make([]*types.ParsedVariable, 0)
	for _, object := range objects {
		if object.IsResource() {
			resources = append(resources, object.ParsedResource)
		} else if object.IsVariable() {
			variables = append(variables, object.ParsedVariable)
		}
	}

	return resources, variables
}
