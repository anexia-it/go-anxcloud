package loadbalancer

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type (
	// Loadbalancer holds the information of a load balancer instance.
	Loadbalancer = v1.LoadBalancer
)

type Definition struct {
	Name      string       `json:"name,omitempty"`
	IpAddress string       `json:"ip_address,omitempty"`
	State     common.State `json:"state,omitempty"`
}
