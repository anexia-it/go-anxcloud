package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// Region represents a region resource in the Object Storage API.
type Region struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`

	Name        string                  `json:"name"`
	State       *GenericAttributeState  `json:"state,omitempty"`
	Description string                  `json:"description"`
	Backend     *common.PartialResource `json:"backend,omitempty"`
}
