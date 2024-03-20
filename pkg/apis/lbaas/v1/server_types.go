package v1

import "go.anx.io/go-anxcloud/pkg/apis/common/gs"

// anxcloud:object

// Server holds the information of a load balancers backend server
type Server struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	IP                 string     `json:"ip"`
	Port               int        `json:"port"`
	Check              string     `json:"check,omitempty"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	Backend Backend `json:"backend" anxencode:"flatten"`
}
