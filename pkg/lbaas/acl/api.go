package acl

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const path = "/api/LBaaS/v1/ACL.json"

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) genericresource.API[ACL, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "ACL"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (ACL, error) {
	name := "ACL"
	object, err := genericresource.GenericGetByID[ACL](ctx, identifier, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (ACL, error) {
	name := "ACL"

	object, err := genericresource.GenericCreate[ACL, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (ACL, error) {
	name := "ACL"
	object, err := genericresource.GenericUpdate[ACL, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return ACL{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "ACL"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
