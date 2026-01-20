package v1

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
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

	// Identifier of an internal v4 prefix (to be) assigned to the cluster. If ManageInternalIPv4Prefix
	// is set to false, the Prefix given in this field is used when creating the cluster, otherwise a new
	// prefix will be created automatically. The API will always return the Prefix for the Cluster,
	// when ManageInternalIPv4Prefix is true, this will be the Prefix that was created automatically.
	InternalIPv4Prefix *common.PartialResource `json:"internal_ipv4_prefix,omitempty"`
	// Identifier of an external v4 prefix (to be) assigned to the cluster. If ManageExternalIPv4Prefix
	// is set to false, the Prefix given in this field is used when creating the cluster, otherwise a new
	// prefix will be created automatically. The API will always return the Prefix for the Cluster,
	// when ManageExternalIPv4Prefix is true, this will be the Prefix that was created automatically.
	ExternalIPv4Prefix *common.PartialResource `json:"external_ipv4_prefix,omitempty"`
	// Identifier of an external v6 prefix (to be) assigned to the cluster. If ManageExternalIPv6Prefix
	// is set to false, the Prefix given in this field is used when creating the cluster, otherwise a new
	// prefix will be created automatically. The API will always return the Prefix for the Cluster,
	// when ManageExternalIPv6Prefix is true, this will be the Prefix that was created automatically.
	ExternalIPv6Prefix *common.PartialResource `json:"external_ipv6_prefix,omitempty"`

	// If set to true an internal v4 prefix is automatically created for the cluster. Defaults to true if not set.
	ManageInternalIPv4Prefix *bool `json:"manage_internal_ipv4_prefix,omitempty"`
	// If set to true an external v4 prefix is automatically created for the cluster. Defaults to true if not set.
	ManageExternalIPv4Prefix *bool `json:"manage_external_ipv4_prefix,omitempty"`
	// If set to true an external v6 prefix is automatically created for the cluster. Defaults to true if not set.
	ManageExternalIPv6Prefix *bool `json:"manage_external_ipv6_prefix,omitempty"`

	ExternalVlan *common.PartialResource `json:"external_vlan,omitempty"`
	InternalVlan *common.PartialResource `json:"internal_vlan,omitempty"`

	KkpProjectID *string `json:"kkp_project_id,omitempty"`
	KkpClusterID *string `json:"kkp_cluster_id,omitempty"`

	// Contains a kubeconfig if available
	KubeConfig *string `json:"kubeconfig,omitempty"`

	// Enable autoscaling for this cluster. You will need to explicitly configure
	// your node pools for autoscaling, please refer to the provided [Autoscaling documentation]
	// for details. Defaults to false if unset.
	//
	// [Autoscaling documentation]: https://engine.anexia-it.com/docs/en/module/kubernetes/user-guide/autoscaling
	EnableAutoscaling *bool `json:"autoscaling,omitempty"`

	// Container Network Interface plugin to be installed on the cluster.
	// Only the default value is supported. Default: canal = Canal
	CniPlugin string `json:"cni_plugin,omitempty"`

	// Space-separated list of IP networks in CIDR notation, which are allowed to access the cluster's API server.
	// If left empty, there will be no IP address-based restrictions.
	ApiServerAllowlist string `json:"apiserver_allowlist,omitempty"`

	// IP families to use for external networking.
	// Addresses from the selected protocols will be allocated for the cluster's Service VMs and load balancing,
	// if those are enabled. One of: IPv4 = IPv4, DualStack = IPv4 & IPv6. Default: DualStack
	ExternalIPFamilies string `json:"external_ip_families,omitempty"`
}
