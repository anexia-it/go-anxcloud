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
	Name      string       `json:"name"`
	IpAddress string       `json:"ip_address"`
	State     common.State `json:"state"`
}
