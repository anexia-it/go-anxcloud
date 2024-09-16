package v1

// anxcloud:object

// ProvisionProgress represents the progress of VM provisioning.
type ProvisionProgress struct {
	// Errors contains error messages if provisioning failed.
	Errors []string `json:"errors"`
	// Identifier is the identifier of the provisioning progress.
	Identifier string `json:"identifier" anxcloud:"identifier"`
	// Percent show the progress of completion.
	Percent int `json:"progress"`
	// Queued is true when the provisioning process is queued.
	Queued bool `json:"queued"`
	// Status contains the status of provisioning.
	// failed = '-1', success = '1', in progress = '2', cancelled = '3'
	Status string `json:"status,omitempty"`
	// VMIdentifier contains the VM identifier once it is available.
	VMIdentifier string `json:"vm_identifier,omitempty"`
}
