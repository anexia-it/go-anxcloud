package nodepool

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

const (
	OSFlatcar = "Flatcar Linux"
	Gibibyte  = 1024 * 1024 * 1024
)

var (
	StateOK   = gs.State{ID: "0", Text: "OK", Type: gs.StateTypeOK}
	StateNoGA = gs.State{ID: "1", Text: "OK", Type: gs.StateTypeOK}
)

type Definition struct {
	GSBase

	CriticalOperationPassword  string `json:"critical_operation_password"`
	CriticalOperationConfirmed bool   `json:"critical_operation_confirmed"`

	Cluster            common.PartialResource `json:"cluster"`
	SyncSource         SyncSource             `json:"syncsource"`
	Replicas           uint                   `json:"replicas"`
	CPUs               uint                   `json:"cpus"`
	CPUType            string                 `json:"cputype"`
	MemoryBytes        uint64                 `json:"memory"`
	DiskSizeBytes      uint64                 `json:"disk_size"`
	OperatingSystem    string                 `json:"operating_system"`
	AutoscalerEnabled  bool                   `json:"autoscaler_enabled"`
	AutoscalerMinNodes bool                   `json:"autoscaler_min_nodes"`
	AutoscalerMaxNodes bool                   `json:"autoscaler_max_nodes"`

	Disks    []NodepoolDisks `json:"disks"`
	Networks []NodepoolDisks `json:"networks"`

	CustomDNSEnabled bool   `json:"customdns_enabled"`
	DNSOverrideIPv4  bool   `json:"dns_override_ipv4"`
	DNSv4_1          string `json:"dns_v4_1"`
	DNSv4_2          string `json:"dns_v4_2"`

	DNSOverrideIPv6 bool   `json:"dns_override_ipv6"`
	DNSv6_1         string `json:"dns_v6_1"`
	DNSv6_2         string `json:"dns_v6_2"`

	Taints      string `json:"taints"`
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	SSHPubKeys  string `json:"sshpubkeys"`
}
