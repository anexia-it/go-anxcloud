package frontend

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const path = "api/LBaaS/v1/frontend.json"

// Frontend represents a LBaaS Frontend.
type Frontend = v1.Frontend

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "Frontend"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Frontend, error) {
	name := "Frontend"
	object, err := genericResource.GenericGetByID[Frontend](ctx, identifier, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Frontend, error) {
	name := "Frontend"

	object, err := genericResource.GenericCreate[Frontend, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Frontend, error) {
	name := "Frontend"
	object, err := genericResource.GenericUpdate[Frontend, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Frontend{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Frontend"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
