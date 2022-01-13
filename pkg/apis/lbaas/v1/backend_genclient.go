package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
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

// FilterAPIRequestBody generates the request body for creating a new Backend, which differs from the Backend object.
func (b *Backend) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate {
		return struct {
			Name         string `json:"name"`
			LoadBalancer string `json:"load_balancer"`
			Mode         Mode   `json:"mode"`
			State        State  `json:"state"`
		}{
			Name:         b.Name,
			Mode:         b.Mode,
			LoadBalancer: b.LoadBalancer.Identifier,
			State:        NewlyCreated,
		}, nil
	} else if op == types.OperationUpdate {
		return struct {
			Backend
			LoadBalancer string `json:"load_balancer"`
		}{
			Backend:      *b,
			LoadBalancer: b.LoadBalancer.Identifier,
		}, nil
	}

	return b, nil
}

func (b *Backend) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
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
