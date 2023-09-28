package target

import (
	"salami/backend/target/terraform"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

var targetFuncsMap = map[string]backendTypes.NewTargetFunc{
	commonTypes.TerraformPlatform: terraform.NewTarget,
}

func ResolveTarget(targetConfig commonTypes.TargetConfig) backendTypes.Target {
	var target backendTypes.Target
	if targetConfig.Platform == commonTypes.TerraformPlatform {
		target = targetFuncsMap[commonTypes.TerraformPlatform]()
	}
	return target
}
