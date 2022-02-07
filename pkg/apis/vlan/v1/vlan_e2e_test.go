package v1

import (
	"context"
	"errors"
	"fmt"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"

	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const locationIdentifier = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"

var _ = Describe("VLAN E2E tests", func() {
	var apiClient api.API

	BeforeEach(func() {
		a, err := e2eApiClient()
		Expect(err).ToNot(HaveOccurred())
		apiClient = a

	})

	Context("with a VLAN created for testing", Ordered, func() {
		var identifier string
		var desc string

		deleted := false

		BeforeAll(func() {
			desc = "go-anxcloud test " + testutils.RandomHostname()
			prepareCreate(desc)

			vlan := VLAN{
				DescriptionCustomer: desc,
				Locations:           []corev1.Location{{Identifier: locationIdentifier}},
			}
			err := apiClient.Create(context.TODO(), &vlan)
			Expect(err).NotTo(HaveOccurred())

			identifier = vlan.Identifier

			DeferCleanup(func() {
				if !deleted {
					err := apiClient.Destroy(context.TODO(), &VLAN{Identifier: identifier})
					if err != nil {
						Fail(fmt.Sprintf("Error destroying VLAN %q created for testing: %v", identifier, err))
					}
				}
			})
		})

		It("should retrieve the VLAN", func() {
			prepareGet(desc, false)

			vlan := VLAN{Identifier: identifier}
			err := apiClient.Get(context.TODO(), &vlan)

			Expect(err).NotTo(HaveOccurred())
			Expect(vlan.DescriptionCustomer).To(Equal(desc))
		})

		It("should list resource using generic API client", func() {
			prepareList(desc, false)

			var oc types.ObjectChannel
			err := apiClient.List(context.TODO(), &VLAN{}, api.ObjectChannel(&oc))
			Expect(err).ToNot(HaveOccurred())

			found := false
			for r := range oc {
				vlan := VLAN{}
				err := r(&vlan)
				Expect(err).NotTo(HaveOccurred())

				if vlan.Identifier == identifier {
					found = true
					Expect(vlan.DescriptionCustomer).To(Equal(desc))
				}
			}

			Expect(found).To(BeTrue())
		})

		It("eventually has StatusActive", func() {
			Eventually(func(g Gomega) {
				prepareEventuallyActive(desc, false)

				vlan := VLAN{Identifier: identifier}
				err := apiClient.Get(context.TODO(), &vlan)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(vlan.Status).To(Equal(StatusActive))
			}, waitTimeout, retryTimeout).Should(Succeed())
		})

		It("updates to VM provisioning enabled", func() {
			prepareUpdate(desc, true)

			vlan := VLAN{
				Identifier:          identifier,
				DescriptionCustomer: desc,
				VMProvisioning:      true,
			}
			err := apiClient.Update(context.TODO(), &vlan)
			Expect(err).NotTo(HaveOccurred())
		})

		It("eventually shows VM provisioning as enabled", func() {
			Eventually(func(g Gomega) {
				prepareGet(desc, true)

				vlan := VLAN{Identifier: identifier}
				err := apiClient.Get(context.TODO(), &vlan)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(vlan.VMProvisioning).To(BeTrue())
			}, waitTimeout, retryTimeout).Should(Succeed())
		})

		It("should destroy the VLAN", func() {
			prepareDelete()

			err := apiClient.Destroy(context.TODO(), &VLAN{Identifier: identifier})
			Expect(err).NotTo(HaveOccurred())
		})

		It("has StatusMarkedForDeletion", func() {
			prepareDeleting()

			vlan := VLAN{Identifier: identifier}
			err := apiClient.Get(context.TODO(), &vlan)

			Expect(err).NotTo(HaveOccurred())
			Expect(vlan.Status).To(Equal(StatusMarkedForDeletion))
		})

		It("eventually is gone", func() {
			Eventually(func(g Gomega) {
				prepareEventuallyDeleted(desc, true)

				vlan := VLAN{Identifier: identifier}
				err := apiClient.Get(context.TODO(), &vlan)
				g.Expect(err).To(HaveOccurred())

				he := api.HTTPError{}
				ok := errors.As(err, &he)
				g.Expect(ok).To(BeTrue())

				g.Expect(he.StatusCode()).To(Equal(404))
			}, waitTimeout, retryTimeout).Should(Succeed())

			deleted = true
		})
	})
})
