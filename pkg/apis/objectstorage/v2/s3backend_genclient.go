package v2

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type S3Backend and the identifier of the given S3Backend.
// It implements the api.Object interface on *S3Backend, making it usable with the generic API client.
func (s *S3Backend) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/s3_backend")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "name,state,endpoint,backend_type,enabled,backend_user,backend_password,reseller,customer")

		filters := make(url.Values)

		if s.State != nil && s.State.ID != "" {
			filters.Add("state", s.State.ID)
		}

		if s.Endpoint.Identifier != "" {
			filters.Add("endpoint", s.Endpoint.Identifier)
		}

		if s.BackendType != nil && s.BackendType.Identifier != "" {
			filters.Add("backend_type", s.BackendType.Identifier)
		}

		if s.CustomerIdentifier != "" {
			filters.Add("customer", s.CustomerIdentifier)
		}

		if s.ResellerIdentifier != "" {
			filters.Add("reseller", s.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for S3Backends, replacing linked Objects with just their identifier.
func (s *S3Backend) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	body := requestBody(ctx, func() interface{} {
		return &struct {
			S3Backend
			Endpoint string `json:"endpoint"`
		}{
			S3Backend: *s,
			Endpoint:  s.Endpoint.Identifier,
		}
	})
	return body, nil
}
