package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
)

func (p *ProvisionProgress) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op != apiTypes.OperationGet {
		return nil, api.ErrOperationNotSupported
	}

	return url.Parse("/api/vsphere/v1/provisioning/progress.json")
}

// TODO: remove
func (p *ProvisionProgress) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	buf := bytes.NewBuffer(nil)
	tee := io.TeeReader(res.Body, buf)
	body, _ := io.ReadAll(tee)
	fmt.Printf("response: %+v\n", string(body))
	res.Body = io.NopCloser(buf)
	return res, nil
}
