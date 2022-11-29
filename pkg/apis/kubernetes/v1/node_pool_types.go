package v1

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// OperatingSystem is a typed string for available OS templates
type OperatingSystem string

const (
	// FlatcarLinux is currently the only available OS template available for Kubernetes Node Pools
	FlatcarLinux OperatingSystem = "Flatcar Linux"
)

// anxcloud:object:hooks=RequestBodyHook

// NodePool represents a Kubernetes node pool
// This resource does not support updates
type NodePool struct {
	gs.GenericService
	gs.HasState

	// Identifier of the node pool
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	// Name of the node pool. Must be an RFC 1123 hostname in lowercase
	Name string `json:"name,omitempty"`
	// Cluster in which the node pool is deployed
	Cluster Cluster `json:"cluster,omitempty" anxcloud:"filterable"`

	// Number of replicas. Can be changed via machine controller.
	// Default: 3 (see FAQ for more details) Optional value can be set via pkg/utils/pointer.Int
	Replicas *int `json:"replicas,omitempty"`
	// Number of computation cores for each node. The provided cores will be "performance" type CPUs. Must be at least 1 and no more than 16
	CPUs int `json:"cpus,omitempty"`
	// RAM size for each node in bytes. Must be a multiple of 1 GiB, at least 2 GiB and no more than 64 GiB
	Memory int `json:"memory,omitempty"`
	// Size of the disk for each node in bytes. Its performance type will be the default Anexia Engine provides for the given location.
	// Must be a multiple of 1 GiB, at least 10 GiB and no more than 1600 GiB
	DiskSize int `json:"disk_size,omitempty"`

	// Operating system for deployment on the nodes. Default: Flatcar Linux
	OperatingSystem OperatingSystem `json:"operating_system,omitempty"`
}
