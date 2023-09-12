package terraform

import "salami/compiler/types"

type TerraformTarget struct {
}

func (t *TerraformTarget) GenerateCodeFiles(
	objectsMap map[string][]types.Object,
) ([]*types.DestinationFile, []error) {
	return []*types.DestinationFile{}, []error{}
}
