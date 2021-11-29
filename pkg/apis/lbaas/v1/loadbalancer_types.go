package v1

// LoadBalancerState describes the status of a given LoadBalancer
type LoadBalancerState struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type int    `json:"type"`
}

// RuleInfo holds the name and identifier of a rule.
type RuleInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

// anxcloud:object:hooks=RequestBodyHook

// LoadBalancer holds the information of a load balancer instance.
type LoadBalancer struct {
	CustomerIdentifier string            `json:"customer_identifier"`
	ResellerIdentifier string            `json:"reseller_identifier"`
	Identifier         string            `json:"identifier" anxcloud:"identifier"`
	Name               string            `json:"name"`
	IpAddress          string            `json:"ip_address"`
	AutomationRules    []RuleInfo        `json:"automation_rules"`
	State              LoadBalancerState `json:"state"`
}
