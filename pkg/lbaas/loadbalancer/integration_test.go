//go:build integration
// +build integration

package loadbalancer

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/loadbalancer client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	It("lists LoadBalancers including our test LoadBalancer", func() {
		found := false
		page := 1

		for !found {
			lbs, err := api.Get(context.TODO(), page, 20)
			Expect(err).To(BeNil())
			Expect(lbs).NotTo(BeEmpty())

			for _, lb := range lbs {
				if lb.Identifier == loadbalancerIdentifier {
					found = true
					break
				}
			}

			page++
		}
	})

	It("retrieves our test LoadBalancer with expected values", func() {
		lb, err := api.GetByID(context.TODO(), loadbalancerIdentifier)
		Expect(err).To(BeNil())
		Expect(lb.Name).To(Equal("go-anxcloud-test"))
		Expect(lb.Identifier).To(Equal(loadbalancerIdentifier))
	})
})
