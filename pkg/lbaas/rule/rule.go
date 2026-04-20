package rule

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const path = "/api/LBaaS/v1/rule.json"

//func (a api) GetPath() string {
//	return path
//}

type Rule = v1.Rule

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "Rule"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Rule, error) {
	name := "Rule"
	rule, err := genericResource.GenericGetByID[Rule](ctx, identifier, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) Create(ctx context.Context, definition Definition) (Rule, error) {
	name := "Rule"

	rule, err := genericResource.GenericCreate[Rule, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Rule, error) {
	name := "Rule"
	rule, err := genericResource.GenericUpdate[Rule, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Rule{}, err
	}
	return *rule, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Rule"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
