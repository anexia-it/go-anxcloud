package backend

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const (
	path = "api/LBaaS/v1/backend.json"
)

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Backend, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Backend"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Backend, error) {
	name := "Backend"
	object, err := genericresource.GenericGetByID[Backend](ctx, identifier, a.client, name, path)
	if err != nil {
		return Backend{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Backend, error) {
	name := "Backend"

	object, err := genericresource.GenericCreate[Backend, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Backend{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Backend, error) {
	name := "Backend"
	object, err := genericresource.GenericUpdate[Backend, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Backend{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Backend"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
