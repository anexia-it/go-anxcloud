package v1_test

import (
	"context"
	"errors"
	"net/http"

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

func testPrefix(test string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	c := new(api.API)
	var apiClient api.API

	BeforeEach(func() {
		a, err := e2eApiClient()
		Expect(err).ToNot(HaveOccurred())
		*c = a
		apiClient = a
	})

	Context(test, Ordered, func() {
		// netmask to the prefix to create for testing
		netmask := netmaskForFamily(fam)

		// if the prefix is expected to contain addresses
		shouldBeEmpty := true

		if fam == ipamv1.FamilyIPv4 && space == ipamv1.AddressSpacePublic {
			shouldBeEmpty = false
		}

		prefixDescription := testutils.TestResourceName()

		matchPrefixExpectation := func() types.GomegaMatcher {
			return SatisfyAll(
				WithTransform(func(p ipamv1.Prefix) string {
					return p.DescriptionCustomer
				}, Equal(prefixDescription)),

				WithTransform(func(p ipamv1.Prefix) int {
					return p.Netmask
				}, Equal(netmask)),

				WithTransform(func(p ipamv1.Prefix) ipamv1.Family {
					return p.Version
				}, Equal(fam)),

				WithTransform(func(p ipamv1.Prefix) []string {
					ret := make([]string, 0, len(p.Locations))
					for _, l := range p.Locations {
						ret = append(ret, l.Identifier)
					}
					return ret
				}, ContainElement(locationIdentifier)),

				WithTransform(func(p ipamv1.Prefix) []string {
					ret := make([]string, 0, len(p.VLANs))
					for _, v := range p.VLANs {
						ret = append(ret, v.Identifier)
					}
					return ret
				}, Equal([]string{vlanIdentifier})),
			)
		}

		var prefix ipamv1.Prefix

		It("creates a prefix", func() {
			var createEmpty *bool = nil

			if fam == ipamv1.FamilyIPv4 {
				createEmpty = &shouldBeEmpty
			}

			preparePrefixCreate(prefixDescription, createEmpty, fam, space)

			prefix = ipamv1.Prefix{
				DescriptionCustomer: prefixDescription,
				Version:             fam,
				Netmask:             netmask,
				Type:                space,
				Locations:           []corev1.Location{{Identifier: locationIdentifier}},
				VLANs:               []vlanv1.VLAN{{Identifier: vlanIdentifier}},
			}

			opts := make([]apiTypes.CreateOption, 0)

			if createEmpty != nil {
				opts = append(opts, ipamv1.CreateEmpty(*createEmpty))
			}

			err := apiClient.Create(context.TODO(), &prefix, opts...)
			Expect(err).NotTo(HaveOccurred())
		})

		It("includes the prefix when listing all prefixes", func() {
			preparePrefixList(prefixDescription, fam, space)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var oc apiTypes.ObjectChannel
			err := apiClient.List(ctx, &ipamv1.Prefix{}, api.ObjectChannel(&oc))
			Expect(err).NotTo(HaveOccurred())

			found := false
			for retriever := range oc {
				var p ipamv1.Prefix
				if err := retriever(&p); err != nil {
					cancel()
					Fail("expected retriever to succeed")
					break
				}

				if p.Identifier == prefix.Identifier {
					found = true
					cancel()
					break
				}
			}

			Expect(found).To(BeTrue())
		})

		It("eventually retrieves the prefix as Active", func() {
			Eventually(func(g types.Gomega) {
				preparePrefixEventuallyActive(prefixDescription, fam, space)

				err := apiClient.Get(context.TODO(), &prefix)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(prefix).To(matchPrefixExpectation())
				g.Expect(prefix.Status).To(Equal(ipamv1.StatusActive))
			}).Should(Succeed())
		})

		It("updates the prefix description", func() {
			prefixDescription += " - Updated!"
			preparePrefixUpdate(prefixDescription, fam, space)

			err := apiClient.Update(context.TODO(), &ipamv1.Prefix{
				Identifier:          prefix.Identifier,
				DescriptionCustomer: prefixDescription,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("retrieves the prefix with updated description", func() {
			preparePrefixGet(prefixDescription, fam, space)

			err := apiClient.Get(context.TODO(), &prefix)
			Expect(err).NotTo(HaveOccurred())

			Expect(prefix).To(matchPrefixExpectation())
			Expect(prefix.Status).To(Equal(ipamv1.StatusActive))
		})

		testAddress(c, shouldBeEmpty, &prefix)

		// we use the same code in the last spec and in AfterAll, AfterAll being there in case any previous
		// test fails
		// ref https://github.com/onsi/ginkgo/issues/933
		cleanup := func() {
			preparePrefixDelete()

			err := apiClient.Destroy(context.TODO(), &ipamv1.Prefix{Identifier: prefix.Identifier})
			Expect(err).NotTo(HaveOccurred())
		}

		AfterAll(func() {
			if prefix.Identifier != "" {
				cleanup()
			}
		})

		It("deletes the test prefix", func() {
			cleanup()
		})

		It("eventually sees the prefix being gone", func() {
			Eventually(func(g types.Gomega) {
				preparePrefixEventuallyDeleted(prefixDescription, fam, space)

				err := apiClient.Get(context.TODO(), &prefix)
				g.Expect(err).To(HaveOccurred())

				var httpError api.HTTPError
				ok := errors.As(err, &httpError)
				g.Expect(ok).To(BeTrue())

				g.Expect(httpError.StatusCode()).To(Equal(http.StatusNotFound))
			}).Should(Succeed())

			prefix.Identifier = ""
		})
	})
}
