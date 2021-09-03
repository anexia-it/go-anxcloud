package backend

import "github.com/anexia-it/go-anxcloud/pkg/lbaas/common"

type Definition struct {
	Name         string       `json:"name"`
	State        common.State `json:"state"`
	LoadBalancer string       `json:"load_balancer"`
	Mode         common.Mode  `json:"mode"`
}
