package acl

import "github.com/anexia-it/go-anxcloud/pkg/lbaas/common"

// Definition describes the ACL object that should be created
type Definition struct {
	Name       string       `json:"name"`
	State      common.State `json:"state"`
	ParentType string       `json:"parent_type"`
	Criterion  string       `json:"criterion"`
	Index      string       `json:"index"`
	Value      string       `json:"value"`
}
