package v1

import (
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getApi() (api.API, error) {
	options := make([]client.Option, 0, 2)

	if !useMock {
		options = append(options, client.AuthFromEnv(false))
	} else {
		initMockServer()

		options = append(options,
			client.BaseURL(mock.server.URL()),
			client.IgnoreMissingToken(),
		)
	}

	return api.NewAPI(api.WithClientOptions(options...))
}

func TestCloudDNS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudDNS tests")

	a, _ := getApi()
	if err := cleanupZones(a); err != nil {
		t.Fatal(err)
	}
}
