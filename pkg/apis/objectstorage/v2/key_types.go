package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"time"
)

// anxcloud:object:hooks=RequestBodyHook

// Key represents a key resource in the Object Storage API.
type Key struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`

	RemoteID   *string                 `json:"remote_id,omitempty"`
	State      *GenericAttributeState  `json:"state,omitempty"`
	Backend    *common.PartialResource `json:"backend,omitempty"`
	Tenant     *common.PartialResource `json:"tenant,omitempty"`
	User       *common.PartialResource `json:"user,omitempty"`
	ExpireDate *time.Time              `json:"expiry_date,omitempty"`
	Name       string                  `json:"name"`
}
