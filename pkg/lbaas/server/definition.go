package server

import "go.anx.io/go-anxcloud/pkg/lbaas/common"

// Definition describes how a server resource should look like
type Definition struct {
	Name    string       `json:"name"`
	State   common.State `json:"state"`
	IP      string       `json:"ip"`
	Port    int          `json:"port"`
	Backend string       `json:"backend"`
}
