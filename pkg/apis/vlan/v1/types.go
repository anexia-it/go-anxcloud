package v1

import (
	"errors"
	"fmt"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

// Status describe the deployment status of a VLAN.
type Status string

const (
	// StatusInvalid is a client-internal, invalid status which is only used to check if a status is set for filtering.
	StatusInvalid Status = ""

	// StatusActive means the VLAN is fully deployed and ready to be used.
	StatusActive Status = "Active"

	// StatusPending is set when a VLAN is newly created and not yet ready to be used.
	StatusPending Status = "Pending"

	// StatusMarkedForDeletion means the VLAN was requested to be deleted, which is still pending.
	StatusMarkedForDeletion Status = "Marked for deletion"
)

var (
	// ErrFilterMultipleLocations is returned when trying to List VLANs with multiple locations configured
	// for filtering - which is not implemented.
	ErrFilterMultipleLocations = fmt.Errorf("%w: cannot filter on multiple locations", apiTypes.ErrInvalidFilter)

	// ErrLocationCount is returned when trying to create a VLAN with no or more than one location.
	ErrLocationCount = errors.New("VLANs have to be created with exactly one Location")
)

// anxcloud:object:hooks=RequestBodyHook

// VLAN describes a virtual network IP prefixes, virtual machines (if VMProvisioning is true) and
// alike can be deployed into.
type VLAN struct {
	Identifier          string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name                string `json:"name,omitempty"`
	DescriptionCustomer string `json:"description_customer,omitempty"`
	RoleText            string `json:"role_text,omitempty"`
	Status              Status `json:"status,omitempty" anxcloud:"filterable"`
	VMProvisioning      bool   `json:"vm_provisioning"`

	// The API returns an array of locations, but there is no way to configure more than one location via the API.
	// Additionally, not even the one location can be updated via API.
	// When creating a VLAN pass a single Location object, only the Identifier needs to be set on it.
	Locations []corev1.Location `json:"locations,omitempty" anxcloud:"filterable"`
}
