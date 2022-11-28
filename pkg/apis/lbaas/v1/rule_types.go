package v1

import "go.anx.io/go-anxcloud/pkg/apis/common/gs"

// anxcloud:object:hooks=RequestBodyHook

// Rule represents an LBaaS Rule
type Rule struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`

	Identifier       string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name             string `json:"name,omitempty"`
	ParentType       string `json:"parent_type,omitempty" anxcloud:"filterable"`
	Condition        string `json:"condition,omitempty" anxcloud:"filterable"`
	ConditionTest    string `json:"condition_test,omitempty"`
	Type             string `json:"type,omitempty" anxcloud:"filterable"`
	Action           string `json:"action,omitempty" anxcloud:"filterable"`
	RedirectionType  string `json:"redirection_type,omitempty" anxcloud:"filterable"`
	RedirectionValue string `json:"redirection_value,omitempty"`
	RedirectionCode  string `json:"redirection_code,omitempty" anxcloud:"filterable"`
	RuleType         string `json:"rule_type,omitempty" anxcloud:"filterable"`

	// Index is *int to allow zero values but also omitempty
	// pkg/utils/pointer can be used to create pointers from primitives
	Index *int `json:"index,omitempty"`

	// Only the name and identifier fields are used and returned.
	Frontend Frontend `json:"frontend,omitempty" anxcloud:"filterable"`
	Backend  Backend  `json:"backend,omitempty" anxcloud:"filterable"`
}
