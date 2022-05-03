package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (obj *VLAN) DeepCopy() types.Object {
	// Initialize arrays
	

	out := &VLAN {
		// Primitives
		Identifier: obj.Identifier,
		Name: obj.Name,
		DescriptionCustomer: obj.DescriptionCustomer,
		RoleText: obj.RoleText,
		Status: obj.Status,
		VMProvisioning: obj.VMProvisioning,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}
