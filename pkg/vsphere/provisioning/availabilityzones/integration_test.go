//go:build integration
// +build integration

package availabilityzones

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("vsphere/provisioning/availabilityzones API client suite", func() {
	var api API

	BeforeEach(func() {
		cli, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	It("lists availabilityzones in ANX04", func() {
		anx04 := "b164595577114876af7662092da89f76"
		_, err := api.List(context.TODO(), anx04)
		Expect(err).NotTo(HaveOccurred())
	})
})
