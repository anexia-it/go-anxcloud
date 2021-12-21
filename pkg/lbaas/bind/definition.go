package bind

import "go.anx.io/go-anxcloud/pkg/lbaas/common"

type Definition struct {
	Name     string       `json:"name"`
	State    common.State `json:"state"`
	Frontend string       `json:"frontend"`
}
