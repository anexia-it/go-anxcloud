package v1

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object

// LoadBalancer holds the information of a load balancer instance.
type LoadBalancer struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	IpAddress          string     `json:"ip_address"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`
}
