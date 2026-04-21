package server

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Server = v1.Server

// Definition describes how a server resource should look like
type Definition struct {
	Name    string       `json:"name"`
	State   common.State `json:"state"`
	IP      string       `json:"ip"`
	Port    int          `json:"port"`
	Backend string       `json:"backend"`
	Check   string       `json:"check"`
}

const (
	CheckEnabled  = "enabled"
	CheckDisabled = "disabled"
)
