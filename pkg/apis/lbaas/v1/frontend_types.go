package v1

import "go.anx.io/go-anxcloud/pkg/apis/internal/gs"

// anxcloud:object:hooks=RequestBodyHook

// Frontend represents a LBaaS Frontend.
type Frontend struct {
	gs.GenericService
	HasState

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	Mode               Mode       `json:"mode"`
	ClientTimeout      string     `json:"client_timeout,omitempty"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	LoadBalancer *LoadBalancer `json:"load_balancer,omitempty"`

	// Only the name and identifier fields are used and returned.
	DefaultBackend *Backend `json:"default_backend,omitempty"`
}
