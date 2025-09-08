package v2

import (
	"context"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Bucket and the identifier of the given Bucket.
// It implements the api.Object interface on *Bucket, making it usable with the generic API client.
func (b *Bucket) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/bucket")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "name,state,region,object_count,object_size,backend,tenant,reseller,customer")

		filters := make(url.Values)

		if b.State != nil && b.State.ID != "" {
			filters.Add("state", b.State.ID)
		}

		if b.Region.Identifier != "" {
			filters.Add("region", b.Region.Identifier)
		}

		if b.Backend.Identifier != "" {
			filters.Add("backend", b.Backend.Identifier)
		}

		if b.Tenant.Identifier != "" {
			filters.Add("tenant", b.Tenant.Identifier)
		}

		if b.CustomerIdentifier != "" {
			filters.Add("customer", b.CustomerIdentifier)
		}

		if b.ResellerIdentifier != "" {
			filters.Add("reseller", b.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for Buckets, replacing linked Objects with just their identifier.
// Only includes required fields and conditionally includes optional fields to avoid 422 errors.
func (b *Bucket) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	body := requestBody(ctx, func() interface{} {
		// Create minimal request body with only required fields
		reqBody := &struct {
			Name               string  `json:"name"`
			Region             string  `json:"region"`
			Backend            string  `json:"backend"`
			Tenant             string  `json:"tenant"`
			CustomerIdentifier string  `json:"customer_identifier"`           // Always include (can be empty)
			ResellerIdentifier *string `json:"reseller_identifier,omitempty"` // Omit if empty
			Share              *bool   `json:"share,omitempty"`               // Omit if false
		}{
			Name:               b.Name,
			Region:             b.Region.Identifier,
			Backend:            b.Backend.Identifier,
			Tenant:             b.Tenant.Identifier,
			CustomerIdentifier: b.CustomerIdentifier, // Always include (even if empty)
		}

		// Conditionally include optional fields only if they have non-default values
		if b.ResellerIdentifier != "" {
			reqBody.ResellerIdentifier = &b.ResellerIdentifier
		}

		if b.Share {
			reqBody.Share = &b.Share
		}

		return reqBody
	})
	return body, nil
}

// FilterAPIRequest modifies the HTTP method for UPDATE operations (Object Storage API uses PATCH, not PUT)
func (b *Bucket) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// For UPDATE operations, use PATCH method instead of PUT
	if op == types.OperationUpdate {
		req.Method = "PATCH"
	}

	return req, nil
}
