package vm

// Definition states the configuration of a VM within the anxcloud.
type Definition struct {
	Location     string `json:"-"`
	TemplateType string `json:"-"`
	TemplateID   string `json:"-"`

	// Required - VM hostname
	// Example: ('my-awesome-new-vm', 'web-001', 'db-001', …)
	// The hostname will be auto prefixed with your customer id.
	Hostname string `json:"hostname"`

	// Required - Memory in MB
	// Example: (1024, 2048, 4096, 8192, …)
	// Default: as given in template.
	Memory int `json:"memory_mb"`

	// Required - Amount of CPUs
	// Example: (1, 2, 3, 4 ,…)
	// Default: as given in template.
	CPUs int `json:"cpus"`

	// Required - Disk capacity in GB
	// Example: (1, 2, 4, 5, …)
	// Default: as given in template.
	Disk int `json:"disk_gb"`

	// Disk category (limits disk performance, e.g. IOPS)
	// Example: ('STD1', 'ENT2','HPC1',…)
	// Default: as defined by data center.
	DiskType string `json:"disk_type,omitempty"`

	// Disks in addition to the primary disk
	AdditionalDisks []AdditionalDisk `json:"additional_disks,omitempty"`

	// CPU type
	// Example: ("best-effort", "standard", "enterprise", "performance")
	// Default: "standard".
	CPUPerformanceType string `json:"cpu_performance_type,omitempty"`

	// Amount of CPU sockets Number of cores have to be a multiple of sockets, as they will be spread evenly across all sockets.
	// Default: number of cores, i.e. one socket per CPU core.
	Sockets int `json:"sockets,omitempty"`

	// Network interfaces
	// IPs are ignored when using template_type "from_scratch".
	Network []Network `json:"network,omitempty"`

	// Primary DNS server
	// Example: '94.16.16.94'
	// Default: as given in template.
	DNS1 string `json:"dns1,omitempty"`

	// Secondary DNS server
	// Example: '94.16.16.16'
	// Default: as given in template.
	DNS2 string `json:"dns2,omitempty"`

	// Tertiary DNS server
	// Example: '2a00:11c0:11c0::2a00'
	// Default: as given in template.
	DNS3 string `json:"dns3,omitempty"`

	// Quaternary DNS server
	// Example: '2a00:11c0:11c0::11c0'
	// Default: as given in template.
	DNS4 string `json:"dns4,omitempty"`

	// Plaintext password
	// Example: ('!anx123mySuperStrongPassword123anx!', 'go3ju0la1ro3', …)
	// USE IT AT YOUR OWN RISK! (or SSH key instead).
	Password string `json:"password,omitempty"`

	// Public key (instead of password, only for Linux systems)
	// Recommended over providing a plaintext password.
	SSH string `json:"ssh,omitempty"`

	// Script to be executed after provisioning
	// Should be base64 encoded
	// Consider the corresponding shebang at the beginning of your script.
	// If you want to use PowerShell, the first line should be: #ps1_sysnative.
	Script string `json:"script,omitempty"`

	// Boot delay in seconds
	// Default: 0.
	BootDelay int `json:"boot_delay,omitempty"`

	// Start the VM into BIOS setup on next boot
	// Default: false.
	EnterBIOSSetup bool `json:"enter_bios_setup,omitempty"`

	// Customer identifier (reseller only).
	Organization string `json:"organization,omitempty"`
}

// Network defines the network configuration of a VM.
type Network struct {
	ID int `json:"id,omitempty"`

	// Example: "vmxnet3"
	NICType string `json:"nic_type,omitempty"`

	// Example: "791e8c171e654b459a7fcbbc07675cf3"
	VLAN string `json:"vlan,omitempty"`

	// Example: [ "identifier1", "identifier2", "10.11.12.13", "1.0.0.1" ]
	IPs []string `json:"ips,omitempty"`
}

// NewDefinition create a VM definition with the mandatory values set.
func (a api) NewDefinition(location, templateType, templateID, hostname string, cpus, memory, disk int, network []Network) Definition {
	return Definition{
		Location:           location,
		TemplateType:       templateType,
		TemplateID:         templateID,
		Hostname:           hostname,
		CPUs:               cpus,
		Memory:             memory,
		Disk:               disk,
		Network:            network,
		CPUPerformanceType: DefaultCPUPerformanceType,
		DiskType:           DefaultDiskType,
	}
}

// Disk represents a disk of a VM.
type Disk struct {
	ID      int    `json:"disk_id,omitempty"`
	Type    string `json:"disk_type"`
	SizeGBs int    `json:"disk_gb"`
}

// AdditionalDisk represents an additional disk which can be defined for provisioning
type AdditionalDisk struct {
	SizeGBs int    `json:"gb"`
	Type    string `json:"type"`
}

// Change contains information about requested VM change request.
type Change struct {
	MemoryMBs          int       `json:"memory_mb,omitempty"`
	CPUs               int       `json:"cpus,omitempty"`
	CPUSockets         int       `json:"sockets,omitempty"`
	CPUPerformanceType string    `json:"cpu_performance_type,omitempty"`
	DeleteDiskIDs      []int     `json:"disk_to_delete,omitempty"`
	AddDisks           []Disk    `json:"disk_to_add,omitempty"`
	ChangeDisks        []Disk    `json:"disk_to_change,omitempty"`
	AddNICs            []Network `json:"network_to_add,omitempty"`
	ChangeNICs         []Network `json:"network_to_change,omitempty"`
	BootDelaySecs      int       `json:"boot_delay,omitempty"`
	EnterBIOSSetup     bool      `json:"enter_bios_setup,omitempty"`
	Reboot             bool      `json:"force_restart_if_needed,omitempty"`
	EnableDangerous    bool      `json:"critical_operation_confirmed,omitempty"`
}

// NewChange create a VM change request with default values.
func NewChange() Change {
	return Change{
		Reboot:          true,
		EnableDangerous: true,
	}
}
