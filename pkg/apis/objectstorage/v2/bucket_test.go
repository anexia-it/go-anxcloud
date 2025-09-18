package v2_test

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	objectstoragev2 "go.anx.io/go-anxcloud/pkg/apis/objectstorage/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bucket filtering", func() {
	DescribeTable("filter parameters",
		func(b objectstoragev2.Bucket, expectedKey string, expectedValue string) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := b.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			q := u.Query()

			// Attributes parameter should always be present
			Expect(q).To(HaveKey("attributes"))
			Expect(q.Get("attributes")).To(Equal("name,state,region,object_count,object_size,backend,tenant,reseller,customer"))

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
		Entry("no filters at all", objectstoragev2.Bucket{}, "", ""),
		Entry("state", objectstoragev2.Bucket{State: &objectstoragev2.GenericAttributeState{ID: "0"}}, "state", "0"),
		Entry("region", objectstoragev2.Bucket{Region: common.PartialResource{Identifier: "region123"}}, "region", "region123"),
		Entry("backend", objectstoragev2.Bucket{Backend: common.PartialResource{Identifier: "backend123"}}, "backend", "backend123"),
		Entry("tenant", objectstoragev2.Bucket{Tenant: common.PartialResource{Identifier: "tenant123"}}, "tenant", "tenant123"),
		Entry("customer", objectstoragev2.Bucket{CustomerIdentifier: "customer123"}, "customer", "customer123"),
		Entry("reseller", objectstoragev2.Bucket{ResellerIdentifier: "reseller456"}, "reseller", "reseller456"),
	)
})

var _ = Describe("Bucket URL generation", func() {
	It("should generate correct base URL", func() {
		b := objectstoragev2.Bucket{}
		ctx := types.ContextWithOperation(context.TODO(), types.OperationCreate)
		u, err := b.EndpointURL(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(u.Path).To(Equal("/api/object_storage/v2/bucket"))
	})
})

var _ = Describe("Bucket embed parameter", func() {
	DescribeTable("embed parameter generation",
		func(embeds []string, operation types.Operation, expectedEmbed string) {
			b := objectstoragev2.Bucket{Embed: embeds}
			ctx := types.ContextWithOperation(context.TODO(), operation)
			u, err := b.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			q := u.Query()
			if expectedEmbed != "" {
				Expect(q).To(HaveKey("embed"))
				Expect(q.Get("embed")).To(Equal(expectedEmbed))
			} else {
				Expect(q).NotTo(HaveKey("embed"))
			}
		},
		Entry("no embed for create operation", []string{"region", "backend"}, types.OperationCreate, ""),
		Entry("no embed for update operation", []string{"region", "backend"}, types.OperationUpdate, ""),
		Entry("single embed for list operation", []string{"region"}, types.OperationList, "region"),
		Entry("multiple embeds for list operation", []string{"region", "backend"}, types.OperationList, "region,backend"),
		Entry("single embed for get operation", []string{"backend"}, types.OperationGet, "backend"),
		Entry("multiple embeds for get operation", []string{"region", "backend"}, types.OperationGet, "region,backend"),
		Entry("empty embed array", []string{}, types.OperationList, ""),
	)
})
