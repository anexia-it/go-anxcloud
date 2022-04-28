package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
func (f *Frontend) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI("/api/LBaaS/v1/frontend.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)

		if f.LoadBalancer != nil && f.LoadBalancer.Identifier != "" {
			filters.Add("load_balancer", f.LoadBalancer.Identifier)
		}

		if f.Mode != "" {
			filters.Add("mode", string(f.Mode))
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}

// FilterAPIRequestBody generates the request body for creating a new Frontend, which differs from the Frontend object.
func (f *Frontend) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		return struct {
			Frontend
			LoadBalancer   string `json:"load_balancer"`
			DefaultBackend string `json:"default_backend"`

			// we never want to send the state field, so making sure to omit it here
			State string `json:"state,omitempty"`
		}{
			Frontend:       *f,
			LoadBalancer:   f.LoadBalancer.Identifier,
			DefaultBackend: f.DefaultBackend.Identifier,
		}, nil
	}

	return f, nil
}

func (f *Frontend) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
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
