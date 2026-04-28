package frontend

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

// Frontend represents a LBaaS Frontend.
type Frontend = v1.Frontend
type Definition struct {
	Name           string             `json:"name"`
	LoadBalancer   string             `json:"load_balancer"`
	DefaultBackend string             `json:"default_backend"`
	Mode           common.Mode        `json:"mode"`
	State          common.State       `json:"state"`
	Enable         common.EnableState `json:"enable"`
}
