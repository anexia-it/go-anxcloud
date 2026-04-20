package loadbalancer

import (
	"context"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

type (
	// RuleInfo holds the name and identifier of a rule.
	RuleInfo = v1.RuleInfo

	// Loadbalancer holds the information of a load balancer instance.
	Loadbalancer = v1.LoadBalancer
)

// LoadBalancerInfo holds the identifier and the name of a load balancer
type LoadBalancerInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

const (
	path = "api/LBaaS/v1/loadbalancer.json"
)

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {
	name := "Loadbalancer"
	return genericResource.GetPagedGeneric(ctx, page, limit, a.client, name, path)
}

func (a api) GetByID(ctx context.Context, identifier string) (Loadbalancer, error) {
	name := "Loadbalancer"
	rule, err := genericResource.GenericGetByID[Loadbalancer](ctx, identifier, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *rule, err
}

func (a api) Create(ctx context.Context, definition Definition) (Loadbalancer, error) {
	name := "Loadbalancer"

	rule, err := genericResource.GenericCreate[Loadbalancer, Definition](ctx, definition, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *rule, err
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Loadbalancer, error) {
	name := "Loadbalancer"
	rule, err := genericResource.GenericUpdate[Loadbalancer, Definition](ctx, identifier, definition, a.client, name, path)
	if err != nil {
		return Loadbalancer{}, err
	}
	return *rule, err
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	name := "Loadbalancer"
	return genericResource.GenericDelete(ctx, identifier, a.client, name, path)
}
