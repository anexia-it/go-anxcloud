package v1

import (
	"errors"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
)

var (
	// ErrLocationCount is returned when trying to create a Prefix with no or more than one location.
	ErrLocationCount = errors.New("Prefixes have to be created with exactly one Location")

	// ErrVLANCount is returned when trying to create a Prefix with more than one VLAN.
	ErrVLANCount = errors.New("Prefixes have to be created with exactly one or none VLAN")
)

// anxcloud:object

type Prefix struct {
	Identifier          string       `json:"identifier,omitempty" anxcloud:"identifier"`
	Name                string       `json:"name,omitempty"`
	DescriptionCustomer string       `json:"description_customer"`
	Version             Family       `json:"version,omitempty" anxcloud:"filterable"`
	Netmask             int          `json:"netmask,omitempty"`
	RoleText            string       `json:"role_text,omitempty"`
	Status              Status       `json:"status,omitempty" anxcloud:"filterable"`
	RouterRedundancy    bool         `json:"router_redundancy"`
	Type                AddressSpace `json:"type,omitempty" anxcloud:"filterable"`

	Locations []corev1.Location `json:"locations,omitempty" anxcloud:"filterable,location,single"`
	VLANs     []vlanv1.VLAN     `json:"vlans,omitempty"`
}
