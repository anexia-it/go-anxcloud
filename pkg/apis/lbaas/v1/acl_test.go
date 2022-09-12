package v1_test

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ACL filtering", func() {
	DescribeTable("filter parameters",
		func(a lbaasv1.ACL, expectedKey string, expectedValue string) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := a.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			q := u.Query()

			if expectedKey != "" {
				Expect(q).To(HaveKey("filters"))
				Expect(q["filters"]).To(HaveLen(1))

				filters, err := url.ParseQuery(q["filters"][0])
				Expect(err).NotTo(HaveOccurred())

				Expect(filters).To(HaveKey(expectedKey))
				Expect(filters.Get(expectedKey)).To(Equal(expectedValue))
			} else {
				Expect(q).To(BeEmpty())
			}
		},
		Entry("no filters at all", lbaasv1.ACL{}, "", ""),
		Entry("parent_type", lbaasv1.ACL{ParentType: "backend"}, "parent_type", "backend"),
		Entry("frontend", lbaasv1.ACL{Frontend: lbaasv1.Frontend{Identifier: "some frontend identifier"}}, "frontend", "some frontend identifier"),
		Entry("backend", lbaasv1.ACL{Backend: lbaasv1.Backend{Identifier: "some backend identifier"}}, "backend", "some backend identifier"),
	)
})
