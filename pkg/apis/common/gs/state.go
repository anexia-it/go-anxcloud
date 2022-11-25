package gs

import (
	"encoding/json"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// State represents the state object available on GS resources
type State struct {
	// programatically usable enum value
	ID string `json:"id"`

	// human readable status text
	Text string `json:"text"`

	Type int `json:"type"`
}

type objectWithStateRetriever interface {
	types.Object
	StateRetriever
}

// StateRetriever is an interface Objects can implement to provide unified state information
type StateRetriever interface {
	StateOK() bool
	StatePending() bool
	StateError() bool
}

const (
	// StateTypeError is used for states of type error
	StateTypeError = iota
	// StateTypeOK is used for states of type OK
	StateTypeOK
	// StateTypePending is used for states of type Pending
	StateTypePending
)

// StateOK checks if the state is one of the successful ones
func (s State) StateOK() bool {
	return s.Type == StateTypeOK
}

// StatePending checks if the state is marking any change currently being applied
func (s State) StatePending() bool {
	return s.Type == StateTypePending
}

// StateError checks if the state is marking any failure
func (s State) StateError() bool {
	return s.Type == StateTypeError
}

func (s State) MarshalJSON() ([]byte, error) {
	// it would be great if one of the proposals in https://github.com/golang/go/issues/11939 would be
	// accepted and we could do something like `if s.ID == "" { omitThisField }` ... but it isn't, so
	// we have to override the field for every LBaaS Object for Create and Update operations.
	return json.Marshal(s.ID)
}

type HasState struct {
	State State `json:"state"`
}

// StateOK checks if the state is one of the successful ones
func (hs HasState) StateOK() bool { return hs.State.StateOK() }

// StatePending checks if the state is marking any change currently being applied
func (hs HasState) StatePending() bool { return hs.State.StatePending() }

// StateError checks if the state is marking any failure
func (hs HasState) StateError() bool { return hs.State.StateError() }
