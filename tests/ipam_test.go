package tests_test

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/prefix"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IPAM API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Address endpoint", func() {

		It("Should list all available addresses", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := address.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Prefix endpoint", func() {

		It("Should list all prefixes", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := prefix.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should create a new prefix and delete it later", func() {
			p := prefix.NewAPI(cli)
			a := address.NewAPI(cli)
			ipV4 := 4
			networkMask := 29
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()

			By("Creating a new prefix")
			summary, err := p.Create(ctx, prefix.NewCreate(locationID, vlanID, ipV4, prefix.TypePrivate, networkMask))
			Expect(err).NotTo(HaveOccurred())

			var info prefix.Info
			By("Waiting for prefix to be 'Active'")
			Eventually(func() string {
				info, err = p.Get(ctx, summary.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.Vlans).NotTo(BeNil())
				Expect(info.PrefixType).To(BeEquivalentTo(prefix.TypePrivate))
				return info.Status
			}, 15*time.Minute, 5*time.Second).Should(Equal("Active"))

			Expect(info.Vlans[0].ID).To(Equal(vlanID))
			filtered, err := a.GetFiltered(ctx, 1, 1000, address.PrefixFilter(info.ID))
			By("checking that all IPs have been created in advance")
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(HaveLen(8)) // we expect all IPs to be already created
			Expect(filtered).ToNot(BeEmpty())

			By("Updating the prefix")
			_, err = p.Update(ctx, summary.ID, prefix.Update{CustomerDescription: "something else"})
			Expect(err).NotTo(HaveOccurred())

			By("Deleting the prefix")
			err = p.Delete(ctx, summary.ID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should create a new empty prefix and delete it later", func() {
			p := prefix.NewAPI(cli)
			a := address.NewAPI(cli)
			ipV4 := 4
			networkMask := 29
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()

			By("Creating a new prefix")
			create := prefix.NewCreate(locationID, vlanID, ipV4, prefix.TypePrivate, networkMask)
			create.CreateEmpty = true
			summary, err := p.Create(ctx, create)
			Expect(err).NotTo(HaveOccurred())

			var info prefix.Info
			By("Waiting for prefix to be 'Active'")
			Eventually(func() string {
				info, err = p.Get(ctx, summary.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.Vlans).NotTo(BeNil())
				Expect(info.PrefixType).To(BeEquivalentTo(prefix.TypePrivate))
				return info.Status
			}, 15*time.Minute, 5*time.Second).Should(Equal("Active"))

			Expect(info.Vlans[0].ID).To(Equal(vlanID))

			By("checking that IPs were not created in advance")
			filtered, err := a.GetFiltered(ctx, 1, 1000, address.PrefixFilter(info.ID))
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(HaveLen(3))

			By("Updating the prefix")
			_, err = p.Update(ctx, summary.ID, prefix.Update{CustomerDescription: "something else"})
			Expect(err).NotTo(HaveOccurred())

			By("Deleting the prefix")
			err = p.Delete(ctx, summary.ID)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
