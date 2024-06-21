//go:build !integration
// +build !integration

package v1_test

import (
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/gomega/ghttp"
)

var (
	mockServer *Server
)

func e2eApiClient() (api.API, error) {
	if mockServer == nil {
		mockServer = NewServer()
	}

	vlanIdentifier = "randomVLANIdentifier"

	return api.NewAPI(
		api.WithClientOptions(
			client.BaseURL(mockServer.URL()),
			client.IgnoreMissingToken(),
		),
	)
}
