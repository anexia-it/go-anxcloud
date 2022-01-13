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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("lbaas/server", Serial, func() {
	Context("CRUD testing", Ordered, func() {
		var api anxAPI.API
		var testingBackend Backend

		BeforeAll(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())
			testingBackend = createBackend(api, Backend{}, true)
		})

		var serverIdentifier string
		var serverName = test.TestResourceName()

		It("Create a server", func() {
			s := Server{
				Name:    serverName,
				IP:      test.RandomHostname(),
				Port:    443,
				Backend: testingBackend,
			}

			serverIdentifier = createServer(api, s, false).Identifier
		})

		It("Read Server by ID", func() {
			b := Server{Identifier: serverIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(serverName))
		})

		It("Update Server", func() {
			b := Server{Identifier: serverIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(serverName))

			newName := test.TestResourceName()
			b.Name = newName
			err = api.Update(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(newName))
		})

		It("Delete Server", func() {
			b := Server{Identifier: serverIdentifier}
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
			s := Server{
				IP:      test.RandomHostname(),
				Port:    443,
				Backend: testingBackend,
			}

			identifiers := make([]interface{}, 3)
			identifiers[0] = createServer(api, s, true).Identifier
			identifiers[1] = createServer(api, s, true).Identifier
			identifiers[2] = createServer(api, s, true).Identifier

			fetchedIdentifiers := make([]string, 0, 3)

			By("Then try finding them in a list")
			var objectChannel types.ObjectChannel
			err := api.List(context.TODO(), &Server{Name: "go-test-%"}, anxAPI.ObjectChannel(&objectChannel))
			Expect(err).NotTo(HaveOccurred())
			for receiver := range objectChannel {
				var currServer Server
				err := receiver(&currServer)
				Expect(err).NotTo(HaveOccurred())
				Expect(currServer.Identifier).NotTo(BeEmpty())
				fetchedIdentifiers = append(fetchedIdentifiers, currServer.Identifier)
			}
			Expect(fetchedIdentifiers).To(ContainElements(identifiers...))
		})
	})
})

func createServer(api anxAPI.API, s Server, cleanup bool) Server {
	if s.Name == "" {
		s.Name = test.TestResourceName()
	}

	err := api.Create(context.TODO(), &s)
	Expect(err).NotTo(HaveOccurred())
	Expect(s.Identifier).NotTo(BeEmpty())
	if cleanup {
		DeferCleanup(func() {
			err := api.Destroy(context.TODO(), &Server{Identifier: s.Identifier})
			Expect(err).NotTo(HaveOccurred())
		})
	}
	return s
}
