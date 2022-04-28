package v1

// anxcloud:object:hooks=RequestBodyHook

// The Backend resource configures settings common for all specific backend Server resources linked to it.
type Backend struct {
	commonMethods
	HasState

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	HealthCheck        string     `json:"health_check,omitempty"`
	Mode               Mode       `json:"mode"`
	ServerTimeout      int        `json:"server_timeout,omitempty"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	LoadBalancer LoadBalancer `json:"load_balancer"`
}
