package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

type LoadBalancerImplementation string

const (
	LoadBalancerImplementationHAProxy LoadBalancerImplementation = "haproxy"
)

// anxcloud:object

// Cluster holds the information of a load balancing cluster
type Cluster struct {
	gs.GenericService
	gs.HasState

	Identifier       string                     `json:"identifier,omitempty" anxcloud:"identifier"`
	Name             string                     `json:"name,omitempty"`
	Implementation   LoadBalancerImplementation `json:"implementation,omitempty"`
	FrontendPrefixes *gs.PartialResourceList    `json:"frontend_prefixes,omitempty"`
	BackendPrefixes  *gs.PartialResourceList    `json:"backend_prefixes,omitempty"`
	Replicas         *int                       `json:"replicas,omitempty"`
}
