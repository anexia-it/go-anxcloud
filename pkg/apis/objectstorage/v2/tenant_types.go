package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// Tenant represents a tenant resource in the Object Storage API.
type Tenant struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`
	AutomationRules    []AutomationRule       `json:"automation_rules,omitempty"`

	Name        string                 `json:"name"`
	State       *GenericAttributeState `json:"state,omitempty"`
	RemoteID    *string                `json:"remote_id,omitempty"`
	Description string                 `json:"description"`
	UserName    string                 `json:"user_name"`
	Password    string                 `json:"password,omitempty"`
	Quota       *float64               `json:"quota,omitempty"`
	Usage       *float64               `json:"usage,omitempty"`
	Backend     common.PartialResource `json:"backend"`
}
