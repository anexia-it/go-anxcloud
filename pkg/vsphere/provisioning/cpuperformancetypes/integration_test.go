//go:build integration
// +build integration

package cpuperformancetype

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("vsphere/provisioning/cpuperformancetype API client", func() {
	var api API

	BeforeEach(func() {
		cli, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	It("lists cpu performance types", func() {
		_, err := api.List(context.TODO())
		Expect(err).NotTo(HaveOccurred())
	})
})
