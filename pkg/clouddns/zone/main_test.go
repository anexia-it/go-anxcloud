package zone

import (
	"fmt"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func init() {
	test.InitFlags()
}

func getClient() client.Client {
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

	c, err := client.New(options...)

	if err != nil {
		Fail(fmt.Sprintf("Error creating client: %v", err))
	}

	return c
}

func TestCloudDNS(t *testing.T) {
	RegisterFailHandler(Fail)

	suite := "CloudDNS unit tests"

	if test.RunAsIntegrationTest {
		suite = "CloudDNS integration tests"
	}

	RunSpecs(t, suite)

	if err := cleanupZones(getClient()); err != nil {
		t.Fatal(err)
	}
}
