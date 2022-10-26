package v1

import (
	"encoding/json"

	"go.anx.io/go-anxcloud/pkg/apis/internal/gs"
)

// anxcloud:object:hooks=RequestBodyHook,ResponseFilterHook

// LoadBalancer holds the information of a load balancer instance.
type LoadBalancer struct {
	gs.GenericService
	State LoadBalancerState `json:"state"`

	CustomerIdentifier string     `json:"customer_identifier,omitempty"`
	ResellerIdentifier string     `json:"reseller_identifier,omitempty"`
	Identifier         string     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name               string     `json:"name"`
	IpAddress          string     `json:"ip_address"`
	AutomationRules    []RuleInfo `json:"automation_rules,omitempty"`
}

// LoadBalancerState is the same as State, but with different states as they are defined differently for
// LoadBalancer resources. It still implements StateRetriever, so you can use it like the other resources.
type LoadBalancerState struct {
	// programatically usable enum value
	ID string `json:"id"`

	// human readable status text
	Text string `json:"text"`

	Type int `json:"type"`
}

// StateSuccess checks if the state is one of the successful ones
func (s LoadBalancerState) StateSuccess() bool {
	return s.ID == LoadBalancerStateOK.ID
}

// StateProgressing checks if the state is marking any change currently being applied
func (s LoadBalancerState) StateProgressing() bool {
	return s.ID == LoadBalancerStatePending.ID || s.ID == LoadBalancerStateCreated.ID
}

// StateFailure checks if the state is marking any failure
func (s LoadBalancerState) StateFailure() bool {
	return s.ID == LoadBalancerStateError.ID
}

func (s LoadBalancerState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}

var (
	LoadBalancerStateOK      = LoadBalancerState{ID: "0", Text: "OK", Type: 0}
	LoadBalancerStateError   = LoadBalancerState{ID: "1", Text: "Error", Type: 1}
	LoadBalancerStatePending = LoadBalancerState{ID: "2", Text: "Pending", Type: 2}
	LoadBalancerStateCreated = LoadBalancerState{ID: "3", Text: "Created", Type: 3}
)
