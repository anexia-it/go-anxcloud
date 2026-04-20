package acl

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const path = "/api/LBaaS/v1/ACL.json"

type ACLInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type ACL struct {
	CustomerIdentifier string                    `json:"customer_identifier"`
	ResellerIdentifier string                    `json:"reseller_identifier"`
	Identifier         string                    `json:"identifier"`
	Name               string                    `json:"name"`
	ParentType         string                    `json:"parent_type"`
	Frontend           *genericResource.Identity `json:"frontend"`
	Backend            *genericResource.Identity `json:"backend"`
	Criterion          string                    `json:"criterion"`
	Index              int                       `json:"index"`
	Value              string                    `json:"value"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "ACL"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (ACL, error) {
	name := "ACL"
	object, err := genericResource.GenericGetByID[ACL](ctx, identifier, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (ACL, error) {
	name := "ACL"

	object, err := genericResource.GenericCreate[ACL, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (ACL, error) {
	name := "ACL"
	object, err := genericResource.GenericUpdate[ACL, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "ACL"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
