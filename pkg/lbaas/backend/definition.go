package backend

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

// The Backend resource configures settings common for all specific backend Server resources linked to it.
type Backend = v1.Backend
type Definition struct {
	Name         string       `json:"name"`
	State        common.State `json:"state"`
	LoadBalancer string       `json:"load_balancer"`
	Mode         common.Mode  `json:"mode"`
	HealthCheck  string       `json:"health_check"`
}
