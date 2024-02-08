//go:build integration
// +build integration

package address

import (
	"context"
	"net"
	"time"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/ipam/prefix"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "166fa87362c8498f8c4aa6d1c5b9042c"
)

var _ = Describe("ipam/address client", func() {
	var api API
	var cli client.Client

	BeforeEach(func() {
		c, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
		cli = c

		api = NewAPI(cli)
	})

	var testPrefixID string
	var testPrefix net.IPNet

	createPrefix := func() {
		papi := prefix.NewAPI(cli)
		p, err := papi.Create(context.TODO(), prefix.Create{
			Location:    locationID,
			VLANID:      vlanID,
			Type:        prefix.TypePrivate,
			IPVersion:   4,
			NetworkMask: 29,

			CreateVLAN:  false,
			CreateEmpty: true,
		})
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			err := papi.Delete(context.TODO(), p.ID)
			Expect(err).NotTo(HaveOccurred())
		})

		poll := func(g Gomega) {
			p, err := papi.Get(context.TODO(), p.ID)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(p.Status).To(Equal("Active"))

			ip, pnet, err := net.ParseCIDR(p.Name)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ip.To4()).NotTo(BeNil())

			// XXX: net.IP.IsPrivate() was added in go 1.17, but we target 1.16 and have to replicate
			// that functionality here
			privateIPv4Ranges := []string{"10.0.0.0/8", "172.20.0.0/14", "192.168.0.0/16"}
			g.Expect(privateIPv4Ranges).To(WithTransform(
				func(ranges []string) []bool {
					results := make([]bool, 0, len(ranges))
					for _, r := range ranges {
						_, net, err := net.ParseCIDR(r)
						g.Expect(err).NotTo(HaveOccurred())
						results = append(results, net.Contains(ip))
					}
					return results
				},
				ContainElement(true),
			))

			mask, size := pnet.Mask.Size()
			g.Expect(mask).To(Equal(29))
			g.Expect(size).To(Equal(32))

			testPrefix = *pnet
			testPrefixID = p.ID
		}
		Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
	}

	Context("with a newly created prefix", func() {
		BeforeEach(createPrefix)

		It("lists addresses in test prefix", func() {
			ips, err := api.List(context.TODO(), 1, 20, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(ips).NotTo(BeEmpty())
		})
	})

	Context("with a prefix created for testing", Ordered, func() {
		BeforeAll(createPrefix)

		var ip net.IP
		var ipID string

		It("creates an address in test prefix", func() {
			ip = testPrefix.IP
			// +0 is network address, +1 is gateway, +2 is first free address
			ip[3] += 2

			Expect(testPrefix.Contains(ip)).To(BeTrue())

			a, err := api.Create(context.TODO(), Create{
				PrefixID:            testPrefixID,
				Address:             ip.String(),
				DescriptionCustomer: "go-anxcloud test IP",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(a.Name).To(Equal(ip.String()))
			ipID = a.ID
		})

		It("eventually retrieves the test address as inactive", func() {
			poll := func(g Gomega) {
				ip, err := api.Get(context.TODO(), ipID)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(ip.Status).To(Equal("Inactive"))
			}
			Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("updates the test address", func() {
			s, err := api.Update(context.TODO(), ipID, Update{
				DescriptionCustomer: "something something IPv4 is exhausted",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(s.ID).To(Equal(ipID))
			Expect(s.Name).To(Equal(ip.String()))
			Expect(s.DescriptionCustomer).To(Equal("something something IPv4 is exhausted"))
		})

		It("eventually retrieves the test address with changed data", func() {
			poll := func(g Gomega) {
				ip, err := api.Get(context.TODO(), ipID)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(ip.Status).To(Equal("Inactive"))
				g.Expect(ip.DescriptionCustomer).To(Equal("something something IPv4 is exhausted"))
			}
			Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("deletes the test address", func() {
			err := api.Delete(context.TODO(), ipID)
			Expect(err).NotTo(HaveOccurred())

			ipID = "" // signal the IP already being deleted, no Cleanup to do
		})

		AfterAll(func() {
			if ipID != "" {
				err := api.Delete(context.TODO(), ipID)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Context("reserving a random IP", Ordered, func() {
		var ip ReservedIP

		It("reserves a random IP", func() {
			ips, err := api.ReserveRandom(context.TODO(), ReserveRandom{
				LocationID: locationID,
				VlanID:     vlanID,
				Count:      1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(ips.TotalItems).To(Equal(1))
			Expect(ips.Data).To(HaveLen(1))
			ip = ips.Data[0]

		})

		It("eventually retrieves the reserved IP being in progress", func() {
			poll := func(g Gomega) {
				ip, err := api.Get(context.TODO(), ip.ID)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(ip.Status).To(Equal("In progress"))
			}
			Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
		})
	})
})
