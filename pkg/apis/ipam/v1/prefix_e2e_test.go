package v1_test

import (
	"context"

	"github.com/onsi/gomega/types"
	"go.anx.io/go-anxcloud/pkg/api"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTableSubtree("Prefix API", func(fam ipamv1.Family, addrSpace ipamv1.AddressType) {
	var a api.API
	var prefix ipamv1.Prefix

	BeforeEach(func(ctx context.Context) {
		a = GetTestAPIClient()
		netmask := 29
		if fam == ipamv1.FamilyIPv6 {
			netmask = 64
		}

		prefix = ipamv1.Prefix{
			DescriptionCustomer: testutils.TestResourceName(),
			Version:             fam,
			Netmask:             netmask,
			Type:                addrSpace,
			Locations:           []corev1.Location{{Identifier: locationIdentifier}},
			VLANs:               []vlanv1.VLAN{{Identifier: vlanIdentifier}},
		}

		preparePrefixCreate(a, prefix)
		preparePrefixGet(a, prefix)

		Expect(a.Create(ctx, &prefix)).To(Succeed())
		Eventually(func(g types.Gomega) {
			GinkgoLogr.Info("waiting for status to be active", "identifier", prefix.Identifier, "current_status", prefix.Status)
			g.Expect(a.Get(ctx, &prefix)).To(Succeed())
			g.Expect(prefix.Status).To(Equal(ipamv1.StatusActive))
		}).Should(Succeed())

		preparePrefixList(a, prefix)
		// Note: We do not delete the prefixes in a DeferCleanup method, as we leave that part to the SynchronizedAfterSuite.
	})

	When("creating a fresh prefix", func() {
		It("retrieves the prefix as Active", func(ctx context.Context) {
			Eventually(func(g types.Gomega, ctx context.Context) {
				g.Expect(a.Get(ctx, &prefix)).To(Succeed())
				g.Expect(prefix.Status).To(Equal(ipamv1.StatusActive))
			}).
				WithContext(ctx).
				Should(Succeed())
		})
		It("includes the prefix when listing all prefixes", func(ctx context.Context) {
			var oc apiTypes.ObjectChannel
			Expect(a.List(ctx, &ipamv1.Prefix{}, api.ObjectChannel(&oc))).To(Succeed())

			found := false
			for retriever := range oc {
				var p ipamv1.Prefix
				if err := retriever(&p); err != nil {
					Fail("expected retriever to succeed")
					break
				}

				if p.Identifier == prefix.Identifier {
					found = true
					break
				}
			}

			Expect(found).To(BeTrue())
		})
	})
	When("updating a prefix", Ordered, func() {
		var expectedNewDescription string
		BeforeEach(func(ctx context.Context) {
			expectedNewDescription = prefix.DescriptionCustomer + " -- updated"
			preparePrefixUpdate(a, prefix, expectedNewDescription)
			preparePrefixList(a, ipamv1.Prefix{
				Identifier:          prefix.Identifier,
				DescriptionCustomer: expectedNewDescription,
				Type:                addrSpace,
			})

			prefix.DescriptionCustomer = expectedNewDescription
			Expect(a.Update(ctx, &prefix)).To(Succeed())
		})

		It("retrieves the prefix with the updated description", func(ctx context.Context) {
			Eventually(func(g types.Gomega, ctx context.Context) {
				g.Expect(a.Get(ctx, &prefix)).To(Succeed())
				g.Expect(prefix.DescriptionCustomer).To(Equal(expectedNewDescription))
			}).
				WithContext(ctx).
				Should(Succeed())
		})
	})
	When("deleting a prefix", func() {
		BeforeEach(func(ctx context.Context) {
			preparePrefixDelete(a)
			Expect(a.Destroy(ctx, &prefix)).To(Succeed())
		})

		It("sees the prefix as deleted", func(ctx context.Context) {
			Eventually(func(g types.Gomega, ctx context.Context) {
				g.Expect(a.Get(ctx, &prefix)).To(MatchError(api.ErrNotFound))
			}).
				WithContext(ctx).
				Should(Succeed())
		})
	})
},
	Entry("private IPv4", ipamv1.FamilyIPv4, ipamv1.TypePrivate),
	Entry("private IPv6", ipamv1.FamilyIPv6, ipamv1.TypePrivate),
	Entry("public IPv4", ipamv1.FamilyIPv4, ipamv1.TypePublic),
	Entry("public IPv6", ipamv1.FamilyIPv6, ipamv1.TypePublic),
)
