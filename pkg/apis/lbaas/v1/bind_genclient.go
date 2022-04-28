package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (b *Bind) EndpointURL(ctx context.Context) (*url.URL, error) {

	// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
	// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("api/LBaaS/v1/bind.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)
		if b.Frontend.Identifier != "" {
			filters.Add("frontend", b.Frontend.Identifier)
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}

// FilterAPIRequestBody generates the request body for Binds, replacing linked Objects with just their identifier.
func (b *Bind) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return requestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			Bind
			Frontend string `json:"frontend"`
		}{
			Bind:     *b,
			Frontend: b.Frontend.Identifier,
		}
	})
}
