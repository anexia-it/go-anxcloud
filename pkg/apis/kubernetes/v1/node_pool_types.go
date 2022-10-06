package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/apis/internal/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// NodePool represents a Kubernetes node pool
// This resource does not support updates
type NodePool struct {
	gs.GenericService
	HasState

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
	OperatingSystem string `json:"operating_system,omitempty"`
}

// AwaitCompletion blocks until the NodePool state is "OK"
func (np *NodePool) AwaitCompletion(ctx context.Context, a api.API) error {
	return gs.AwaitCompletion(ctx, a, np)
}
