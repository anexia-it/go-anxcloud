package gs

import (
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
	StateSuccess() bool
	StateProgressing() bool
	StateFailure() bool
}
