package v1

// anxcloud:object:hooks=RequestBodyHook

// The Backend resource configures settings common for all specific backend Server resources linked to it.
type Backend struct {
	CustomerIdentifier string     `json:"customer_identifier"`
	ResellerIdentifier string     `json:"reseller_identifier"`
	Identifier         string     `json:"identifier" anxcloud:"identifier"`
	Name               string     `json:"name"`
	HealthCheck        string     `json:"health_check"`
	Mode               Mode       `json:"mode"`
	ServerTimeout      int        `json:"server_timeout"`
	State              State      `json:"state"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	LoadBalancer LoadBalancer `json:"load_balancer"`
}
