//go:build integration
// +build integration

package vlan

import (
	"context"
	"time"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
)

var _ = Describe("vlan client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	Context("with a VLAN created for testing", Ordered, func() {
		var vlanID string

		BeforeAll(func() {
			def := CreateDefinition{
				Location:            locationID,
				VMProvisioning:      false,
				CustomerDescription: "go SDK integration test",
			}
			summary, err := api.Create(context.TODO(), def)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := api.Delete(context.TODO(), summary.Identifier)
				Expect(err).NotTo(HaveOccurred())
			})

			vlanID = summary.Identifier
		})

		It("lists VLANs including test VLAN", func() {
			found := false
			page := 1

			for !found {
				vs, err := api.List(context.TODO(), page, 20, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(vs).NotTo(BeEmpty())

				for _, v := range vs {
					if v.Identifier == vlanID {
						found = true
						break
					}
				}

				page++
			}
		})

		It("eventually retrieves test VLAN with expected data and being active", func() {
			pollCheck := func(g Gomega) {
				vlan, err := api.Get(context.TODO(), vlanID)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(vlan.Locations).To(HaveLen(1))
				g.Expect(vlan.Locations[0].Identifier).To(Equal(locationID))

				g.Expect(vlan.VMProvisioning).To(BeFalse())
				g.Expect(vlan.CustomerDescription).To(Equal("go SDK integration test"))

				g.Expect(vlan.Status).To(Equal("Active"))
			}

			Eventually(pollCheck, 5*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("updates test VLAN with changed data", func() {
			def := UpdateDefinition{
				CustomerDescription: "go SDK integration test updated",
				VMProvisioning:      true,
			}

			err := api.Update(context.TODO(), vlanID, def)
			Expect(err).NotTo(HaveOccurred())
		})

		It("eventually retrieves test VLAN with expected updated data", func() {
			pollCheck := func(g Gomega) {
				vlan, err := api.Get(context.TODO(), vlanID)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(vlan.Locations).To(HaveLen(1))
				g.Expect(vlan.Locations[0].Identifier).To(Equal(locationID))

				g.Expect(vlan.VMProvisioning).To(BeTrue())
				g.Expect(vlan.CustomerDescription).To(Equal("go SDK integration test updated"))
				g.Expect(vlan.Status).To(Equal("Active"))
			}

			Eventually(pollCheck, 5*time.Minute, 10*time.Second).Should(Succeed())
		})
	})
})
