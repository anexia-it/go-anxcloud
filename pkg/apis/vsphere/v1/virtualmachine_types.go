package v1

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

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
	// Identifier is the unique 32 character identifier of this VM resource.
	Identifier string `json:"identifier" anxcloud:"identifier"`
	// Name is the hostname of the VM.
	Name string `json:"name"`
	// CustomName is custom name that can be added.
	CustomName string `json:"custom_name,omitempty"`
	// GuestOS contains the guest os full name.
	GuestOS string `json:"guest_os"`
	// Firmware contains the firmware used to boot the VM.
	Firmware string `json:"firmware"`
	// Status contains the power status of the VM.
	Status Status `json:"status"`
	// RAM specifies the amount of RAM available for the VM.
	RAM int `json:"ram"`
	// CPU contains the amount of CPUs (sockets) available for the VM.
	// By default, when only specifying 'Cores', CPU = Cores. Cores must be a multiple of CPU (e.g. 1 CPU, 2 Cores).
	CPU int `json:"cpu"`
	// CPUClockRate specifies the maximum clock rate of a CPU core.
	CPUClockRate int `json:"cpu_clock_rate"`
	// CPUPerformanceType specifies a performance category that relates to the CPUClockRate.
	// When creating a VM, specify the type ('best-effort', 'standard', 'enterprise', 'performance')
	// combined with '-intel' or '-amd', e.g. 'standard-amd'.
	CPUPerformanceType string `json:"cpu_performance_type"`
	// VTPMEnabled specifies if vTPM should be enabled for the VM.
	VTPMEnabled bool `json:"vtpm_enabled"`
	// Cores specifies the total number of CPU cores available for the VM.
	Cores int `json:"cores"`
	// Disks contains the amount of disk devices available for the VM.
	Disks int `json:"disks"`
	// DiskInfo specifies disk devices available for the VM.
	DiskInfo []DiskInfo `json:"disk_info"`
	// Networks specifies network interfaces attached to the VM.
	Networks []Network `json:"networks"`
	// VersionTools specifies the guest OS tools version.
	VersionTools string `json:"version_tools"`
	// GuestToolsStatus contains the status of the guest OS tools (Active, Inactive)
	GuestToolsStatus string `json:"guest_tools_status"`
	// Location specifies the datacenter location of the VM.
	Location v1.Location `json:"location"`
	// ProvisioningLocationIdentifier contains the location identifier of the provisioning API.
	ProvisioningLocationIdentifier string `json:"provisioning_location_identifier"`
	// TemplateID contains the template identifier for creating the VM.
	TemplateID string `json:"template_id"`
	// Password is the initial root password when creating a VM. Use 'ssh_key' instead.
	Password string `json:"password,omitempty"`
	// SSHKey is a public key that is granted root access when creating a VM. Preferred over 'password'.
	SSHKey string `json:"ssh_key,omitempty"`
	// StartScript is a base64 encoded shell script that is executed after provisioning, when creating a VM.
	StartScript string `json:"start_script,omitempty"`
}

type CPUPerformanceType string

// TODO: use 'enum' or not?
const (
	CPUPerformanceBestEffort  string = "best-effort"
	CPUPerformanceStandard    string = "standard"
	CPUPerformanceEnterprise  string = "enterprise"
	CPUPerformancePerformance string = "performance"
)

// DiskInfo contains meta information of attached disks to a VM.
type DiskInfo struct {
	// BusType contains the disk bus that is used.
	BusType string `json:"bus_type"`
	// BusTypeLabel contains the actual name of the bus on the system.
	BusTypeLabel string `json:"bus_type_label"`
	// DiskGB specifies the disk size in GB.
	DiskGB int `json:"disk_gb"`
	// DiskID contains the numeric disk ID.
	DiskID int `json:"disk_id"`
	// DiskType specifies the disk type to be used.
	// 'STDx' for HDD, 'ENTx' for SSD, 'HPCx' for high IOPS and low latency SSD, 'LOCx' for local storage SSD.
	DiskType string `json:"disk_type"`
	// IOPS specifies the upper limit of allowed I/O operations.
	IOPS int `json:"iops"`
	// Latency contains the expected latency in milliseconds.
	Latency int `json:"latence"`
	// StorageType contains the disk type used ('HDD', 'SSD').
	StorageType string `json:"storage_type"`
}

// Network contains meta information of attached NICs to a VM.
type Network struct {
	// BandwidthLimit is the limit in MBit for the interface (-1, 100, 1000 or 10000).
	BandwidthLimit BandwidthLimit `json:"bandwidth_limit"`
	// ID is the numeric network ID (not the prefix identifier).
	ID int `json:"id"`
	// IPsv4 contains the IPs v4 that are assigned to the VM.
	IPsv4 []string `json:"ips_v4"`
	// IPsv6 contains the IPs v6 that are assigned to the VM.
	IPsv6 []string `json:"ips_v6"`
	// MACAddress contains the MAC of the VM.
	MACAddress string `json:"mac_address"`
	// NIC is the numeric NIC type ID.
	NIC int `json:"nic"`
	// NICType is only used when creating a VM.
	NICType string `json:"nic_type,omitempty"`
	// VLAN is an identifier of a VLAN that needs to exist.
	VLAN string `json:"vlan"`
	// IPs is only used when creating a VM and can contain IPv4, IPv6 or identifiers. TODO: use ips_v4+ips_v6 instead?
	IPs []string `json:"ips,omitempty"`
	// Connected specifies if this network is connected.
	Connected bool `json:"connected"`
	// ConnectedAtStart specifies if this network should be connected on startup of the VM.
	ConnectedAtStart bool `json:"connected_at_start"`
}

type BandwidthLimit int

const (
	BandwidthUnlimited BandwidthLimit = -1
	Bandwidth100MBit   BandwidthLimit = 100
	Bandwidth1GBit     BandwidthLimit = 1000
	Bandwidth10Git     BandwidthLimit = 10000
)
