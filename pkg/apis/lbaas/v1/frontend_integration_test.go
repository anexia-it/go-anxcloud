//go:build integration
// +build integration

package v1

import (
	"context"
	"errors"

	anxAPI "go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	anxClient "go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("lbaas/frontend", Serial, func() {

	Context("CRUD testing", Ordered, func() {
		var api anxAPI.API
		var testingBackend Backend

		BeforeAll(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())
			testingBackend = createBackend(api, Backend{}, true)
		})

		var frontendIdentifier string
		var frontendName = test.TestResourceName()

		It("Create a frontend", func() {
			f := Frontend{
				Name:           frontendName,
				DefaultBackend: &testingBackend,
			}

			frontendIdentifier = createFrontend(api, f, false).Identifier
		})

		It("Read Frontend by ID", func() {
			b := Frontend{Identifier: frontendIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(frontendName))
		})

		It("Update Frontend", func() {
			b := Frontend{Identifier: frontendIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(frontendName))

			newName := test.TestResourceName()
			b.Name = newName
			err = api.Update(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(newName))
		})

		It("Delete Frontend", func() {
			b := Frontend{Identifier: frontendIdentifier}
			err := api.Destroy(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(errors.Is(err, anxAPI.ErrNotFound)).To(BeTrue())
		})
	})

	Context("Listing objects", Ordered, func() {
		var api anxAPI.API
		var testingBackend Backend
		BeforeAll(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())

			testingBackend = createBackend(api, Backend{}, true)
		})

		It("Test listing with name", func() {
			By("Creating some backends first")
			f := Frontend{
				Mode:           TCP,
				DefaultBackend: &testingBackend,
				LoadBalancer:   &LoadBalancer{Identifier: LoadBalancerIdentifier},
			}

			identifiers := make([]interface{}, 3)
			identifiers[0] = createFrontend(api, f, true).Identifier
			identifiers[1] = createFrontend(api, f, true).Identifier
			identifiers[2] = createFrontend(api, f, true).Identifier

			fetchedIdentifiers := make([]string, 0, 3)

			By("Then try finding them in a list")
			var objectChannel types.ObjectChannel
			err := api.List(context.TODO(), &Frontend{Name: "go-test-%"}, anxAPI.ObjectChannel(&objectChannel))
			Expect(err).NotTo(HaveOccurred())
			for receiver := range objectChannel {
				var currFrontend Frontend
				err := receiver(&currFrontend)
				Expect(err).NotTo(HaveOccurred())
				Expect(currFrontend.Identifier).NotTo(BeEmpty())
				fetchedIdentifiers = append(fetchedIdentifiers, currFrontend.Identifier)
			}
			Expect(fetchedIdentifiers).To(ContainElements(identifiers...))
		})
	})
})

func createFrontend(api anxAPI.API, frontend Frontend, cleanup bool) Frontend {
	if frontend.Name == "" {
		frontend.Name = test.TestResourceName()
	}
	if frontend.Mode == "" {
		frontend.Mode = TCP
	}
	if frontend.DefaultBackend == nil || frontend.DefaultBackend.Identifier == "" {
		backend := createBackend(api, Backend{}, true)
		frontend.DefaultBackend = &backend
	}
	if frontend.LoadBalancer == nil || frontend.LoadBalancer.Identifier == "" {
		frontend.LoadBalancer = &LoadBalancer{Identifier: LoadBalancerIdentifier}
	}

	err := api.Create(context.TODO(), &frontend)
	Expect(err).NotTo(HaveOccurred())
	Expect(frontend.Identifier).NotTo(BeEmpty())

	if cleanup {
		DeferCleanup(func() {
			err := api.Destroy(context.TODO(), &Frontend{Identifier: frontend.Identifier})
			Expect(err).NotTo(HaveOccurred())
		})
	}
	return frontend
}
