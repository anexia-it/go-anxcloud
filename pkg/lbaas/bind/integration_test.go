// +build integration
// go:build integration

package bind

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"

	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/frontend"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/bind client", func() {
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

	Context("with a Bind created for testing", func() {
		var definition Definition
		var bind Bind

		BeforeEach(func() {
			definition = Definition{
				Name:     test.TestResourceName(),
				Frontend: frontendIdentifier,
				State:    common.NewlyCreated,
			}

			b, err := api.Create(context.TODO(), definition)
			Expect(err).NotTo(HaveOccurred())

			Expect(b.Name).To(Equal(definition.Name))
			Expect(b.Frontend.Identifier).To(Equal(frontendIdentifier))

			DeferCleanup(func() {
				err := api.DeleteByID(context.TODO(), b.Identifier)
				Expect(err).NotTo(HaveOccurred())
			})

			bind = b
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
})
