package v1_test

import (
	"context"
	"net/netip"

	"github.com/onsi/gomega/types"
	"go.anx.io/go-anxcloud/pkg/api"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = DescribeTableSubtree("Address API", func(ipVersion ipamv1.Family, addrType ipamv1.AddressType) {
	var a api.API
	var prefix ipamv1.Prefix
	var addr ipamv1.Address

	BeforeEach(func(ctx context.Context) {
		a = GetTestAPIClient()
		prefix = GetTestPrefix(ipVersion, addrType)

		gatewayAddr := netip.MustParsePrefix(prefix.Name).Addr().Next()
		addr = ipamv1.Address{
			DescriptionCustomer: testutils.TestResourceName(),
			Version:             ipVersion,
			Prefix:              prefix,
			Type:                addrType,

			// Gateway addresses are not editable in the engine, therefore we choose the next one after it.
			Name:     gatewayAddr.Next().String(),
			Location: corev1.Location{Identifier: locationIdentifier},
			VLAN:     vlanv1.VLAN{Identifier: vlanIdentifier},

			// We have to set the address to reserved, so it'll get marked as "Active" by the engine.
			RoleText: "Reserved",
		}

		// Note: Since we delete the whole prefix in the SynchronizedAfterSuite, we do not clean up individually IP addresses.

		prepareAddressCreate(a, &addr)
		prepareAddressGet(a, addr)

		Expect(a.Create(ctx, &addr, ipamv1.CreateEmpty(true))).To(Succeed())
	})

	When("creating a fresh address", func() {
		It("retrieves the address as Active", func(ctx context.Context) {
			Eventually(func(g types.Gomega) {
				g.Expect(a.Get(ctx, &addr)).To(Succeed())
				g.Expect(addr.Status).To(Equal(ipamv1.StatusActive))
			}).Should(Succeed())
		})

		Context("listing", func() {
			var addresses []ipamv1.Address

			BeforeEach(func(ctx context.Context) {
				prepareAddressList(a, prefix, addr)

				addresses = []ipamv1.Address{}

				var oc apiTypes.ObjectChannel
				Expect(a.List(ctx, &ipamv1.Address{
					Version: ipVersion,
					Prefix:  addr.Prefix,
					Type:    addrType,
				}, api.ObjectChannel(&oc))).To(Succeed())

				for retriever := range oc {
					var p ipamv1.Address
					if err := retriever(&p); err != nil {
						Fail("expected retriever to succeed")
						break
					}

					addresses = append(addresses, p)
				}
			})

			It("finds the address", func() {
				Expect(addresses).To(ContainElement(MatchFields(IgnoreExtras, Fields{
					"Name": Equal(addr.Name),
				})))
			})
		})
	})

	When("updating an address", func() {
		var expectedNewDescription string
		BeforeEach(func(ctx context.Context) {
			// In order for an update to work, we have to wait until the address is in an active state.
			Eventually(func(g types.Gomega) {
				g.Expect(a.Get(ctx, &addr)).To(Succeed())
				GinkgoLogr.Info("waiting for address to be active", "addr", addr.Name, "current_status", addr.Status)
				g.Expect(addr.Status).To(Equal(ipamv1.StatusActive))
			}).Should(Succeed())

			expectedNewDescription = addr.DescriptionCustomer + " -- updated"

			updatedAddr := ipamv1.Address{
				Identifier:          addr.Identifier,
				Name:                addr.Name,
				DescriptionCustomer: expectedNewDescription,
				Version:             ipVersion,
				RoleText:            addr.RoleText,
				Status:              addr.Status,
				VLAN:                addr.VLAN,
				Prefix:              prefix,
				Location:            addr.Location,
				Type:                addrType,
			}
			prepareAddressUpdate(a, addr, expectedNewDescription)
			prepareAddressList(a, prefix, updatedAddr)
			prepareAddressGet(a, updatedAddr)

			addr.DescriptionCustomer = expectedNewDescription
			Expect(a.Update(ctx, &addr)).To(Succeed())
		})

		It("updates the description", func(ctx context.Context) {
			Eventually(func(g types.Gomega) {
				g.Expect(a.Get(ctx, &addr)).To(Succeed())
				g.Expect(addr.DescriptionCustomer).To(Equal(expectedNewDescription))
			}).Should(Succeed())
		})
		It("retrieves the address with the updated description", func(ctx context.Context) {
			var oc apiTypes.ObjectChannel
			Expect(a.List(ctx, &ipamv1.Address{
				Version: ipVersion,
				Type:    addrType,
				Prefix:  prefix,
			}, api.ObjectChannel(&oc))).To(Succeed())

			for retriever := range oc {
				var p ipamv1.Prefix
				if err := retriever(&p); err != nil {
					Fail("expected retriever to succeed")
					break
				}

				if p.Identifier == addr.Identifier {
					Expect(p.DescriptionCustomer).To(Equal(expectedNewDescription))
					return
				}
			}

			Fail("Address not found in list operation.")
		})
	})
	When("deleting an address", func() {
		BeforeEach(func(ctx context.Context) {
			prepareAddressDelete(a, addr)
			Eventually(func(g types.Gomega, ctx context.Context) {
				g.Expect(a.Get(ctx, &addr)).To(Succeed())
				g.Expect(addr.Status).To(Equal(ipamv1.StatusActive))
			}).
				WithContext(ctx).
				Should(Succeed())

			Expect(a.Destroy(ctx, &addr)).To(Succeed())
		})

		It("sees the address as deleted", func(ctx context.Context) {
			Eventually(func(g types.Gomega) {
				g.Expect(a.Get(ctx, &addr)).To(MatchError(api.ErrNotFound))
			}).Should(Succeed())
		})
	})
},
	Entry("private IPv4", ipamv1.FamilyIPv4, ipamv1.TypePrivate),
	Entry("private IPv6", ipamv1.FamilyIPv6, ipamv1.TypePrivate),
	Entry("public IPv4", ipamv1.FamilyIPv4, ipamv1.TypePublic),
	Entry("public IPv6", ipamv1.FamilyIPv6, ipamv1.TypePublic),
)
