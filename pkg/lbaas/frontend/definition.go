package frontend

import "github.com/anexia-it/go-anxcloud/pkg/lbaas/common"

type Definition struct {
	Name           string       `json:"name"`
	LoadBalancer   string       `json:"load_balancer"`
	DefaultBackend string       `json:"default_backend"`
	Mode           common.Mode  `json:"mode"`
	State          common.State `json:"state"`
}
