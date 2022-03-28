package v1

// anxcloud:object:hooks=RequestBodyHook

// Frontend represents a LBaaS Frontend.
type Frontend struct {
	HasState

	CustomerIdentifier string     `json:"customer_identifier"`
	ResellerIdentifier string     `json:"reseller_identifier"`
	Identifier         string     `json:"identifier" anxcloud:"identifier"`
	Name               string     `json:"name"`
	Mode               Mode       `json:"mode"`
	ClientTimeout      string     `json:"client_timeout"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	LoadBalancer *LoadBalancer `json:"load_balancer,omitempty"`

	// Only the name and identifier fields are used and returned.
	DefaultBackend *Backend `json:"default_backend,omitempty"`
}
