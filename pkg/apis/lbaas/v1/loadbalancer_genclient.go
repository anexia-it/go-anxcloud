package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the URL where to retrieve objects of type LoadBalancer and the identifier of the given Loadbalancer.
// It implements the api.Object interface on *LoadBalancer, making it usable with the generic API client.
func (lb *LoadBalancer) EndpointURL(ctx context.Context) (*url.URL, error) {
	url, err := url.ParseRequestURI("/api/LBaaS/v1/loadbalancer.json")
	return url, err
}
