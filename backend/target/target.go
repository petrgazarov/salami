package target

import (
	"salami/backend/target/terraform"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

func ResolveTarget(targetConfig commonTypes.TargetConfig) backendTypes.Target {
	var target backendTypes.Target
	if targetConfig.Platform == commonTypes.TerraformPlatform {
		target = terraform.NewTarget()
	}
	return target
}
