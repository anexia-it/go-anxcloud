package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rule filtering", func() {
	DescribeTable("filter parameters",
		func(r Rule, expectedKey string, expectedValue string) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := r.EndpointURL(ctx)
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
		Entry("no filters at all", Rule{}, "", ""),
		Entry("parent_type", Rule{ParentType: "backend"}, "parent_type", "backend"),
		Entry("condition", Rule{Condition: "if"}, "condition", "if"),
		Entry("type", Rule{Type: "connection"}, "type", "connection"),
		Entry("action ", Rule{Action: "accept"}, "action", "accept"),
		Entry("redirection_type", Rule{RedirectionType: "foo"}, "redirection_type", "foo"),
		Entry("redirection_code", Rule{RedirectionCode: "403"}, "redirection_code", "403"),
		Entry("rule_type", Rule{RuleType: "foo"}, "rule_type", "foo"),
		Entry("frontend", Rule{Frontend: Frontend{Identifier: "some frontend identifier"}}, "frontend", "some frontend identifier"),
		Entry("backend", Rule{Backend: Backend{Identifier: "some backend identifier"}}, "backend", "some backend identifier"),
	)
})
