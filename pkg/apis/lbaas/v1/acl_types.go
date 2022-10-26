package v1

import "go.anx.io/go-anxcloud/pkg/apis/internal/gs"

// anxcloud:object:hooks=RequestBodyHook

// ACL represents an LBaaS ACL
type ACL struct {
	gs.GenericService
	HasState

	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`

	Identifier      string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name            string     `json:"name,omitempty"`
	ParentType      string     `json:"parent_type,omitempty" anxcloud:"filterable"`
	Criterion       string     `json:"criterion,omitempty"`
	Value           string     `json:"value,omitempty"`
	AutomationRules []RuleInfo `json:"automation_rules,omitempty"`

	// Index is *int to allow zero values but also omitempty
	// pkg/utils/pointer can be used to create pointers from primitives
	Index *int `json:"index,omitempty"`

	// Only the name and identifier fields are used and returned.
	Frontend Frontend `json:"frontend,omitempty" anxcloud:"filterable"`
	Backend  Backend  `json:"backend,omitempty" anxcloud:"filterable"`
}
