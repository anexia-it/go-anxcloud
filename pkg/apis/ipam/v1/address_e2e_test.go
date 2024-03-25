package v1_test

import (
	"context"
	"net"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testAddress(c *api.API, shouldBeEmpty bool, p *ipamv1.Prefix) {
	var apiClient api.API
	var prefix ipamv1.Prefix

	BeforeEach(func() {
		apiClient = *c
		prefix = *p
	})

	testForIP := func(ip *net.IP) {
		var address ipamv1.Address

		BeforeAll(func() {
			address = ipamv1.Address{
				Name:                ip.String(),
				DescriptionCustomer: testutils.TestResourceName(),
				Version:             prefix.Version,
				Prefix:              prefix,
			}

			if len(prefix.VLANs) == 1 {
				address.VLAN = prefix.VLANs[0]
			}
		})

		if shouldBeEmpty {
			It("creates the test address", func() {
				prepareAddressCreate(prefix, address.DescriptionCustomer, *ip)

				err := apiClient.Create(context.TODO(), &address)
				Expect(err).NotTo(HaveOccurred())
			})
		}

		It("finds the test address", func() {
			prepareAddressList(prefix, shouldBeEmpty, address.DescriptionCustomer, *ip)

			var oc types.ObjectChannel
			err := apiClient.List(
				context.TODO(),
				&ipamv1.Address{
					Prefix:  prefix,
					Version: prefix.Version,
					Type:    prefix.Type,
				},
				api.ObjectChannel(&oc),
			)
			Expect(err).NotTo(HaveOccurred())

			addressCount := 0
			for retriever := range oc {
				var addr ipamv1.Address
				err := retriever(&addr)
				Expect(err).NotTo(HaveOccurred())

				addressCount++

				if net.ParseIP(addr.Name).Equal(*ip) {
					address.Identifier = addr.Identifier
				}
			}

			// network address, gateway and our test IP
			expectedIPs := 3

			if prefix.Version == ipamv1.FamilyIPv4 {
				// for IPv4 we additionally get the broadcast address
				expectedIPs++

				if !shouldBeEmpty {
					// we test with /29 prefixes, so there should be 8 addresses when not created empty
					expectedIPs = 8
				}
			}

			Expect(addressCount).To(Equal(expectedIPs))
		})

		It("retrieves the test address", func() {
			prepareAddressGet(prefix, address.DescriptionCustomer, *ip)

			err := apiClient.Get(context.TODO(), &address)
			Expect(err).NotTo(HaveOccurred())

			Expect(address.Name).To(Equal(ip.String()))
		})

		It("updates the test address description", func() {
			address.DescriptionCustomer += " - Updated!"
			prepareAddressUpdate(prefix, address.DescriptionCustomer, *ip)

			err := apiClient.Update(context.TODO(), &ipamv1.Address{
				Identifier:          address.Identifier,
				DescriptionCustomer: address.DescriptionCustomer,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("retrieves the test address with the new description", func() {
			prepareAddressGet(prefix, address.DescriptionCustomer, *ip)

			err := apiClient.Get(context.TODO(), &address)
			Expect(err).NotTo(HaveOccurred())

			Expect(address.Name).To(Equal(ip.String()))
		})

		It("deletes the test address", func() {
			prepareAddressDelete(prefix, address.DescriptionCustomer, *ip)

			err := apiClient.Destroy(context.TODO(), &address)
			Expect(err).NotTo(HaveOccurred())
		})
	}

	Context("fixed address", Ordered, func() {
		ip := new(net.IP)

		BeforeAll(func() {
			i, _, err := net.ParseCIDR(prefix.Name)
			Expect(err).NotTo(HaveOccurred(), "expected parsable prefix")

			i[len(i)-1] += 3
			*ip = i
		})

		testForIP(ip)
	})
}
