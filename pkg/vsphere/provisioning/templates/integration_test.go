// +build integration
// go:build integration

package templates

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
)

var _ = Describe("vsphere/provisioning/templates API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("lists all available templates", func() {
		_, err := NewAPI(cli).List(context.TODO(), locationID, TemplateTypeTemplates, 1, 50)
		Expect(err).NotTo(HaveOccurred())
	})
})
