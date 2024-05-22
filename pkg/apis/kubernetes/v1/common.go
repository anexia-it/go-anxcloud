package v1

import (
	"context"
	"fmt"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
)

func endpointURL(ctx context.Context, o types.Object, resourcePathName string) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationUpdate {
		return nil, api.ErrOperationNotSupported
	}

	env := api.GetEnvironmentPathSegment(ctx, "kubernetes/v1", "kubernetes")

	// we can ignore the error since the URL is hard-coded known as valid
	u, _ := url.Parse(fmt.Sprintf("/api/%s/v1/%s.json", env, resourcePathName))

	if op == types.OperationList {
		helper, err := filter.NewHelper(o)
		if err != nil {
			return nil, err
		}

		filters := helper.BuildQuery().Encode()

		if filters != "" {
			query := u.Query()
			query.Set("filters", filters)
			u.RawQuery = query.Encode()
		}
	}

	return u, nil
}

// commonRequestBody is embedded in the request body types of kubernetes resources
type commonRequestBody struct {
	// this allows us to provide the state as string
	State string `json:"state,omitempty"`
}

// requestBody removes the request body for read operations on kubernetes resources
func requestBody(ctx context.Context, br func() interface{}) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		response := br()

		return response, nil
	}

	return nil, nil
}
