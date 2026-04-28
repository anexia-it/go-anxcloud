package rule

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const path = "/api/LBaaS/v1/rule.json"

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Rule, Definition] {
	return &api{c}
}

func (a api) Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error) {
	name := "Rule"
	return genericresource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Rule, error) {
	name := "Rule"
	rule, err := genericresource.GenericGetByID[Rule](ctx, identifier, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) Create(ctx context.Context, definition Definition) (Rule, error) {
	name := "Rule"

	rule, err := genericresource.GenericCreate[Rule, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Rule, error) {
	name := "Rule"
	rule, err := genericresource.GenericUpdate[Rule, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Rule"
	return genericresource.GenericDelete(ctx, identifier, a.client, name, path)
}
