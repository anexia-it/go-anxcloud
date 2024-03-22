package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Server and the identifier of the given Server.
// It implements the api.Object interface on *Server, making it usable with the generic API client.
func (s *Server) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/LBaaS/v1/server.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)
		if s.Backend.Identifier != "" {
			filters.Add("backend", s.Backend.Identifier)
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}
