package server

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Server = v1.Server

// Definition describes how a server resource should look like
type Definition struct {
	Name    string       `json:"name,omitempty"`
	State   common.State `json:"state,omitempty"`
	IP      string       `json:"ip,omitempty"`
	Port    int          `json:"port,omitempty"`
	Backend string       `json:"backend,omitempty"`
	Check   string       `json:"check,omitempty"`
}

const (
	CheckEnabled  = "enabled"
	CheckDisabled = "disabled"
)
