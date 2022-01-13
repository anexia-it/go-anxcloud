package v1

import (
	"bytes"
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"io"
	"net/http"
	"net/url"
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

// FilterAPIRequestBody generates the request body for creating a new Server, which differs from the Server object.
func (s *Server) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate {
		return struct {
			Server
			State   State  `json:"state"`
			Backend string `json:"backend"`
		}{
			Server:  *s,
			State:   NewlyCreated,
			Backend: s.Backend.Identifier,
		}, nil
	} else if op == types.OperationUpdate {
		return struct {
			Server
			State   State  `json:"state"`
			Backend string `json:"backend"`
		}{
			Server:  *s,
			Backend: s.Backend.Identifier,
			State:   Updating,
		}, nil
	}

	return s, nil
}

func (s *Server) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return res, err
	}

	if op == types.OperationDestroy {
		err = res.Body.Close()
		if err != nil {
			return res, err
		}

		res.Body = io.NopCloser(bytes.NewReader([]byte("{}")))
		return res, nil
	}
	return res, nil
}
