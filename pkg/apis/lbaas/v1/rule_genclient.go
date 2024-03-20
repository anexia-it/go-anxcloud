package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
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
		helper, err := filter.NewHelper(r)
		if err != nil {
			return nil, err
		}

		filters := helper.BuildQuery().Encode()

		if filters != "" {
			query := u.Query()
			query.Set("filters", filters)
			u.RawQuery = query.Encode()
		}
	}

	return u, nil
}
