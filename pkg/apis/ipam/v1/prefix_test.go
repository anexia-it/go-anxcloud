package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prefix filtered listing", func() {
	DescribeTable("correctly sets query parameters",
		func(p Prefix, queryData ...string) {
			if len(queryData)%2 != 0 {
				panic("invalid test case! queryData is a list of key-value-pairs, meaning it needs a even number of entries")
			}

			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := p.EndpointURL(ctx)
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
		Entry("no filter set", Prefix{}),

		Entry("version filter set for IPv4", Prefix{Version: FamilyIPv4}, "version", "4"),
		Entry("version filter set for IPv6", Prefix{Version: FamilyIPv6}, "version", "6"),

		Entry("type filter set for public", Prefix{Type: TypePublic}, "type", "0"),
		Entry("type filter set for private", Prefix{Type: TypePrivate}, "type", "1"),
		Entry("role_text filter set", Prefix{RoleText: "test_role"}, "role_text", "test_role"),

		Entry("status filter set for StatusActive", Prefix{Status: StatusActive}, "status", "Active"),

		Entry("location filter set for a single Location", Prefix{Locations: []corev1.Location{{Identifier: "foo"}}}, "location", "foo"),
	)
})
