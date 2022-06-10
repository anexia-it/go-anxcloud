package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ACL filtering", func() {
	DescribeTable("filter parameters",
		func(a ACL, expectedKey string, expectedValue string) {
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
		Entry("no filters at all", ACL{}, "", ""),
		Entry("parent_type", ACL{ParentType: "backend"}, "parent_type", "backend"),
		Entry("frontend", ACL{Frontend: Frontend{Identifier: "some frontend identifier"}}, "frontend", "some frontend identifier"),
		Entry("backend", ACL{Backend: Backend{Identifier: "some backend identifier"}}, "backend", "some backend identifier"),
	)
})
