package vm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
)

const (
	provisioningPathPrefix = "/api/vsphere/v1/provisioning/vm.json"
	// DefaultCPUPerformanceType to be used if a VM definition is created
	// by NewDefinition.
	DefaultCPUPerformanceType = "performance"
	// DefaultDiskType to be used if a VM definition is created
	// by NewDefinition.
	DefaultDiskType = "ENT2"
)

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
	// Example: "vmxnet3"
	NICType string `json:"nic_type,omitempty"`

	// Example: "791e8c171e654b459a7fcbbc07675cf3"
	VLAN string `json:"vlan,omitempty"`

	// Example: [ "identifier1", "identifier2", "10.11.12.13", "1.0.0.1" ]
	IPs []string `json:"ips,omitempty"`
}

// NewDefinition create a VM definition with the mandatory values set.
func NewDefinition(location, templateType, templateID, hostname string, cpus, memory, disk int, network []Network) Definition {
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

// ProvisioningResponse contains information returned by the API regarding a newly created VM.
type ProvisioningResponse struct {
	Progress   int      `json:"progress"`
	Errors     []string `json:"errors"`
	Identifier string   `json:"identifier"`
	Queued     bool     `json:"queued"`
}

// ProvisionVM issues a request to provision a new VM using the given VM definition.
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// definition contains the definition of the VM to be created.
// client is the HTTP to be used for the request.
//
// If the API returns errors, they are raised as ResponseError error.
// The returned ProvisioningResponse is still valid in this case.
func ProvisionVM(ctx context.Context, definition Definition, c client.Client) (ProvisioningResponse, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&definition); err != nil {
		panic(fmt.Sprintf("could not encode definition: %v", err))
	}

	url := fmt.Sprintf(
		"https://%s%s/%s/%s/%s",
		client.DefaultHost,
		provisioningPathPrefix,
		definition.Location,
		definition.TemplateType,
		definition.TemplateID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not create VM provisioning request: %w", err)
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not execute VM provisioning request: %w", err)
	}
	var responsePayload ProvisioningResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not decode VM provisioning response: %w", err)
	}

	if len(responsePayload.Errors) != 0 {
		err = &ProvisioningError{responsePayload.Errors}
	}

	return responsePayload, err
}
