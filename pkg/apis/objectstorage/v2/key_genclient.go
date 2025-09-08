package v2

import (
	"context"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Key and the identifier of the given Key.
// It implements the api.Object interface on *Key, making it usable with the generic API client.
func (k *Key) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/key")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "name,remote_id,state,backend,tenant,user,expire_date,reseller,customer,secret,secret_url")

		filters := make(url.Values)

		if k.State != nil && k.State.ID != "" {
			filters.Add("state", k.State.ID)
		}

		if k.Backend != nil && k.Backend.Identifier != "" {
			filters.Add("backend", k.Backend.Identifier)
		}

		if k.Tenant != nil && k.Tenant.Identifier != "" {
			filters.Add("tenant", k.Tenant.Identifier)
		}

		if k.User != nil && k.User.Identifier != "" {
			filters.Add("user", k.User.Identifier)
		}

		if k.CustomerIdentifier != "" {
			filters.Add("customer", k.CustomerIdentifier)
		}

		if k.ResellerIdentifier != "" {
			filters.Add("reseller", k.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for Keys, replacing linked Objects with just their identifier.
// Only includes required fields and conditionally includes optional fields to avoid API errors.
func (k *Key) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	reqBody := requestBody(ctx, func() interface{} {
		// Create minimal request body with only required fields
		reqBody := &struct {
			Name               string  `json:"name"`
			Backend            string  `json:"backend"`
			Tenant             string  `json:"tenant"`
			User               string  `json:"user"`
			ExpireDate         *string `json:"expiry_date,omitempty"`
			CustomerIdentifier string  `json:"customer_identifier"`
			ResellerIdentifier *string `json:"reseller_identifier,omitempty"`
			Share              *bool   `json:"share,omitempty"`
		}{
			Name:               k.Name,
			CustomerIdentifier: k.CustomerIdentifier, // Always include (even if empty)
		}

		// Set required linked object identifiers
		if k.Backend != nil {
			reqBody.Backend = k.Backend.Identifier
		}
		if k.Tenant != nil {
			reqBody.Tenant = k.Tenant.Identifier
		}
		if k.User != nil {
			reqBody.User = k.User.Identifier
		}

		// Handle expiry date
		if k.ExpireDate != nil {
			expireStr := k.ExpireDate.Format("2006-01-02T15:04:05Z")
			reqBody.ExpireDate = &expireStr
		}

		// Conditionally include optional fields only if they have non-default values
		if k.ResellerIdentifier != "" {
			reqBody.ResellerIdentifier = &k.ResellerIdentifier
		}

		if k.Share {
			reqBody.Share = &k.Share
		}

		return reqBody
	})
	return reqBody, nil
}

// FilterAPIRequest modifies the HTTP method for UPDATE operations (Object Storage API uses PATCH, not PUT)
func (k *Key) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
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
