//go:build integration
// +build integration

package backend

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/lbaas/common"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/backend client", Label("old client", "slow"), func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	createBackend := func(definition Definition) Backend {
		b, err := api.Create(context.TODO(), definition)
		Expect(err).NotTo(HaveOccurred())

		Expect(b.Name).To(Equal(definition.Name))
		Expect(b.Mode).To(Equal(definition.Mode))
		Expect(b.LoadBalancer).NotTo(BeNil())
		Expect(b.LoadBalancer.Identifier).To(Equal(definition.LoadBalancer))
		Expect(b.Identifier).ToNot(BeEmpty())

		DeferCleanup(func() {
			err := api.DeleteByID(context.TODO(), b.Identifier)
			Expect(err).NotTo(HaveOccurred())
		})

		return b
	}

	Context("with a backend created", func() {
		var definition Definition
		var backend Backend

		BeforeEach(func() {
			definition = Definition{
				Name:         test.TestResourceName(),
				State:        common.NewlyCreated,
				LoadBalancer: loadbalancerIdentifier,
				Mode:         common.TCP,
			}

			backend = createBackend(definition)
		})

		It("lists backends including the test backend", func() {
			found := false
			page := 1

			for !found {
				bs, err := api.Get(context.TODO(), page, 20)
				Expect(err).NotTo(HaveOccurred())
				Expect(bs).NotTo(BeEmpty())

				for _, b := range bs {
					if b.Identifier == backend.Identifier {
						found = true
						break
					}
				}

				page++
			}
		})

		It("retrieves test backend with expected values", func() {
			b, err := api.GetByID(context.TODO(), backend.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal(backend))
		})

		It("updates the test backend", func() {
			definition := Definition{
				Name:         test.TestResourceName(),
				State:        common.Updated,
				Mode:         backend.Mode,
				LoadBalancer: backend.LoadBalancer.Identifier,
			}
			b, err := api.Update(context.TODO(), backend.Identifier, definition)
			Expect(err).NotTo(HaveOccurred())

			Expect(b.Identifier).To(Equal(backend.Identifier))
			Expect(b.Name).To(Equal(definition.Name))
		})
	})

	Context("with some backends created for testing", func() {
		const numberOfTestBackends = 5

		BeforeEach(func() {
			for i := 0; i < numberOfTestBackends; i++ {
				createBackend(Definition{
					Name:         test.TestResourceName(),
					State:        common.NewlyCreated,
					LoadBalancer: loadbalancerIdentifier,
					Mode:         common.TCP,
				})
			}
		})

		It("iterates through pages as expected", func() {
			page, err := api.GetPage(context.TODO(), 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTestBackends))

			// we already had the first page
			for i := 2; i < numberOfTestBackends+1; i++ {
				page, err = api.NextPage(context.TODO(), page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})
	})
})
