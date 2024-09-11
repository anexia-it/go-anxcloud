package v1_test

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address filtered listing", func() {
	DescribeTable("correctly sets query parameters",
		func(a ipamv1.Address, queryData ...string) {
			if len(queryData)%2 != 0 {
				panic("invalid test case! queryData is a list of key-value-pairs, meaning it needs a even number of entries")
			}

			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := a.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			expQuery := make(url.Values)

			for i, k := range queryData {
				if i%2 == 1 {
					continue
				}

				v := queryData[i+1]

				expQuery.Set(k, v)
			}

			Expect(u.Query()).To(BeEquivalentTo(expQuery))
		},
		Entry("no filter set", ipamv1.Address{}),

		Entry("version filter set for IPv4", ipamv1.Address{Version: ipamv1.FamilyIPv4}, "version", "4"),
		Entry("version filter set for IPv6", ipamv1.Address{Version: ipamv1.FamilyIPv6}, "version", "6"),

		Entry("type filter set for public", ipamv1.Address{Type: ipamv1.TypePublic}, "private", "false"),
		Entry("type filter set for private", ipamv1.Address{Type: ipamv1.TypePrivate}, "private", "true"),

		Entry("status filter set for StatusActive", ipamv1.Address{Status: ipamv1.StatusActive}, "status", "Active"),

		Entry("location filter set for a test Location", ipamv1.Address{Location: corev1.Location{Identifier: "foo"}}, "location", "foo"),

		Entry("prefix filter set for a test prefix", ipamv1.Address{Prefix: ipamv1.Prefix{Identifier: "foo"}}, "prefix", "foo"),

		Entry("vlan filter set for a test vlan", ipamv1.Address{VLAN: vlanv1.VLAN{Identifier: "foo"}}, "vlan", "foo"),
	)
})
