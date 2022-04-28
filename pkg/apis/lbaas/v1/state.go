package v1

import "encoding/json"

type StateRetriever interface {
	StateSuccess() bool
	StateProgressing() bool
	StateFailure() bool
}

type State struct {
	// programatically usable enum value
	ID string `json:"id"`

	// human readable status text
	Text string `json:"text"`

	Type int `json:"type"`
}

// StateSuccess checks if the state is one of the successful ones
func (s State) StateSuccess() bool {
	return s.ID == Updated.ID || s.ID == Deployed.ID
}

// StateProgressing checks if the state is marking any change currently being applied
func (s State) StateProgressing() bool {
	return s.ID == Updating.ID || s.ID == NewlyCreated.ID
}

// StateFailure checks if the state is marking any failure
func (s State) StateFailure() bool {
	return s.ID == DeploymentError.ID
}

func (s State) MarshalJSON() ([]byte, error) {
	// it would be great if one of the proposals in https://github.com/golang/go/issues/11939 would be
	// accepted and we could do something like `if s.ID == "" { omitThisField }` ... but it isn't, so
	// we have to override the field for every LBaaS Object for Create and Update operations.
	return json.Marshal(s.ID)
}

var (
	Updating        = State{ID: "0", Text: "Updating", Type: 0}
	Updated         = State{ID: "1", Text: "Updated", Type: 1}
	DeploymentError = State{ID: "2", Text: "DeploymentError", Type: 2}
	Deployed        = State{ID: "3", Text: "Deployed", Type: 3}
	NewlyCreated    = State{ID: "4", Text: "NewlyCreated", Type: 4}
)

type HasState struct {
	State State `json:"state"`
}

// StateSuccess checks if the state is one of the successful ones
func (hs HasState) StateSuccess() bool { return hs.State.StateSuccess() }

// StateProgressing checks if the state is marking any change currently being applied
func (hs HasState) StateProgressing() bool { return hs.State.StateProgressing() }

// StateFailure checks if the state is marking any failure
func (hs HasState) StateFailure() bool { return hs.State.StateFailure() }
