package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/apis/internal/gs"
)

// anxcloud:object

// Cluster represents a Kubernetes cluster
// This resource does not support updates
type Cluster struct {
	gs.GenericService
	gs.HasState

	// Identifier of the cluster
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	// Name of the cluster. Must be an RFC 1123 hostname in lowercase
	Name string `json:"name,omitempty"`
	// Kubernetes version to be used for the cluster. We recommend to use the default value.
	Version string `json:"version,omitempty"`
	// Location where the cluster will be deployed
	Location corev1.Location `json:"location,omitempty"`
	// If set to true, Service VMs providing load balancers and outbound masquerade are created for this cluster.
	// Default: true. Optional value can be set via pkg/utils/pointer.Bool
	NeedsServiceVMs *bool `json:"needs_service_vms,omitempty"`
	// If enabled, Service VMs are configured as NAT gateways connecting the internal cluster network to the internet.
	// Requires Service VMs.
	EnableNATGateways *bool `json:"enable_nat_gateways,omitempty"`
	// If enabled, Service VMs are set up as LBaaS hosts enabling K8s services of type LoadBalancer.
	// Requires Service VMs.
	EnableLBaaS *bool `json:"enable_lbaas,omitempty"`

	// Contains a kubeconfig if available
	KubeConfig string `json:"kubeconfig,omitempty"`
}

// AwaitCompletion blocks until the Cluster state is "OK"
func (c *Cluster) AwaitCompletion(ctx context.Context, a api.API) error {
	return gs.AwaitCompletion(ctx, a, c)
}
