package terraform

import (
	"os/exec"
	"salami/common/errors"
)

func (t *Terraform) VerifyPeerDependencies() error {
	if _, err := exec.LookPath("terraform"); err != nil {
		return &errors.ConfigError{Message: "'terraform' could not be found in your PATH"}
	}
	return nil
}
