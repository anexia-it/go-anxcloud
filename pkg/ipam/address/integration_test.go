//go:build integration
// +build integration

package address

import (
	"context"
	"errors"
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

		poll := func() error {
			p, err := papi.Get(context.TODO(), p.ID)
			if err != nil {
				return err
			}
			if p.Status != "Active" {
				return errors.New("prefix not active")
			}

			ip, pnet, err := net.ParseCIDR(p.Name)
			if err != nil {
				return err
			}
			if ip.To4() == nil {
				return errors.New("not IPv4 address")
			}

			// XXX: net.IP.IsPrivate() was added in go 1.17, but we target 1.16 and have to replicate
			// that functionality here
			privateIPv4Ranges := []string{"10.0.0.0/8", "172.20.0.0/14", "192.168.0.0/16"}
			isPrivate := false
			for _, r := range privateIPv4Ranges {
				_, net, err := net.ParseCIDR(r)
				if err != nil {
					return err
				}
				if net.Contains(ip) {
					isPrivate = true
					break
				}
			}
			if !isPrivate {
				return errors.New("IP is not in private range")
			}

			mask, size := pnet.Mask.Size()
			if mask != 29 {
				return errors.New("unexpected mask size")
			}
			if size != 32 {
				return errors.New("unexpected address size")
			}

			testPrefix = *pnet
			testPrefixID = p.ID
			return nil
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
			poll := func() error {
				ip, err := api.Get(context.TODO(), ipID)
				if err != nil {
					return err
				}
				if ip.Status != "Inactive" {
					return errors.New("IP status not inactive")
				}
				return nil
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
			poll := func() error {
				ip, err := api.Get(context.TODO(), ipID)
				if err != nil {
					return err
				}
				if ip.Status != "Inactive" {
					return errors.New("IP status not inactive")
				}
				if ip.DescriptionCustomer != "something something IPv4 is exhausted" {
					return errors.New("description not updated")
				}
				return nil
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
			poll := func() error {
				ip, err := api.Get(context.TODO(), ip.ID)
				if err != nil {
					return err
				}
				if ip.Status != "In progress" {
					return errors.New("IP status not in progress")
				}
				return nil
			}
			Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
		})
	})
})
