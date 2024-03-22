package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Backend and the identifier of the given Backend.
// It implements the api.Object interface on *Backend, making it usable with the generic API client.
func (b *Backend) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/LBaaS/v1/backend.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)

		if b.LoadBalancer.Identifier != "" {
			filters.Add("load_balancer", b.LoadBalancer.Identifier)
		}

		if b.Mode != "" {
			filters.Add("mode", string(b.Mode))
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, err
}
