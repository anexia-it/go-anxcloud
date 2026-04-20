package bind

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const (
	path = "api/LBaaS/v1/bind.json"
)

type Bind = v1.Bind

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "Bind"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Bind, error) {
	name := "Bind"
	object, err := genericResource.GenericGetByID[Bind](ctx, identifier, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Bind, error) {
	name := "Bind"

	object, err := genericResource.GenericCreate[Bind, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Bind, error) {
	name := "Bind"
	object, err := genericResource.GenericUpdate[Bind, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Bind"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
