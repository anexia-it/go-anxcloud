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

var _ = Describe("lbaas/bind", Serial, func() {
	Context("CRUD testing", Ordered, func() {
		var api anxAPI.API
		var testingFrontend Frontend
		var bindIdentifier string
		var bindName = test.TestResourceName()

		BeforeAll(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())
			testingFrontend = createFrontend(api, Frontend{}, true)
		})

		It("Create a bind", func() {
			b := Bind{
				Name:     bindName,
				Frontend: testingFrontend,
			}

			bindIdentifier = createBind(api, b, false).Identifier
		})

		It("Read Bind by ID", func() {
			b := Bind{Identifier: bindIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(bindName))
		})

		It("Update Bind", func() {
			b := Bind{Identifier: bindIdentifier}
			err := api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(bindName))

			newName := test.TestResourceName()
			b.Name = newName
			err = api.Update(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())
			Expect(b.Name).To(BeEquivalentTo(newName))
		})

		It("Delete Bind", func() {
			b := Bind{Identifier: bindIdentifier}
			err := api.Destroy(context.TODO(), &b)
			Expect(err).NotTo(HaveOccurred())

			err = api.Get(context.TODO(), &b)
			Expect(errors.Is(err, anxAPI.ErrNotFound)).To(BeTrue())
		})
	})

	Context("Listing objects", Ordered, func() {
		var api anxAPI.API
		var testingFrontend Frontend

		BeforeAll(func() {
			var err error
			api, err = anxAPI.NewAPI(anxAPI.WithClientOptions(anxClient.TokenFromEnv(false)))
			Expect(err).NotTo(HaveOccurred())

			testingFrontend = createFrontend(api, Frontend{}, true)
		})

		It("Test listing with name", func() {
			By("Creating some backends first")
			b := Bind{
				Frontend: testingFrontend,
			}
			identifiers := make([]interface{}, 3)
			identifiers[1] = createBind(api, b, true).Identifier
			identifiers[0] = createBind(api, b, true).Identifier
			identifiers[2] = createBind(api, b, true).Identifier

			fetchedIdentifiers := make([]string, 0, 3)

			By("Then try finding them in a list")
			var objectChannel types.ObjectChannel
			err := api.List(context.TODO(), &Bind{Name: "go-test-%"}, anxAPI.ObjectChannel(&objectChannel))
			Expect(err).NotTo(HaveOccurred())
			for receiver := range objectChannel {
				var currBind Bind
				err := receiver(&currBind)
				Expect(err).NotTo(HaveOccurred())
				Expect(currBind.Identifier).NotTo(BeEmpty())
				fetchedIdentifiers = append(fetchedIdentifiers, currBind.Identifier)
			}
			Expect(fetchedIdentifiers).To(ContainElements(identifiers...))
		})
	})
})

func createBind(api anxAPI.API, bind Bind, cleanup bool) Bind {
	if bind.Name == "" {
		bind.Name = test.TestResourceName()
	}

	if bind.Frontend.Identifier == "" {
		bind.Frontend.Identifier = createFrontend(api, Frontend{}, true).Identifier
	}

	err := api.Create(context.TODO(), &bind)
	Expect(err).NotTo(HaveOccurred())
	Expect(bind.Identifier).NotTo(BeEmpty())

	if cleanup {
		DeferCleanup(func() {
			err := api.Destroy(context.TODO(), &Bind{Identifier: bind.Identifier})
			Expect(err).NotTo(HaveOccurred())
		})
	}

	return bind
}
