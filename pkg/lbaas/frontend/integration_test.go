// +build integration
// go:build integration

package frontend

import (
	"context"

	lbaasBackend "github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/frontend client", func() {
	var cli client.Client
	var api API
	var backend lbaasBackend.Backend

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
		api = NewAPI(cli)

		backendAPI := lbaasBackend.NewAPI(cli)

		b, err := backendAPI.Create(context.TODO(), lbaasBackend.Definition{
			Name:         test.TestResourceName(),
			State:        common.NewlyCreated,
			LoadBalancer: loadbalancerIdentifier,
			Mode:         common.TCP,
		})
		Expect(err).To(BeNil())

		DeferCleanup(func() {
			err := backendAPI.DeleteByID(context.TODO(), b.Identifier)
			Expect(err).To(BeNil())
		})

		backend = b
	})

	Context("with a frontend created for testing", func() {
		var definition Definition
		var frontend Frontend

		BeforeEach(func() {
			definition = Definition{
				Name:           test.TestResourceName(),
				LoadBalancer:   loadbalancerIdentifier,
				DefaultBackend: backend.Identifier,
				Mode:           common.TCP,
				State:          common.NewlyCreated,
			}

			f, err := api.Create(context.TODO(), definition)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := api.DeleteByID(context.TODO(), f.Identifier)
				Expect(err).NotTo(HaveOccurred())
			})

			Expect(f.Name).To(Equal(definition.Name))
			Expect(f.Mode).To(Equal(definition.Mode))
			Expect(f.LoadBalancer).NotTo(BeNil())
			Expect(f.LoadBalancer.Identifier).To(Equal(loadbalancerIdentifier))

			frontend = f
		})

		It("lists frontends including our test frontend", func() {
			found := false
			page := 1

			for !found {
				fs, err := api.Get(context.TODO(), page, 20)
				Expect(err).To(BeNil())
				Expect(fs).NotTo(BeEmpty())

				for _, f := range fs {
					if f.Identifier == frontend.Identifier {
						found = true
						break
					}
				}

				page++
			}
		})

		It("retrieves test frontend with expected values", func() {
			f, err := api.GetByID(context.TODO(), frontend.Identifier)

			Expect(err).To(BeNil())
			Expect(f).To(Equal(frontend))
		})

		It("updates test frontend with changed values", func() {
			definition := Definition{
				Name:           test.TestResourceName(),
				LoadBalancer:   loadbalancerIdentifier,
				DefaultBackend: backend.Identifier,
				Mode:           common.TCP,
				State:          common.Updating,
			}

			f, err := api.Update(context.TODO(), frontend.Identifier, definition)
			Expect(err).To(BeNil())

			Expect(f.Identifier).To(Equal(frontend.Identifier))
			Expect(f.Name).To(Equal(definition.Name))
			Expect(f.Mode).To(Equal(common.TCP))
		})
	})
})
