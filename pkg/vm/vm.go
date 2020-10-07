// Package vm contains
package vm

import (
	"fmt"
)

// ProvisioningError wraps errors that the anxcloud API may return when the provisioning of a VM is requested.
type ProvisioningError struct {
	Errors []string
}

func (p ProvisioningError) Error() string {
	return fmt.Sprintf("could not provision VM: %q", p.Errors)
}
