package v1

// anxcloud:object:hooks=RequestBodyHook

// Bind represents an LBaaS FrontendBind
type Bind struct {
	HasState

	CustomerIdentifier string     `json:"customer_identifier"`
	ResellerIdentifier string     `json:"reseller_identifier"`
	Identifier         string     `json:"identifier" anxcloud:"identifier"`
	Name               string     `json:"name"`
	Address            string     `json:"address"`
	Port               int        `json:"port"`
	SSL                bool       `json:"ssl"`
	SslCertificatePath string     `json:"ssl_certificate_path,omitempty"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`

	// Only the name and identifier fields are used and returned.
	Frontend Frontend `json:"frontend"`
}
