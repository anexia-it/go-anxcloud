//go:build integration
// +build integration

package bind

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/lbaas/backend"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
	"go.anx.io/go-anxcloud/pkg/lbaas/frontend"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/bind client", Label("old client", "slow"), func() {
	var cli client.Client
	var api API

	var frontendIdentifier string

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)

		backendAPI := backend.NewAPI(cli)
		b, err := backendAPI.Create(context.TODO(), backend.Definition{
			Name:         test.TestResourceName(),
			LoadBalancer: loadbalancerIdentifier,
			Mode:         common.HTTP,
			State:        common.NewlyCreated,
		})
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			err := backendAPI.DeleteByID(context.TODO(), b.Identifier)
			Expect(err).NotTo(HaveOccurred())
		})

		frontendAPI := frontend.NewAPI(cli)
		f, err := frontendAPI.Create(context.TODO(), frontend.Definition{
			Name:           test.TestResourceName(),
			LoadBalancer:   loadbalancerIdentifier,
			DefaultBackend: b.Identifier,
			Mode:           common.HTTP,
			State:          common.NewlyCreated,
		})
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			err := frontendAPI.DeleteByID(context.TODO(), f.Identifier)
			Expect(err).NotTo(HaveOccurred())
		})

		frontendIdentifier = f.Identifier
	})

	createBind := func(definition Definition) Bind {
		b, err := api.Create(context.TODO(), definition)
		Expect(err).NotTo(HaveOccurred())

		Expect(b.Name).To(Equal(definition.Name))
		Expect(b.Frontend.Identifier).To(Equal(frontendIdentifier))

		DeferCleanup(func() {
			err := api.DeleteByID(context.TODO(), b.Identifier)
			Expect(err).NotTo(HaveOccurred())
		})

		return b
	}

	Context("with a Bind created for testing", func() {
		var definition Definition
		var bind Bind

		BeforeEach(func() {
			definition = Definition{
				Name:     test.TestResourceName(),
				Frontend: frontendIdentifier,
				State:    common.NewlyCreated,
			}

			bind = createBind(definition)
		})

		It("lists Binds including our test Bind", func() {
			found := false
			page := 1

			for !found {
				bs, err := api.Get(context.TODO(), page, 20)
				Expect(err).To(BeNil())
				Expect(bs).NotTo(BeEmpty())

				for _, b := range bs {
					if b.Identifier == bind.Identifier {
						found = true
						break
					}
				}

				page++
			}
		})

		It("retrieves test Bind with expected values", func() {
			b, err := api.GetByID(context.TODO(), bind.Identifier)
			Expect(err).To(BeNil())
			Expect(b).To(Equal(bind))
		})

		It("updates test bind with new values", func() {
			definition := Definition{
				Name:     test.TestResourceName(),
				State:    common.Updating,
				Frontend: frontendIdentifier,
			}

			b, err := api.Update(context.TODO(), bind.Identifier, definition)
			Expect(err).To(BeNil())

			Expect(b.Identifier).To(Equal(bind.Identifier))
			Expect(b.Name).To(Equal(definition.Name))
		})
	})

	Context("with some binds created for testing", func() {
		const numberOfTestBinds = 5

		BeforeEach(func() {
			for i := 0; i < numberOfTestBinds; i++ {
				createBind(Definition{
					Name:     test.TestResourceName(),
					Frontend: frontendIdentifier,
					State:    common.NewlyCreated,
				})
			}
		})

		It("iterates through pages as expected", func() {
			page, err := api.GetPage(context.TODO(), 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTestBinds))

			// we already had the first page
			for i := 2; i < numberOfTestBinds+1; i++ {
				page, err = api.NextPage(context.TODO(), page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})
	})
})
