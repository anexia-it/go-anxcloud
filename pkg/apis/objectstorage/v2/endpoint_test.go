package v2_test

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	objectstoragev2 "go.anx.io/go-anxcloud/pkg/apis/objectstorage/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Endpoint filtering", func() {
	DescribeTable("filter parameters",
		func(e objectstoragev2.Endpoint, expectedKey string, expectedValue string) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := e.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			q := u.Query()

			// Attributes parameter should always be present
			Expect(q).To(HaveKey("attributes"))
			Expect(q.Get("attributes")).To(Equal("url,state,endpoint_user,endpoint_password,enabled,reseller,customer"))

			if expectedKey != "" {
				Expect(q).To(HaveKey("filters"))
				Expect(q["filters"]).To(HaveLen(1))

				filters, err := url.ParseQuery(q["filters"][0])
				Expect(err).NotTo(HaveOccurred())

				Expect(filters).To(HaveKey(expectedKey))
				Expect(filters.Get(expectedKey)).To(Equal(expectedValue))
			} else {
				// Only attributes parameter should be present, no filters
				Expect(q).To(HaveLen(1))
			}
		},
		Entry("no filters at all", objectstoragev2.Endpoint{}, "", ""),
		Entry("state", objectstoragev2.Endpoint{State: &objectstoragev2.GenericAttributeState{ID: "0"}}, "state", "0"),
		Entry("endpoint_user", objectstoragev2.Endpoint{EndpointUser: "testuser"}, "endpoint_user", "testuser"),
		Entry("customer", objectstoragev2.Endpoint{CustomerIdentifier: "customer123"}, "customer", "customer123"),
		Entry("reseller", objectstoragev2.Endpoint{ResellerIdentifier: "reseller456"}, "reseller", "reseller456"),
	)
})

var _ = Describe("Endpoint URL generation", func() {
	It("should generate correct base URL", func() {
		e := objectstoragev2.Endpoint{}
		ctx := types.ContextWithOperation(context.TODO(), types.OperationCreate)
		u, err := e.EndpointURL(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(u.Path).To(Equal("/api/object_storage/v2/endpoint"))
	})
})
