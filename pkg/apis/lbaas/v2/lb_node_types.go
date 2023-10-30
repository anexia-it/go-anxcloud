package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object

// Node holds the information of a load balancing node within a Cluster
type Node struct {
	gs.GenericService
	gs.HasState

	Identifier string                  `json:"identifier,omitempty" anxcloud:"identifier"`
	Name       string                  `json:"name,omitempty"`
	Cluster    *common.PartialResource `json:"cluster,omitempty" anxcloud:"filterable"`
}
