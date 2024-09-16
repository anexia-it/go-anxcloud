package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
)

func (i *IPs) GetIdentifier(ctx context.Context) (string, error) {
	return "", nil
}

func (i *IPs) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op != apiTypes.OperationList {
		return nil, api.ErrOperationNotSupported
	}

	//return url.Parse("/api/vsphere/v1/provisioning/ips.json/" + i.LocationIdentifier + "/" + i.VLANIdentifier)
	return url.Parse("/api/vsphere/v1/provisioning/ips.json/LOCething1234567aaaaaaaabbbbbbbb/VLANthing1234567aaaaaaaabbbbbbbb")
}

func (i *IPs) HasPagination(ctx context.Context) (bool, error) {
	return true, nil
}

//// FilterRequestURL removes the Identifier from URL on Get operations (template needs to be parsed from list response)
//func (i *IPs) FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error) {
//	op, err := apiTypes.OperationFromContext(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	if op == apiTypes.OperationGet {
//		url.Path = path.Dir(url.Path)
//	}
//
//	q := url.Query()
//	q.Set("page", "1")
//	q.Set("limit", "1000")
//	url.RawQuery = q.Encode()
//
//	return url, nil
//}

func (i *IPs) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	err := json.NewDecoder(data).Decode(&i)
	if err != nil {
		return err
	}

	return nil
}

// TODO: remove
func (i *IPs) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	buf := bytes.NewBuffer(nil)
	tee := io.TeeReader(res.Body, buf)
	body, _ := io.ReadAll(tee)
	fmt.Printf("response: %+v\n", string(body))
	res.Body = io.NopCloser(buf)
	return res, nil
}
