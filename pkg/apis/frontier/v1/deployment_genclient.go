package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the base URL path of the resources API
// Deployments are created via the (frontier) api API
func (d *Deployment) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if op == types.OperationCreate {
		return url.Parse(fmt.Sprintf("/api/frontier/v1/api.json/%s/deploy", d.APIIdentifier))
	}

	return url.Parse("/api/frontier/v1/deployment.json")
}

// FilterAPIRequestBody is a hook that handles the deviating request format expected when creating
// a Deployment
func (d *Deployment) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate {
		return &struct {
			Slug string `json:"slug"`
		}{d.Slug}, nil
	}

	return d, nil
}

// DecodeAPIResponse is a hook that extracts the identifier from the deviating response
// format when creating a Deployment
func (d *Deployment) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return err
	}

	if op == types.OperationCreate {
		var frontierAPI API
		if err := json.NewDecoder(data).Decode(&frontierAPI); err != nil {
			return err
		}

		d.Identifier = frontierAPI.DeploymentIdentifier
		return nil
	}

	return json.NewDecoder(data).Decode(d)
}
