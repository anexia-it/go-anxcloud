package frontend

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

// Frontend represents a LBaaS Frontend.
type Frontend = v1.Frontend
type Definition struct {
	Name           string             `json:"name,omitempty"`
	LoadBalancer   string             `json:"load_balancer,omitempty"`
	DefaultBackend string             `json:"default_backend,omitempty"`
	Mode           common.Mode        `json:"mode,omitempty"`
	State          common.State       `json:"state,omitempty"`
	Enable         common.EnableState `json:"enable,omitempty"`
}
