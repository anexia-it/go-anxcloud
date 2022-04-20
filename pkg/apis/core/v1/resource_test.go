package v1

import (
	"context"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resource.Resource", func() {
	Context("ResourceWithTags", func() {
		var apiClient api.API
		var rwt *ResourceWithTag

		BeforeEach(func() {
			if !isIntegrationTest {
				Skip("integration build-flag not set")
			}

			a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
			Expect(err).ToNot(HaveOccurred())
			apiClient = a
			rwt = &ResourceWithTag{Identifier: "94a4e6561ba944dfb9f5d2dfd7f10d78", Tag: "abc"}
		})

		It("tags a resource", func() {
			err := apiClient.Create(context.TODO(), rwt)
			Expect(err).ToNot(HaveOccurred())
		})

		It("untags a resource", func() {
			err := apiClient.Destroy(context.TODO(), rwt)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("doing unsupported operations on Info objects", func() {
		var apiClient api.API

		BeforeEach(func() {
			a, err := api.NewAPI(api.WithClientOptions(client.IgnoreMissingToken()))
			Expect(err).ToNot(HaveOccurred())
			apiClient = a
		})

		It("throws an error for Create operation", func() {
			err := apiClient.Create(context.TODO(), &Resource{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
		It("throws an error for Update operation", func() {
			err := apiClient.Update(context.TODO(), &Resource{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
		It("throws an error for Destroy operation", func() {
			err := apiClient.Destroy(context.TODO(), &Resource{Identifier: "foo"})
			Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
		})
	})

	It("decodes correctly", func() {
		msg := `{
	"name":"test",
	"identifier":"some identifier string",
	"resource_type":{
		"identifier":"some other identifier string",
		"name":"Service Resource"
	},
	"service_name":"Service",
	"deleted_at":null,
	"updated_at":"2022-03-23 12:19:00",
	"created_at":"2022-03-23 12:18:59",
	"reseller":{
		"customer_id":"421337",
		"demo":false,
		"identifier":"yet another identifier",
		"name":"Some reseller name",
		"name_slug":"some_reseller_name",
		"reseller":null
	},
	"customer":{
		"customer_id":"133742",
		"demo":false,
		"identifier":"even yet another identifier",
		"name":"Some customer name",
		"name_slug":"some_customer_name",
		"reseller":"yet another identifier"
	},
	"billing_contract":null,
	"managed_status":"unmanaged",
	"shared_by":null,
	"shared_at":null,
	"resource_pools":[],
	"attributes":null,
	"tags":[
		{
			"name":"some-tag",
			"identifier":"we sure have a lot of identifiers"
		},
		{
			"name":"some-other-tag",
			"identifier":"we sure have a lot of identifiers ..."
		}
	]
}`

		r := Resource{}
		err := r.DecodeAPIResponse(
			types.ContextWithOperation(context.TODO(), types.OperationGet),
			strings.NewReader(msg),
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(r.Name).To(Equal("test"))
		Expect(r.Identifier).To(Equal("some identifier string"))
		Expect(r.Type.Name).To(Equal("Service Resource"))
		Expect(r.Type.Identifier).To(Equal("some other identifier string"))
		Expect(r.Tags).To(Equal([]string{"some-tag", "some-other-tag"}))

	})
})
