//go:build integration
// +build integration

package v1

import (
	"context"
	"errors"
	anxAPI "github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	anxClient "github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const LoadBalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = BeforeSuite(func() {
	rand.Seed(time.Now().Unix())
})

var _ = Describe("lbaas/backend", func() {
	Context("CRUD testing", Ordered, func() {
		var backendIdentifier string
		var backendName = test.TestResourceName()
		var api anxAPI.API

		BeforeEach(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Create a backend", func() {
			b := Backend{
				Name: backendName,
			}
			backendIdentifier = createBackend(api, b, false).Identifier
		})

		It("Read Backend by ID", func() {
			b := Backend{Identifier: backendIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(backendName))
		})

		It("Update Backend", func() {
			b := Backend{Identifier: backendIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(backendName))

			newName := "Updated-Name"
			b.Name = newName
			err = api.Update(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(newName))
		})

		It("Delete Backend", func() {
			b := Backend{Identifier: backendIdentifier}
			err := api.Destroy(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(errors.Is(err, anxAPI.ErrNotFound)).To(BeTrue())
		})
	})

	Context("Listing objects", Ordered, func() {
		var api anxAPI.API

		BeforeEach(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Test listing with name", func() {
			By("Creating some backends first")

			identifiers := make([]interface{}, 3)
			identifiers[0] = createBackend(api, Backend{}, true).Identifier
			identifiers[1] = createBackend(api, Backend{}, true).Identifier
			identifiers[2] = createBackend(api, Backend{}, true).Identifier

			fetchedIdentifiers := make([]string, 0, 3)

			By("Then try finding them in a list")
			var objectChannel types.ObjectChannel
			err := api.List(context.TODO(), &Backend{Name: "go-test-%"}, anxAPI.ObjectChannel(&objectChannel))
			Expect(err).NotTo(HaveOccurred())
			for receiver := range objectChannel {
				var currBackend Backend
				err := receiver(&currBackend)
				Expect(err).NotTo(HaveOccurred())
				Expect(currBackend.Identifier).NotTo(BeEmpty())
				fetchedIdentifiers = append(fetchedIdentifiers, currBackend.Identifier)
			}
			Expect(fetchedIdentifiers).To(ContainElements(identifiers...))
		})
	})
})

func createBackend(api anxAPI.API, backend Backend, cleanup bool) Backend {
	if backend.Name == "" {
		backend.Name = test.TestResourceName()
	}
	if backend.Mode == "" {
		backend.Mode = TCP
	}

	if backend.LoadBalancer.Identifier == "" {
		backend.LoadBalancer.Identifier = LoadBalancerIdentifier
	}

	err := api.Create(context.TODO(), &backend)
	Expect(err).NotTo(HaveOccurred())
	Expect(backend.Identifier).NotTo(BeEmpty())
	if cleanup {
		DeferCleanup(func() {
			err = api.Destroy(context.TODO(), &Backend{Identifier: backend.Identifier})
			Expect(err).NotTo(HaveOccurred())
		})
	}
	return backend
}
