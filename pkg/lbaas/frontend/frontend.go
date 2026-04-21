package frontend

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const path = "api/LBaaS/v1/frontend.json"

// Frontend represents a LBaaS Frontend.
type Frontend = v1.Frontend

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Frontend"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Frontend, error) {
	name := "Frontend"
	object, err := genericresource.GenericGetByID[Frontend](ctx, identifier, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Frontend, error) {
	name := "Frontend"

	object, err := genericresource.GenericCreate[Frontend, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Frontend, error) {
	name := "Frontend"
	object, err := genericresource.GenericUpdate[Frontend, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Frontend"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
