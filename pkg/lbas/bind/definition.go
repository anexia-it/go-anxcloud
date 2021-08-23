package bind

import "github.com/anexia-it/go-anxcloud/pkg/lbas/common"

type Definition struct {
	Name     string       `json:"name"`
	State    common.State `json:"state"`
	Frontend string       `json:"frontend"`
}
