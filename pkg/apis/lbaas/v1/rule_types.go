package v1

// anxcloud:object:hooks=RequestBodyHook

// Rule represents an LBaaS Rule
type Rule struct {
	commonMethods
	HasState

	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`

	Identifier       string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name             string `json:"name,omitempty"`
	ParentType       string `json:"parent_type,omitempty"`
	Index            int    `json:"index,omitempty"`
	Condition        string `json:"condition,omitempty"`
	ConditionTest    string `json:"condition_test,omitempty"`
	Type             string `json:"type,omitempty"`
	Action           string `json:"action,omitempty"`
	RedirectionType  string `json:"redirection_type,omitempty"`
	RedirectionValue string `json:"redirection_value,omitempty"`
	RedirectionCode  string `json:"redirection_code,omitempty"`
	RuleType         string `json:"rule_type,omitempty"`

	// Only the name and identifier fields are used and returned.
	Frontend Frontend `json:"frontend,omitempty"`
	Backend  Backend  `json:"backend,omitempty"`
}
