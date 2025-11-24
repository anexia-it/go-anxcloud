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

type Definition struct {
	Name                     string     `json:"name"`
	State                    gs.State   `json:"state"`
	Location                 Location   `json:"location"`
	Version                  string     `json:"version"`
	Autoscaling              IntBoolean `json:"autoscaling"`
	CNIPlugin                string     `json:"cni_plugin"`
	APIServerAllowList       string     `json:"api_server_allow_list"`
	MaintenanceWindowStart   string     `json:"maintenance_window_start"`
	MaintenanceWindowLength  string     `json:"maintenance_window_length"`
	ManageInternalIPv4Prefix IntBoolean `json:"manage_internal_ipv4_prefix"`
	InternalIPv4Prefix       string     `json:"internal_ipv4_prefix"`
	NeedsServiceVMs          IntBoolean `json:"needs_service_vms"`
	EnableNATGateways        IntBoolean `json:"enable_nat_gateways"`
	EnableLBaaS              IntBoolean `json:"enable_lbaas"`
	ExternalIPFamilies       string     `json:"external_ip_families"`
	ManageExternalIPv4Prefix IntBoolean `json:"manage_external_ipv4_prefix"`
	ExternalIPv4Prefix       string     `json:"external_ipv4_prefix"`
	ManageExternalIPv6Prefix IntBoolean `json:"manage_external_ipv6_prefix"`
	ExternalIPv6Prefix       string     `json:"external_ipv6_prefix"`
}
