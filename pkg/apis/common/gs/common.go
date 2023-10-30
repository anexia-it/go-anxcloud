package gs

import (
	"bytes"
	"context"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
	"io"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// GenericService has methods overridden for all GS objects in the same way
type GenericService struct{}

// FilterAPIResponse replaces the API response with an empty json object
func (gs *GenericService) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationDestroy {
		err = res.Body.Close()
		if err != nil {
			return nil, err
		}

		res.Body = io.NopCloser(bytes.NewReader([]byte("{}")))
		return res, nil
	}
	return res, nil
}

// RequestBody prevents decoding of delete responses as they are not compatible with the
// objects type
func RequestBody(ctx context.Context, br func() interface{}) (interface{}, error) {
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

// EndpointURL is a helper function which can be wrapped by API bindings to enable the filter helper
func EndpointURL(ctx context.Context, obj types.Object, resourcePath string) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(resourcePath)
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		helper, err := filter.NewHelper(obj)
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
