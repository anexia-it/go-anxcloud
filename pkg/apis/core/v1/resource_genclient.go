package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-logr/logr"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (r Resource) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.ParseRequestURI("/api/core/v1/resource.json")
	if err != nil {
		return nil, err
	}

	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	switch op {
	// OperationCreate is not supported because the API does not exist in the engine.
	// OperationDestroy and OperationUpdate is not yet implemented
	case types.OperationCreate, types.OperationDestroy, types.OperationUpdate:
		return nil, api.ErrOperationNotSupported
	}

	if op == types.OperationList {
		query := u.Query()

		if len(r.Tags) > 1 {
			logr.FromContextOrDiscard(ctx).Info("Listing with multiple tags isn't supported. Only first one used")
		}

		if len(r.Tags) > 0 {
			query.Add("tag_name", r.Tags[0])
		}
		u.RawQuery = query.Encode()
	}
	return u, err
}

func (r *Resource) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return err
	}

	if op != types.OperationGet {
		return json.NewDecoder(data).Decode(r)
	} else {
		type apiResource struct {
			*Resource
			Tags []struct {
				Name       string
				Identifier string
			} `json:"tags"`
		}

		res := apiResource{
			Resource: r,
		}

		if err := json.NewDecoder(data).Decode(&res); err != nil {
			return err
		}

		r.Tags = make([]string, 0, len(res.Tags))
		for _, tag := range res.Tags {
			r.Tags = append(r.Tags, tag.Name)
		}
	}

	return nil
}

func (rwt ResourceWithTag) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op != types.OperationCreate && op != types.OperationDestroy {
		return nil, fmt.Errorf("%w: ResourceWithTag only support Create and Destroy operations", api.ErrOperationNotSupported)
	}

	return url.Parse(fmt.Sprintf("/api/core/v1/resource.json/%v/tags/%v", rwt.Identifier, rwt.Tag))
}

func (rwt ResourceWithTag) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	endpointURL, err := rwt.EndpointURL(ctx)
	if err != nil {
		return nil, err
	}

	baseURL := strings.TrimSuffix(req.URL.String(), endpointURL.String())

	url := fmt.Sprintf("%v/api/core/v1/resource.json/%v/tags/%v", baseURL, rwt.Identifier, rwt.Tag)

	newRequest, err := http.NewRequest(req.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating modified Request: %w", err)
	}

	for h, vs := range req.Header {
		if h != "Content-Type" && h != "Content-Length" {
			for _, v := range vs {
				newRequest.Header.Add(h, v)
			}
		}
	}

	return newRequest, nil
}

func (rwt ResourceWithTag) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	if res.StatusCode == http.StatusOK {
		res.StatusCode = http.StatusNoContent
		res.Body.Close()
		res.Body = ioutil.NopCloser(&bytes.Buffer{})
	}

	return res, nil
}
