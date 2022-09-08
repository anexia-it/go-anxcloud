//go:build integration
// +build integration

package v1_test

import (
	"context"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	waitTimeout  = 5 * time.Minute
	retryTimeout = 15 * time.Second
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
			err := apiClient.List(ctx, &corev1.Resource{}, api.Paged(1, 100, &pageIter))
			Expect(err).ToNot(HaveOccurred())

			var resInfo []corev1.Resource
			Expect(pageIter.Next(&resInfo)).To(BeTrue())
			Expect(resInfo).ToNot(BeEmpty())
			Expect(resInfo[0].Identifier).ToNot(BeEmpty())
		})
	})
})

var _ = Describe("api.Create AutoTag", func() {
	var a api.API

	BeforeEach(func() {
		var err error
		a, err = api.NewAPI(
			api.WithClientOptions(client.AuthFromEnv(false)),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("todo", func() {
		It("can auto tag resources on api.Create", func() {
			vlan := vlanv1.VLAN{
				DescriptionCustomer: "go-anxcloud test api.Create AutoTag " + testutils.RandomHostname(),
				Locations: []corev1.Location{
					{Identifier: "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"},
				},
			}

			ctx := context.TODO()

			err := a.Create(ctx, &vlan, api.AutoTag("foo", "bar", "baz"))
			Expect(err).ToNot(HaveOccurred())

			tags, err := corev1.ListTags(ctx, a, &vlan)
			Expect(err).ToNot(HaveOccurred())
			Expect(tags).To(ContainElements("foo", "bar", "baz"))

			Eventually(func(g Gomega) {
				err := a.Get(ctx, &vlan)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(vlan.Status).To(Equal(vlanv1.StatusActive))
			}, waitTimeout, retryTimeout).Should(Succeed())

			err = a.Destroy(ctx, &vlan)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
