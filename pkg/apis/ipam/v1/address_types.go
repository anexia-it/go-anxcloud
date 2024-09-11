package v1

import (
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
)

// anxcloud:object

type Address struct {
	Identifier          string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name                string `json:"name,omitempty"`
	DescriptionCustomer string `json:"description_customer"`
	Version             Family `json:"version" anxcloud:"filterable"`
	RoleText            string `json:"role_text,omitempty" anxcloud:"filterable"`
	Status              Status `json:"status,omitempty" anxcloud:"filterable"`

	VLAN vlanv1.VLAN `json:"-" anxcloud:"filterable,vlan"`

	// Prefix of the address, only for filtering and creating and not returned by the API.
	Prefix Prefix `json:"-" anxcloud:"filterable,prefix"`

	// Location of this address, only for filtering and not returned by the API.
	Location corev1.Location `json:"-" anxcloud:"filterable,location"`

	// Type of this address (public or private), only for filtering and not returned by the API.
	Type AddressType `json:"-" anxcloud:"filterable,type"`
}
