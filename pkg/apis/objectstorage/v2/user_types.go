package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// User represents a user resource in the Object Storage API.
type User struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`

	UserName string                 `json:"user_name"`
	State    *GenericAttributeState `json:"state,omitempty"`
	Enabled  *bool                  `json:"enabled,omitempty"`
	FullName string                 `json:"full_name"`
	Backend  common.PartialResource `json:"backend"`
	Tenant   common.PartialResource `json:"tenant"`
	RemoteID *string                `json:"remote_id,omitempty"`
}
