package acl

import "go.anx.io/go-anxcloud/pkg/lbaas/common"

// Definition describes the ACL object that should be created
type Definition struct {
	Name       string       `json:"name"`
	State      common.State `json:"state"`
	ParentType string       `json:"parent_type"`
	Criterion  string       `json:"criterion"`
	Index      int          `json:"index"`
	Value      string       `json:"value"`
	Frontend   *string      `json:"frontend,omitempty"`
	Backend    *string      `json:"backend,omitempty"`
}
