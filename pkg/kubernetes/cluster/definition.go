package cluster

import "go.anx.io/go-anxcloud/pkg/apis/common/gs"

type Location string

const (
	LocationANX04 = Location("52b5f6b2fd3a4a7eaaedf1a7c019e9ea")
	LocationANX25 = Location("97b3cdf37368470496a249aff0b84845")
)

type IntBoolean int

const (
	Disable IntBoolean = 0
	Enable  IntBoolean = 1
)

var (
	StateNewlyCreated = gs.State{ID: "2", Text: "NewlyCreated", Type: gs.StateTypePending}
)

type Definition struct {
	Name                           string     `json:"name,omitempty"`
	State                          gs.State   `json:"state,omitempty"`
	Location                       Location   `json:"location,omitempty"`
	Version                        string     `json:"version,omitempty"`
	Autoscaling                    IntBoolean `json:"autoscaling,omitempty"`
	CNIPlugin                      string     `json:"cni_plugin,omitempty"`
	APIServerAllowList             string     `json:"api_server_allow_list,omitempty"`
	MaintenanceWindowStart         string     `json:"maintenance_window_start,omitempty"`
	MaintenanceWindowLength        string     `json:"maintenance_window_length,omitempty"`
	ManageInternalIPv4Prefix       IntBoolean `json:"manage_internal_ipv4_prefix,omitempty"`
	InternalIPv4Prefix             string     `json:"internal_ipv4_prefix,omitempty"`
	NeedsServiceVMs                IntBoolean `json:"needs_service_vms,omitempty"`
	EnableNATGateways              IntBoolean `json:"enable_nat_gateways,omitempty"`
	EnableLBaaS                    IntBoolean `json:"enable_lbaas,omitempty"`
	ExternalIPFamilies             string     `json:"external_ip_families,omitempty"`
	ManageExternalIPv4Prefix       IntBoolean `json:"manage_external_ipv4_prefix,omitempty"`
	ExternalIPv4Prefix             string     `json:"external_ipv4_prefix,omitempty"`
	ManageExternalIPv6Prefix       IntBoolean `json:"manage_external_ipv6_prefix,omitempty"`
	ExternalIPv6Prefix             string     `json:"external_ipv6_prefix,omitempty"`
	ServiceVM01InternalIPv4Address string     `json:"service_vm_01_internal_ipv4_address,omitempty"`
	ServiceVM02InternalIPv4Address string     `json:"service_vm_02_internal_ipv4_address,omitempty"`
	ServiceVM01ExternalIPv4Address string     `json:"service_vm_01_external_ipv4_address,omitempty"`
	ServiceVM02ExternalIPv4Address string     `json:"service_vm_02_external_ipv4_address,omitempty"`
	ServiceVM01ExternalIPv6Address string     `json:"service_vm_01_external_ipv6_address,omitempty"`
	ServiceVM02ExternalIPv6Address string     `json:"service_vm_02_external_ipv6_address,omitempty"`
	ServiceUser                    string     `json:"service_user,omitempty"`
	KKPProjectID                   string     `json:"kp_project_id,omitempty"`
	KKPClusterID                   string     `json:"kp_cluster_id,omitempty"`
}
