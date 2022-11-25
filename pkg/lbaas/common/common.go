package common

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
)

type (
	Mode  = v1.Mode
	State = gs.State
)

const (
	HTTP = v1.HTTP
	TCP  = v1.TCP
)

var (
	Updating        = v1.Updating
	Updated         = v1.Updated
	DeploymentError = v1.DeploymentError
	Deployed        = v1.Deployed
	NewlyCreated    = v1.NewlyCreated
)
