// Package vm implements API functions residing under /provisioning/vm.
// This path contains methods for provision- and deprovisioning VMs.
package vm

const (
	pathPrefix = "/api/vsphere/v1/provisioning/vm.json"
	// DefaultCPUPerformanceType to be used if a VM definition is created
	// by NewDefinition.
	DefaultCPUPerformanceType = "performance"
	// DefaultDiskType to be used if a VM definition is created
	// by NewDefinition.
	DefaultDiskType = "ENT2"
)
