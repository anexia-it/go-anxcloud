package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

type omitResponseDecodeOnDestroy struct{}

func (*omitResponseDecodeOnDestroy) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationDestroy && res.StatusCode == http.StatusOK {
		res.StatusCode = http.StatusNoContent
		res.Body.Close()
		res.Body = io.NopCloser(&bytes.Buffer{})
	}

	return res, nil
}
