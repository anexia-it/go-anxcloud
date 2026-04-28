package acl

import (
	"go.anx.io/go-anxcloud/pkg/genericresource"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type ACL struct {
	CustomerIdentifier string                    `json:"customer_identifier"`
	ResellerIdentifier string                    `json:"reseller_identifier"`
	Identifier         string                    `json:"identifier"`
	Name               string                    `json:"name"`
	ParentType         string                    `json:"parent_type"`
	Frontend           *genericresource.Identity `json:"frontend"`
	Backend            *genericresource.Identity `json:"backend"`
	Criterion          string                    `json:"criterion"`
	Index              int                       `json:"index"`
	Value              string                    `json:"value"`
}

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
