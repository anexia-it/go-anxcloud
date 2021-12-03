package v1

import (
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	test.InitFlags()
}

func getApi() (api.API, error) {
	options := make([]client.Option, 0, 2)

	if test.RunAsIntegrationTest {
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

	suite := "Generic CloudDNS unit tests"

	if test.RunAsIntegrationTest {
		suite = "Generic CloudDNS integration tests"
	}

	RunSpecs(t, suite)

	a, _ := getApi()
	if err := cleanupZones(a); err != nil {
		t.Fatal(err)
	}
}
