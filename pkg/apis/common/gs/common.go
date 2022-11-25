package gs

import (
	"bytes"
	"context"
	"io"
	"net/http"

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
