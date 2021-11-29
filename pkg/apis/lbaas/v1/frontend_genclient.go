package v1

import (
	"context"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
func (f *Frontend) EndpointURL(ctx context.Context) (*url.URL, error) {
	url, err := url.ParseRequestURI("/api/LBaaS/v1/frontend.json")
	return url, err
}

// FilterAPIRequestBody generates the request body for creating a new Frontend, which differs from the Frontend object.
func (f *Frontend) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate {
		return struct {
			Name           string `json:"name"`
			LoadBalancer   string `json:"load_balancer"`
			DefaultBackend string `json:"default_backend"`
			State          State  `json:"state"`
		}{
			Name:           f.Name,
			LoadBalancer:   f.LoadBalancer.Identifier,
			DefaultBackend: f.DefaultBackend.Identifier,
			State:          NewlyCreated,
		}, nil
	} else if op == types.OperationUpdate {
		return struct {
			Frontend
			LoadBalancer   string `json:"load_balancer"`
			DefaultBackend string `json:"default_backend"`
		}{
			Frontend:       *f,
			LoadBalancer:   f.LoadBalancer.Identifier,
			DefaultBackend: f.DefaultBackend.Identifier,
		}, nil
	}

	return f, nil
}
