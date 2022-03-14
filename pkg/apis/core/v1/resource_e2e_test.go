//go:build integration
// +build integration

package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resource E2E tests", func() {
	var apiClient api.API

	BeforeEach(func() {
		a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
		Expect(err).ToNot(HaveOccurred())
		apiClient = a
	})

	Context("with at least one resource existing", func() {
		ctx := context.TODO()

		JustBeforeEach(func() {
			// TODO: create a resource and take care to remove it after the test
		})

		It("should list resource using generic API client", func() {
			var pageIter types.PageInfo
			err := apiClient.List(ctx, &Resource{}, api.Paged(1, 100, &pageIter))
			Expect(err).ToNot(HaveOccurred())

			var resInfo []Resource
			Expect(pageIter.Next(&resInfo)).To(BeTrue())
			Expect(resInfo).ToNot(BeEmpty())
			Expect(resInfo[0].Identifier).ToNot(BeEmpty())
		})
	})
})
