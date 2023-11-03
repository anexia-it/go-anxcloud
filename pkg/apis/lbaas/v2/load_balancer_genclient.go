package v2

import (
	"context"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"net/url"
)

// FilterAPIRequestBody adds the CommonRequestBody
func (lb *LoadBalancer) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return gs.RequestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			LoadBalancer
		}{
			LoadBalancer: *lb,
		}
	})
}

// EndpointURL returns the common URL for operations on LoadBalancer resource
func (lb *LoadBalancer) EndpointURL(ctx context.Context) (*url.URL, error) {
	return gs.EndpointURL(ctx, lb, "/api/LBaaSv2/v1/load_balancers.json")
}
