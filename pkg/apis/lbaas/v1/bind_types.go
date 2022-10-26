package v1

import "go.anx.io/go-anxcloud/pkg/apis/internal/gs"

// anxcloud:object:hooks=RequestBodyHook

// Bind represents an LBaaS FrontendBind
type Bind struct {
	gs.GenericService
	HasState

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	Address            string     `json:"address"`
	Port               int        `json:"port"`
	SSL                bool       `json:"ssl"`
	SslCertificatePath string     `json:"ssl_certificate_path,omitempty"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	Frontend Frontend `json:"frontend"`
}
