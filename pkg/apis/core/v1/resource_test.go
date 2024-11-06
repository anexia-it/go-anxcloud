package v1_test

import (
	"context"
	"path"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	gomegaTypes "github.com/onsi/gomega/types"
)

var _ = Describe("resource.Resource", func() {
	Context("ResourceWithTags", func() {
		rwt := &corev1.ResourceWithTag{ResourceIdentifier: "test-identifier", Tag: "test-tag"}

		DescribeTable("Test EndpointURL and FilterRequestURL for all operations", func(op types.Operation, errorMatcher gomegaTypes.GomegaMatcher, expectedPath string) {
			singleObjectOperation := op == types.OperationGet || op == types.OperationUpdate || op == types.OperationDestroy
			ctxWithOperation := types.ContextWithOperation(
				context.TODO(),
				op,
			)

			url, err := rwt.EndpointURL(ctxWithOperation)
			Expect(err).To(errorMatcher)

			if err == nil {
				Expect(url.Path).To(BeEquivalentTo(expectedPath))
				// API client appends objects identifier to path on singleObjectOperation which should be removed by FilterRequestURLHook
				if singleObjectOperation {
					url.Path = path.Join(url.Path, rwt.Identifier)
				}
				filteredURL, err := rwt.FilterRequestURL(ctxWithOperation, url)
				Expect(err).To(errorMatcher)
				Expect(filteredURL.Path).To(BeEquivalentTo(expectedPath))
			}
		},
			Entry("When operation is Create", types.OperationCreate, BeNil(), "/api/core/v1/resource.json/test-identifier/tags/test-tag"),
			Entry("When operation is Destroy", types.OperationDestroy, BeNil(), "/api/core/v1/resource.json/test-identifier/tags/test-tag"),
			Entry("When operation is Get", types.OperationGet, MatchError(api.ErrOperationNotSupported), ""),
			Entry("When operation is List", types.OperationList, MatchError(api.ErrOperationNotSupported), ""),
			Entry("When operation is Update", types.OperationUpdate, MatchError(api.ErrOperationNotSupported), ""),
		)
	})

	Context("RetryResourceTagging", func() {
		var a api.API
		var srv *ghttp.Server
		BeforeEach(func() {
			srv = ghttp.NewServer()
			var err error
			a, err = api.NewAPI(api.WithClientOptions(
				client.BaseURL(srv.URL()),
				client.IgnoreMissingToken(),
			))
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			srv.Close()
		})
		It("Can retry tagging when first attempt fails", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Can retry tagging one tag when first 2 attempts fail", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Can retry tagging one failed tag", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(200, ""),
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag", "test-tagfail")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Can retry tagging multiple failed once tags in a row", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(200, ""),
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(200, ""),
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag-fail", "test-tag-fail", "test-tag-nofail")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Fails tagging one tag when 3 attempts fail", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(409, ""),
				ghttp.RespondWith(409, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag")
			Expect(err).To(HaveOccurred())
		})
		It("Can still tag when receiving 422 for multiple tags", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(422, ""),
				ghttp.RespondWith(200, ""),
				ghttp.RespondWith(422, ""),
				ghttp.RespondWith(422, ""),
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag", "test-tag2", "test-tag3", "test-tag4", "test-tag5")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Does not retry when response code is 422 for one tag", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(422, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag")
			Expect(err).ToNot(HaveOccurred())
		})
		It("Does not retry when it should not on tagging", func() {
			srv.AppendHandlers(
				ghttp.RespondWith(200, ""),
			)
			var err = corev1.Tag(context.TODO(), a, &corev1.Resource{Identifier: "Test"}, "test-tag")
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
			err := apiClient.Create(context.TODO(), &corev1.Resource{Identifier: "foo"})
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
		})
		It("throws an error for Update operation", func() {
			err := apiClient.Update(context.TODO(), &corev1.Resource{Identifier: "foo"})
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
		})
		It("throws an error for Destroy operation", func() {
			err := apiClient.Destroy(context.TODO(), &corev1.Resource{Identifier: "foo"})
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
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

		r := corev1.Resource{}
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
