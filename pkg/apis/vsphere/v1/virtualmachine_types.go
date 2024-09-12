package v1

import v1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"

type Status string

const (
	// StatusPoweredOn means the Virtual Machine is powered on.
	StatusPoweredOn Status = "poweredOn"

	// StatusActive means the Virtual Machine is powered off.
	StatusPoweredOff Status = "poweredOff"
)

// anxcloud:object:hooks=RequestBodyHook,FilterRequestURLHook

// VirtualMachine represents an Anexia Dynamic Compute resource used for vm provisioning
// should contain everything we need to interact with the API. Does not need to represent the Engine API object as a whole.
type VirtualMachine struct {
	Identifier                     string      `json:"identifier" anxcloud:"identifier"`
	Name                           string      `json:"name"`
	CustomName                     string      `json:"custom_name"`
	GuestOS                        string      `json:"guest_os"`
	Firmware                       string      `json:"firmware"`
	Status                         Status      `json:"status"`
	RAM                            int         `json:"ram"`
	CPU                            int         `json:"cpu"`
	CPUClockRate                   int         `json:"cpu_clock_rate"`
	CPUPerformanceType             string      `json:"cpu_performance_type"` // "-"
	VTPMEnabled                    bool        `json:"vtpm_enabled"`
	Cores                          int         `json:"cores"`
	Disks                          int         `json:"disks"`
	DiskInfo                       []DiskInfo  `json:"disk_info"`
	Network                        []Network   `json:"network"`
	VersionTools                   string      `json:"version_tools"`
	GuestToolsStatus               string      `json:"guest_tools_status"` // Active,Inactive
	Location                       v1.Location `json:"location"`
	ProvisioningLocationIdentifier string      `json:"provisioning_location_identifier"`
	TemplateID                     string      `json:"template_id"`
	ResourceSalesperson            string      `json:"resource_salesperson"`
}

// DiskInfo contains meta information of attached disks to a VM.
type DiskInfo struct {
	BusType      string  `json:"bus_type"`
	BusTypeLabel string  `json:"bus_type_label"`
	DiskGB       float64 `json:"disk_gb"`
	DiskID       int     `json:"disk_id"`
	DiskType     string  `json:"disk_type"`
	IOPS         int     `json:"iops"`
	Latency      int     `json:"latence"`
	StorageType  string  `json:"storage_type"`
}

// Network contains meta information of attached NICs to a VM.
type Network struct {
	BandwidthLimit int      `json:"bandwidth_limit"`
	ID             int      `json:"id"`
	IPv4           []string `json:"ips_v4"`
	IPv6           []string `json:"ips_v6"`
	MACAddress     string   `json:"mac_address"`
	NIC            int      `json:"nic"`
	VLAN           string   `json:"vlan"`
}
