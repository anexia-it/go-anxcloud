package pagination

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	path2 "path"
	"reflect"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/param"
)

var (
	ErrConditionNeverMet = errors.New("looped all items and the condition was never met")
)

type Page interface {
	Num() int
	Size() int
	Total() int
	Options() []param.Parameter
	Content() interface{}
}

// Pageable should be implemented by evey struct that supports pagination.
type Pageable interface {
	GetPage(ctx context.Context, page, limit int, opts ...param.Parameter) (Page, error)
	NextPage(ctx context.Context, page Page) (Page, error)
}

// HasNext is a helper function which checks whether there are more pages to fetch
func HasNext(page Page) bool {
	return page.Num() < page.Total()
}

type UntilTrueFunc func(interface{}) (bool, error)

// LoopUntil takes a pageable and loops over it until untilFunc returns true or an error.
func LoopUntil(ctx context.Context, pageable Pageable, untilFunc UntilTrueFunc, opts ...param.Parameter) error {
	page, err := pageable.GetPage(ctx, 1, 10, opts...)
	if err != nil {
		return err
	}

	content := reflect.ValueOf(page.Content())
	switch content.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		stopped := false
		for i := 0; i < content.Len(); i++ {
			content.Interface()
			stop, err := untilFunc(content.Index(i).Interface())
			if err != nil {
				return err
			}
			if stop {
				stopped = true
				break
			}
		}
		if !stopped {
			return ErrConditionNeverMet
		}
	default:
		panic(fmt.Sprintf("The page content is supposed to be of type slice or array but was %T", page.Content()))
	}
	return nil
}

type CancelFunc func()

// AsChan takes a Pageable and returns its Pageable.Content via a channel until there are no more pages or
// CancelFunc gets called by the consumer.
func AsChan(ctx context.Context, pageable Pageable, opts ...param.Parameter) (chan interface{}, CancelFunc) {
	consumer := make(chan interface{})
	done := make(chan interface{})
	cancel := func() {
		close(done)
	}

	go func() {
		defer close(consumer)
		_ = LoopUntil(ctx, pageable, func(i interface{}) (bool, error) {
			select {
			case consumer <- i:
			case _, ok := <-done:
				if !ok {
					return true, nil
				}
			}
			return false, nil
		}, opts...)
	}()
	return consumer, cancel
}

func GetPage(ctx context.Context, page int, limit int, parameters []param.Parameter, client client.Client, path string) (Page, error) {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, path)
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

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer frontends %s", response.Status)
	}

	payload := struct {
		Page GenericPage `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer frontend list response: %w", err)
	}

	payload.Page.PageOptions = parameters
	return payload.Page, nil
}

type Identifiable interface {
	GetIdentifier() string
	GetName() string
}

type Identity struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (g Identity) GetIdentifier() string {
	return g.Identifier
}

func (g Identity) GetName() string {
	return g.Name
}

type GenericPage struct {
	Page        int        `json:"page"`
	TotalItems  int        `json:"total_items"`
	TotalPages  int        `json:"total_pages"`
	Limit       int        `json:"limit"`
	Data        []Identity `json:"data"`
	PageOptions []param.Parameter
}

func (f GenericPage) Options() []param.Parameter {
	return f.PageOptions
}

func (f GenericPage) Num() int {
	return f.Page
}

func (f GenericPage) Size() int {
	return f.Limit
}

func (f GenericPage) Total() int {
	return f.TotalPages
}

func (f GenericPage) Content() interface{} {
	return f.Data
}
