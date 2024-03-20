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

	env, err := api.GetEnvironmentPathSegment(ctx, "kubernetes/v1", "kubernetes")
	if err != nil {
		return nil, fmt.Errorf("get environment: %w", err)
	}

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
