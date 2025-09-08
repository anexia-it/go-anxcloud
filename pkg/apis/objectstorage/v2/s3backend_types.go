package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// S3Backend represents an S3 backend resource in the Object Storage API.
type S3Backend struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`

	Name            string                  `json:"name"`
	State           *GenericAttributeState  `json:"state,omitempty"`
	Endpoint        common.PartialResource  `json:"endpoint"`
	BackendType     *GenericAttributeSelect `json:"backend_type,omitempty"`
	Enabled         *bool                   `json:"enabled,omitempty"`
	BackendUser     string                  `json:"backend_user,omitempty"`
	BackendPassword string                  `json:"backend_password,omitempty"`
}
