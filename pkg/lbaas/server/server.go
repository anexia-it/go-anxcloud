package server

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const (
	path = "/api/LBaaS/v1/server.json"
)

type Server = v1.Server

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "Server"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Server, error) {
	name := "Server"
	object, err := genericResource.GenericGetByID[Server](ctx, identifier, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Server, error) {
	name := "Server"

	object, err := genericResource.GenericCreate[Server, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Server, error) {
	name := "Server"
	object, err := genericResource.GenericUpdate[Server, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Server{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Server"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
