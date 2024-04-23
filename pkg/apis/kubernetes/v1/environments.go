package v1

import (
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func EnvironmentProduction(override bool) types.Option {
	return api.EnvironmentOption{
		APIGroup:       "kubernetes/v1",
		EnvPathSegment: "kubernetes",
		Override:       override,
	}
}
