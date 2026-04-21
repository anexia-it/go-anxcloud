package bind

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/genericresource"

	"go.anx.io/go-anxcloud/pkg/client"
)

const (
	path = "api/LBaaS/v1/bind.json"
)

type api struct {
	client client.Client
}

// NewAPI creates a new bind API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Bind, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Bind"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Bind, error) {
	name := "Bind"
	object, err := genericresource.GenericGetByID[Bind](ctx, identifier, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Bind, error) {
	name := "Bind"

	object, err := genericresource.GenericCreate[Bind, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Bind, error) {
	name := "Bind"
	object, err := genericresource.GenericUpdate[Bind, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Bind{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Bind"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
