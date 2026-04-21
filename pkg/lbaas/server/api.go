package server

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const (
	path = "/api/LBaaS/v1/server.json"
)

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend server API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Server, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Server"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Server, error) {
	name := "Server"
	object, err := genericresource.GenericGetByID[Server](ctx, identifier, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Server, error) {
	name := "Server"

	object, err := genericresource.GenericCreate[Server, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Server, error) {
	name := "Server"
	object, err := genericresource.GenericUpdate[Server, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Server"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
