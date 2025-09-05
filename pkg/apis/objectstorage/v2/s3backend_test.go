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

var _ = Describe("S3Backend filtering", func() {
	DescribeTable("filter parameters",
		func(s objectstoragev2.S3Backend, expectedKey string, expectedValue string) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)
			u, err := s.EndpointURL(ctx)
			Expect(err).NotTo(HaveOccurred())

			q := u.Query()

			// Attributes parameter should always be present
			Expect(q).To(HaveKey("attributes"))
			Expect(q.Get("attributes")).To(Equal("name,state,endpoint,backend_type,enabled,backend_user,backend_password,reseller,customer"))

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
		Entry("no filters at all", objectstoragev2.S3Backend{}, "", ""),
		Entry("state", objectstoragev2.S3Backend{State: &objectstoragev2.GenericAttributeState{ID: "0"}}, "state", "0"),
		Entry("endpoint", objectstoragev2.S3Backend{Endpoint: common.PartialResource{Identifier: "endpoint123"}}, "endpoint", "endpoint123"),
		Entry("backend_type", objectstoragev2.S3Backend{BackendType: &objectstoragev2.GenericAttributeSelect{Identifier: "minio"}}, "backend_type", "minio"),
		Entry("customer", objectstoragev2.S3Backend{CustomerIdentifier: "customer123"}, "customer", "customer123"),
		Entry("reseller", objectstoragev2.S3Backend{ResellerIdentifier: "reseller456"}, "reseller", "reseller456"),
	)
})

var _ = Describe("S3Backend URL generation", func() {
	It("should generate correct base URL", func() {
		s := objectstoragev2.S3Backend{}
		ctx := types.ContextWithOperation(context.TODO(), types.OperationCreate)
		u, err := s.EndpointURL(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(u.Path).To(Equal("/api/object_storage/v2/s3_backend"))
	})
})
