package loadbalancer

import (
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Definition struct {
	Name      string       `json:"name"`
	IpAddress string       `json:"ip_address"`
	State     common.State `json:"state"`
}
