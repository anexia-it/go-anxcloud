package v2

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Region and the identifier of the given Region.
// It implements the api.Object interface on *Region, making it usable with the generic API client.
func (r *Region) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/region")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "name,state,description,backend,reseller,customer")

		filters := make(url.Values)

		if r.State != nil && r.State.ID != "" {
			filters.Add("state", r.State.ID)
		}

		if r.Backend != nil && r.Backend.Identifier != "" {
			filters.Add("backend", r.Backend.Identifier)
		}

		if r.CustomerIdentifier != "" {
			filters.Add("customer", r.CustomerIdentifier)
		}

		if r.ResellerIdentifier != "" {
			filters.Add("reseller", r.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for Regions, replacing linked Objects with just their identifier.
func (r *Region) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	reqBody := requestBody(ctx, func() interface{} {
		body := &struct {
			Region
			Backend *string `json:"backend,omitempty"`
		}{
			Region: *r,
		}

		if r.Backend != nil {
			backendID := r.Backend.Identifier
			body.Backend = &backendID
		}

		return body
	})
	return reqBody, nil
}
