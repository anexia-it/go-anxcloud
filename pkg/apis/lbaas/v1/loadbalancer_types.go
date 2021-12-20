package v1

// anxcloud:object:hooks=RequestBodyHook

// LoadBalancer holds the information of a load balancer instance.
type LoadBalancer struct {
	CustomerIdentifier string     `json:"customer_identifier"`
	ResellerIdentifier string     `json:"reseller_identifier"`
	Identifier         string     `json:"identifier" anxcloud:"identifier"`
	Name               string     `json:"name"`
	IpAddress          string     `json:"ip_address"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`
	State              State      `json:"state"`
}
