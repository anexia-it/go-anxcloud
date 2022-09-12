package v1_test

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rule filtering", func() {
	DescribeTable("filter parameters",
		func(r lbaasv1.Rule, expectedKey string, expectedValue string) {
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
		Entry("no filters at all", lbaasv1.Rule{}, "", ""),
		Entry("parent_type", lbaasv1.Rule{ParentType: "backend"}, "parent_type", "backend"),
		Entry("condition", lbaasv1.Rule{Condition: "if"}, "condition", "if"),
		Entry("type", lbaasv1.Rule{Type: "connection"}, "type", "connection"),
		Entry("action ", lbaasv1.Rule{Action: "accept"}, "action", "accept"),
		Entry("redirection_type", lbaasv1.Rule{RedirectionType: "foo"}, "redirection_type", "foo"),
		Entry("redirection_code", lbaasv1.Rule{RedirectionCode: "403"}, "redirection_code", "403"),
		Entry("rule_type", lbaasv1.Rule{RuleType: "foo"}, "rule_type", "foo"),
		Entry("frontend", lbaasv1.Rule{Frontend: lbaasv1.Frontend{Identifier: "some frontend identifier"}}, "frontend", "some frontend identifier"),
		Entry("backend", lbaasv1.Rule{Backend: lbaasv1.Backend{Identifier: "some backend identifier"}}, "backend", "some backend identifier"),
	)
})
