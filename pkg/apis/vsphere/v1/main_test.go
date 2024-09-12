package v1_test

import (
	"testing"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getApi() (api.API, error) {
	options := make([]client.Option, 0, 2)

	if isIntegrationTest {
		options = append(options, client.AuthFromEnv(false))
	} else {
		initMockServer()

		options = append(options,
			client.BaseURL(mock.URL()),
			client.IgnoreMissingToken(),
		)
	}

	return api.NewAPI(api.WithClientOptions(options...))
}

func TestVSphere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generic vSphere API tests")
}
