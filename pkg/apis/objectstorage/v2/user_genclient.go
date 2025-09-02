package v2

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type User and the identifier of the given User.
// It implements the api.Object interface on *User, making it usable with the generic API client.
func (user *User) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/user")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "user_name,state,enabled,full_name,backend,tenant,remote_id,reseller,customer")

		filters := make(url.Values)

		if user.State != nil && user.State.ID != "" {
			filters.Add("state", user.State.ID)
		}

		if user.Enabled != nil {
			filters.Add("enabled", strconv.FormatBool(*user.Enabled))
		}

		if user.FullName != "" {
			filters.Add("full_name", user.FullName)
		}

		if user.Backend.Identifier != "" {
			filters.Add("backend", user.Backend.Identifier)
		}

		if user.Tenant.Identifier != "" {
			filters.Add("tenant", user.Tenant.Identifier)
		}

		if user.CustomerIdentifier != "" {
			filters.Add("customer", user.CustomerIdentifier)
		}

		if user.ResellerIdentifier != "" {
			filters.Add("reseller", user.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for Users, replacing linked Objects with just their identifier.
func (user *User) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	body := requestBody(ctx, func() interface{} {
		return &struct {
			User
			Backend string `json:"backend"`
			Tenant  string `json:"tenant"`
		}{
			User:    *user,
			Backend: user.Backend.Identifier,
			Tenant:  user.Tenant.Identifier,
		}
	})
	return body, nil
}

// FilterAPIRequest modifies the HTTP method for UPDATE operations (Object Storage API uses PATCH, not PUT)
func (user *User) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
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
