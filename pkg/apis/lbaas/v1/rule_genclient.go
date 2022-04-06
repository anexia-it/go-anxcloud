package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Rule and the identifier of the given Rule.
// It implements the api.Object interface on *Rule, making it usable with the generic API client.
func (r *Rule) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/LBaaS/v1/rule.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := r.buildListFilters()
		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}

func (r *Rule) buildListFilters() url.Values {
	filters := make(url.Values)

	if r.RuleType != "" {
		filters.Add("rule_type", r.RuleType)
	}
	if r.ParentType != "" {
		filters.Add("parent_type", r.ParentType)
	}
	if r.Frontend.Identifier != "" {
		filters.Add("frontend", r.Frontend.Identifier)
	}
	if r.Backend.Identifier != "" {
		filters.Add("backend", r.Backend.Identifier)
	}
	if r.Condition != "" {
		filters.Add("condition", r.Condition)
	}
	if r.Type != "" {
		filters.Add("type", r.Type)
	}
	if r.Action != "" {
		filters.Add("action", r.Action)
	}
	if r.RedirectionType != "" {
		filters.Add("redirection_type", r.RedirectionType)
	}
	if r.RedirectionCode != "" {
		filters.Add("redirection_code", r.RedirectionCode)
	}

	return filters
}

// FilterAPIRequestBody generates the request body for Rules, replacing linked Objects with just their identifier.
func (r *Rule) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return requestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			Rule
			Backend  string `json:"backend,omitempty"`
			Frontend string `json:"frontend,omitempty"`
		}{
			Rule:     *r,
			Backend:  r.Backend.Identifier,
			Frontend: r.Frontend.Identifier,
		}
	})
}
