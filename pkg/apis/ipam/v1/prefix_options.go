package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

const (
	createEmptyOptionKey            = "ipamv1/prefix/createEmpty"
	enableRouterRedundancyOptionKey = "ipamv1/prefix/enableRouterRedundancy"
	enableVMProvisioningKey         = "ipamv1/prefix/enableVMProvisioning"

	vlanProvisioningEnableKey      = "ipamv1/prefix/vlanProvisioning/enable"
	vlanProvisioningDescriptionKey = "ipamv1/prefix/vlanProvisioning/name"
)

type boolOpt struct {
	key   string
	value bool
}

func (bo boolOpt) ApplyToCreate(opts *types.CreateOptions) error {
	// Set only returns an error when the requested key is already set. Since we
	// overwrite the value, we ignore it deliberately.
	_ = opts.Set(bo.key, bo.value, true)
	return nil
}

// CreateEmpty can be used to define if a Prefix is to be created with all Address objects
// in it created (false) or only Address objects created that are actually in use (true).
func CreateEmpty(empty bool) types.CreateOption {
	return boolOpt{key: createEmptyOptionKey, value: empty}
}

// EnableVMProvisioning allows this prefix to be used for provisioning of Anexia Dynamic Compute instances.
func EnableVMProvisioning(v bool) types.CreateOption {
	return boolOpt{key: enableVMProvisioningKey, value: v}
}

// EnableRouterRedundancy allows you to redundancy features provided by the Engine.
func EnableRouterRedundancy(v bool) types.CreateOption {
	return boolOpt{key: enableRouterRedundancyOptionKey, value: v}
}

type genVLAN struct {
	enabled     bool
	description string
}

func (g genVLAN) ApplyToCreate(opts *types.CreateOptions) error {
	// Set only returns an error when the requested key is already set. Since we
	// overwrite the value, we ignore it deliberately.
	_ = opts.Set(vlanProvisioningEnableKey, g.enabled, true)
	_ = opts.Set(vlanProvisioningDescriptionKey, g.description, true)
	return nil
}

// GenerateVLAN allows you to generate a new VLAN for the given prefix with [description].
// The description is optional and can be left empty.
func GenerateVLAN(description string) types.CreateOption {
	return genVLAN{enabled: true, description: description}
}
