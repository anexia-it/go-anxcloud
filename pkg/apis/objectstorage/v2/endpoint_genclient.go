package v2

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// EndpointURL returns the URL where to retrieve objects of type Endpoint and the identifier of the given Endpoint.
// It implements the api.Object interface on *Endpoint, making it usable with the generic API client.
func (e *Endpoint) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/endpoint")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "url,state,endpoint_user,endpoint_password,enabled,reseller,customer")

		filters := make(url.Values)

		if e.State != nil && e.State.ID != "" {
			filters.Add("state", e.State.ID)
		}

		if e.EndpointUser != "" {
			filters.Add("endpoint_user", e.EndpointUser)
		}

		if e.CustomerIdentifier != "" {
			filters.Add("customer", e.CustomerIdentifier)
		}

		if e.ResellerIdentifier != "" {
			filters.Add("reseller", e.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for Endpoints, replacing linked Objects with just their identifier.
func (e *Endpoint) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	body := requestBody(ctx, func() interface{} {
		return &struct {
			*Endpoint
			// Exclude conflicting fields from commonRequestBody
			Tags     gs.PartialResourceList `json:"tags,omitempty"`
			Reseller string                 `json:"reseller,omitempty"`
			Customer string                 `json:"customer,omitempty"`
			Share    bool                   `json:"share,omitempty"`
		}{
			Endpoint: e,
			Tags:     e.Tags,
			Reseller: e.Reseller,
			Customer: e.Customer,
			Share:    e.Share,
		}
	})
	return body, nil
}
