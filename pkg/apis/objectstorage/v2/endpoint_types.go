package v2

import "go.anx.io/go-anxcloud/pkg/apis/common/gs"

// anxcloud:object:hooks=RequestBodyHook

// Endpoint represents an endpoint resource in the Object Storage API.
type Endpoint struct {
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

	Name             string                 `json:"name,omitempty"`
	URL              string                 `json:"url,omitempty"`
	State            *GenericAttributeState `json:"state,omitempty"`
	EndpointUser     string                 `json:"endpoint_user,omitempty"`
	EndpointPassword string                 `json:"endpoint_password,omitempty"`
	Enabled          bool                   `json:"enabled,omitempty"`
}
