package zone

import (
	"fmt"
	"testing"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getClient() client.Client {
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

	c, err := client.New(options...)

	if err != nil {
		Fail(fmt.Sprintf("Error creating client: %v", err))
	}

	return c
}

func TestCloudDNS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudDNS tests")

	if err := cleanupZones(getClient()); err != nil {
		t.Fatal(err)
	}
}
