package loadbalancer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/pagination"
	"github.com/anexia-it/go-anxcloud/pkg/utils/param"
	"net/http"
	"net/url"
	"strconv"
)

type LoadBalancerPage struct {
	Page        int                `json:"page"`
	TotalItems  int                `json:"total_items"`
	TotalPages  int                `json:"total_pages"`
	Limit       int                `json:"limit"`
	Data        []LoadBalancerInfo `json:"data"`
	pageOptions []param.Parameter
}

func (f LoadBalancerPage) Options() []param.Parameter {
	return f.pageOptions
}

func (f LoadBalancerPage) Num() int {
	return f.Page
}

func (f LoadBalancerPage) Size() int {
	return f.Limit
}

func (f LoadBalancerPage) Total() int {
	return f.TotalPages
}

func (a api) GetPage(ctx context.Context, page, limit int, parameters ...param.Parameter) (pagination.Page, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path
	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))
	for _, parameter := range parameters {
		parameter(query)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer frontends %s", response.Status)
	}

	payload := struct {
		Page LoadBalancerPage `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer frontend list response: %w", err)
	}

	payload.Page.pageOptions = parameters
	return payload.Page, nil
}

func (a api) NextPage(ctx context.Context, page pagination.Page) (pagination.Page, error) {
	return a.GetPage(ctx, page.Num()+1, page.Size(), page.Options()...)
}

func (f LoadBalancerPage) Content() interface{} {
	return f.Data
}
