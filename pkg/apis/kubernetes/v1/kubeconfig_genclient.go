package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

const (
	requesetKubeConfigRuleIdentifier = "12277a581e1c47cba72338425a008aa3"
	removeKubeConfigRuleIdentifier   = "eec87131729e44fa91a4b7ee8c365a26"
)

func (k *kubeconfig) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var ruleIdentifier string

	switch op {
	case types.OperationCreate:
		ruleIdentifier = requesetKubeConfigRuleIdentifier
	case types.OperationDestroy:
		ruleIdentifier = removeKubeConfigRuleIdentifier
	default:
		return nil, api.ErrOperationNotSupported
	}

	return url.Parse(fmt.Sprintf("/api/kubernetes/v1/cluster.json/%s/rule/%s", k.Cluster, ruleIdentifier))
}

func (k *kubeconfig) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationDestroy {
		req.Method = "POST"
		req.URL.Path = path.Dir(req.URL.Path)
	}

	req.Body = nil
	req.ContentLength = 0

	return req, nil
}

func (k *kubeconfig) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader([]byte("{}")))
	return res, nil
}
