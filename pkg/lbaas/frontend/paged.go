package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/pagination"
	"net/http"
	"net/url"
	"strconv"
)

type FrontendPage struct {
	Page        int            `json:"page"`
	TotalItems  int            `json:"total_items"`
	TotalPages  int            `json:"total_pages"`
	Limit       int            `json:"limit"`
	Data        []FrontendInfo `json:"data"`
	pageOptions []pagination.Option
}

func (f FrontendPage) Options() []pagination.Option {
	return f.pageOptions
}

func (f FrontendPage) Num() int {
	return f.Page
}

func (f FrontendPage) Size() int {
	return f.Limit
}

func (f FrontendPage) Total() int {
	return f.TotalPages
}

func (a api) GetPage(ctx context.Context, page, limit int, opts ...pagination.Option) (pagination.Page, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path
	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))
	if option := pagination.GetOption(opts, common.OptNameSearch); option != nil {
		query.Set("search", option.Value)
	}
	if option := pagination.GetOption(opts, common.OptNameFilter); option != nil {
		query.Set("filter", option.Value)
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

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer frontends %s", response.Status)
	}

	payload := struct {
		Page FrontendPage `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer frontend list response: %w", err)
	}

	payload.Page.pageOptions = opts
	return payload.Page, nil
}

func (a api) NextPage(ctx context.Context, page pagination.Page) (pagination.Page, error) {
	return a.GetPage(ctx, page.Num()+1, page.Size(), page.Options()...)
}

func (f FrontendPage) Content() interface{} {
	return f.Data
}
