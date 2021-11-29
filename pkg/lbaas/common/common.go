package common

import (
	v1 "github.com/anexia-it/go-anxcloud/pkg/apis/lbaas/v1"
)

type (
	Mode  = v1.Mode
	State = v1.State
)

const (
	HTTP = v1.HTTP
	TCP  = v1.TCP

	Updating        = v1.Updating
	Updated         = v1.Updated
	DeploymentError = v1.DeploymentError
	Deployed        = v1.Deployed
	NewlyCreated    = v1.NewlyCreated
)
