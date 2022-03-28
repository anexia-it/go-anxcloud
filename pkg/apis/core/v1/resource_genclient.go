package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/url"

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
