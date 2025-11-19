package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/pagination"
	"go.anx.io/go-anxcloud/pkg/utils/param"
)

type ClusterPage struct {
	Page        int           `json:"page"`
	TotalItems  int           `json:"total_items"`
	TotalPages  int           `json:"total_pages"`
	Limit       int           `json:"limit"`
	Data        []ClusterInfo `json:"data"`
	pageOptions []param.Parameter
}

func (f ClusterPage) Options() []param.Parameter {
	return f.pageOptions
}

func (f ClusterPage) Num() int {
	return f.Page
}

func (f ClusterPage) Size() int {
	return f.Limit
}

func (f ClusterPage) Total() int {
	return f.TotalPages
}

func (f ClusterPage) Content() interface{} {
	return f.Data
}

func (a *api) GetPage(ctx context.Context, page, limit int, parameters ...param.Parameter) (pagination.Page, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)
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
		return nil, fmt.Errorf("could not get kubernetes clusters %s", response.Status)
	}

	payload := struct {
		Page ClusterPage `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes cluster list response: %w", err)
	}

	payload.Page.pageOptions = parameters
	return payload.Page, nil
}

func (a *api) NextPage(ctx context.Context, page pagination.Page) (pagination.Page, error) {
	return a.GetPage(ctx, page.Num()+1, page.Size(), page.Options()...)
}
