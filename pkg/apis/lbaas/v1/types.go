package v1

import "encoding/json"

// Mode is an enum for the supported LoadBalancer protocols.
type Mode string

const (
	TCP  Mode = "tcp"
	HTTP Mode = "http"
)

type State struct {
	// programatically usable enum value
	ID string `json:"id"`

	// human readable status text
	Text string `json:"text"`

	Type int `json:"type"`
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}

var (
	Updating        = State{ID: "0", Text: "Updating", Type: 0}
	Updated         = State{ID: "1", Text: "Updated", Type: 1}
	DeploymentError = State{ID: "2", Text: "DeploymentError", Type: 2}
	Deployed        = State{ID: "3", Text: "Deployed", Type: 3}
	NewlyCreated    = State{ID: "4", Text: "NewlyCreated", Type: 4}
)

// RuleInfo holds the name and identifier of a rule.
type RuleInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}
