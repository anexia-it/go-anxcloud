package backend

import "go.anx.io/go-anxcloud/pkg/lbaas/common"

type Definition struct {
	Name         string       `json:"name"`
	State        common.State `json:"state"`
	LoadBalancer string       `json:"load_balancer"`
	Mode         common.Mode  `json:"mode"`
	HealthCheck  string       `json:"health_check"`
}
