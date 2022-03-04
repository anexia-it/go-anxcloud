package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address filtered listing", func() {
	DescribeTable("correctly sets query parameters",
		func(a Address, queryData ...string) {
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
		Entry("no filter set", Address{}),

		Entry("version filter set for IPv4", Address{Version: FamilyIPv4}, "version", "4"),
		Entry("version filter set for IPv6", Address{Version: FamilyIPv6}, "version", "6"),

		Entry("type filter set for public", Address{Type: AddressSpacePublic}, "private", "false"),
		Entry("type filter set for private", Address{Type: AddressSpacePrivate}, "private", "true"),

		Entry("status filter set for StatusActive", Address{Status: StatusActive}, "status", "Active"),

		Entry("location filter set for a test Location", Address{Location: corev1.Location{Identifier: "foo"}}, "location", "foo"),

		Entry("prefix filter set for a test prefix", Address{Prefix: Prefix{Identifier: "foo"}}, "prefix", "foo"),

		Entry("vlan filter set for a test vlan", Address{VLAN: vlanv1.VLAN{Identifier: "foo"}}, "vlan", "foo"),
	)
})
