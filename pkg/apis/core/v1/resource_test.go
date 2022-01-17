package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resource.Info", func() {
	When("doing unsupported operations on Info objects", func() {
		var apiClient api.API

		BeforeEach(func() {
			a, err := api.NewAPI(api.WithClientOptions(client.IgnoreMissingToken()))
			Expect(err).ToNot(HaveOccurred())
			apiClient = a
		})

		It("throws an error for Create operation", func() {
			err := apiClient.Create(context.TODO(), &Info{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
		It("throws an error for Update operation", func() {
			err := apiClient.Update(context.TODO(), &Info{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
		It("throws an error for Destroy operation", func() {
			err := apiClient.Destroy(context.TODO(), &Info{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
	})
})
