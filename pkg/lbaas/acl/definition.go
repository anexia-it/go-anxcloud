package acl

import (
	"go.anx.io/go-anxcloud/pkg/genericresource"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type ACL struct {
	CustomerIdentifier string                    `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                    `json:"reseller_identifier,omitempty"`
	Identifier         string                    `json:"identifier,omitempty"`
	Name               string                    `json:"name,omitempty"`
	ParentType         string                    `json:"parent_type,omitempty"`
	Frontend           *genericresource.Identity `json:"frontend,omitempty"`
	Backend            *genericresource.Identity `json:"backend,omitempty"`
	Criterion          string                    `json:"criterion,omitempty"`
	Index              int                       `json:"index,omitempty"`
	Value              string                    `json:"value,omitempty"`
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
