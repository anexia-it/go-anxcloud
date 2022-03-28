package v1

// Mode is an enum for the supported LoadBalancer protocols.
type Mode string

const (
	TCP  Mode = "tcp"
	HTTP Mode = "http"
)

// RuleInfo holds the name and identifier of a rule.
type RuleInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}
