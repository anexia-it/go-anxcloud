//go:build integration
// +build integration

package location

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("vsphere/provisioning/location API client", func() {
	var api API

	BeforeEach(func() {
		cli, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	It("lists vsphere locations including ANX04", func() {
		found := false
		page := 1

		for !found {
			ls, err := api.List(context.TODO(), page, 20, "", "")
			Expect(err).To(BeNil())
			Expect(ls).NotTo(BeEmpty())

			for _, l := range ls {
				if l.Code == "ANX04" {
					found = true
					break
				}
			}

			page++
		}
	})
})
