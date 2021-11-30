package v1

type Mode string

const (
	TCP  = Mode("tcp")
	HTTP = Mode("http")
)

type State string

const (
	Updating        = State("0")
	Updated         = State("1")
	DeploymentError = State("2")
	Deployed        = State("3")
	NewlyCreated    = State("4")
)

// StateObject describes the status of a given resource, including programatic usable and human readable values.
type StateObject struct {
	// programatically usable enum value
	ID State `json:"id"`

	// human readable status text
	Text string `json:"text"`

	Type int `json:"type"`
}

// RuleInfo holds the name and identifier of a rule.
type RuleInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}
