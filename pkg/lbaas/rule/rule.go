package rule

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

const path = "/api/LBaaS/v1/rule.json"

//func (a api) GetPath() string {
//	return path
//}

type Rule = v1.Rule

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
