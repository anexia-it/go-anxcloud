package loadbalancer

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const (
	path = "api/LBaaS/v1/loadbalancer.json"
)

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Loadbalancer, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Loadbalancer"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Loadbalancer, error) {
	name := "Loadbalancer"
	object, err := genericresource.GenericGetByID[Loadbalancer](ctx, identifier, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *object, err
}

func (a api) Create(ctx context.Context, definition Definition) (Loadbalancer, error) {
	name := "Loadbalancer"

	object, err := genericresource.GenericCreate[Loadbalancer, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *object, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Loadbalancer, error) {
	name := "Loadbalancer"
	object, err := genericresource.GenericUpdate[Loadbalancer, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *object, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Loadbalancer"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
