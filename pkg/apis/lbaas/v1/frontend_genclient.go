package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
func (f *Frontend) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI("/api/LBaaS/v1/frontend.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)

		if f.LoadBalancer != nil && f.LoadBalancer.Identifier != "" {
			filters.Add("load_balancer", f.LoadBalancer.Identifier)
		}

		if f.Mode != "" {
			filters.Add("mode", string(f.Mode))
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}
