package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the URL where to retrieve objects of type LoadBalancer and the identifier of the given Loadbalancer.
// It implements the api.Object interface on *LoadBalancer, making it usable with the generic API client.
func (lb *LoadBalancer) EndpointURL(ctx context.Context) (*url.URL, error) {
	url, err := url.ParseRequestURI("/api/LBaaS/v1/loadbalancer.json")
	return url, err
}

// FilterAPIRequestBody generates the request body for creating a new LoadBalancer, which differs from the LoadBalancer object.
func (lb *LoadBalancer) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate {
		return map[string]string{
			"name":       lb.Name,
			"ip_address": lb.IpAddress,
			"state":      "2",
		}, nil
	}

	return lb, nil
}

func (l *LoadBalancer) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
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
