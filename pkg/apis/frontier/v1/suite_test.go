package v1

import (
	"context"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

func TestFrontierAPIBindings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Frontier API Bindings Suite")
}

var _ = Describe("mock", func() {
	var (
		engine api.API
		srv    *ghttp.Server
	)

	BeforeEach(func() {
		var err error
		srv = ghttp.NewServer()
		engine, err = api.NewAPI(
			api.WithClientOptions(
				client.BaseURL(srv.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		srv.Close()
	})

	Context("API", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/api.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &API{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/api.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"description": "bar",
					"transfer_protocol": "http",
					"deployment_identifier": "baz"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &API{
				Name:                 "foo",
				Description:          pointer.String("bar"),
				TransferProtocol:     "http",
				DeploymentIdentifier: "baz",
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/frontier/v1/api.json/fake-api-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &API{Identifier: "fake-api-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Endpoint", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/endpoint.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Endpoint{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/endpoint.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"path": "bar",
					"api_identifier": "baz"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Endpoint{
				Name:          "foo",
				Path:          "bar",
				APIIdentifier: "baz",
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/frontier/v1/endpoint.json/fake-endpoint-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &Endpoint{Identifier: "fake-endpoint-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Action", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/action.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Action{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/action.json"),
				ghttp.VerifyJSON(`{
					"endpoint_identifier": "foo",
					"http_request_method": "get",
					"type": "mock_response",
					"meta": {
						"mock_response_body": "foo bar baz",
						"mock_response_language": "plaintext"
					}
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Action{
				EndpointIdentifier: "foo",
				HTTPRequestMethod:  "get",
				Type:               "mock_response",
				Meta: &ActionMeta{
					ActionMetaMockResponse: &ActionMetaMockResponse{
						Body:     "foo bar baz",
						Language: "plaintext",
					},
				},
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/frontier/v1/action.json/fake-action-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &Action{Identifier: "fake-action-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Deployment", func() {
		It("creates deployments via frontiers api deploy endpoint", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/frontier/v1/api.json/foo/deploy"),
				ghttp.VerifyJSON(`{
					"slug": "bar"	
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{
					"identifier":            "fake-api-id",
					"deployment_identifier": "fake-deployment-id",
				}),
			))

			deployment := Deployment{
				APIIdentifier: "foo",
				Slug:          "bar",

				// other fields are ignored
				Name:  "ignored",
				State: "ignored",
			}
			err := engine.Create(context.TODO(), &deployment)
			Expect(err).ToNot(HaveOccurred())
			Expect(deployment.Identifier).To(Equal("fake-deployment-id"))
		})

		It("retrieves deployments via the deployment endpoint", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/frontier/v1/deployment.json/fake-deployment-id"),
				ghttp.RespondWithJSONEncoded(200, map[string]any{
					"api_identifier": "fake-api-id",
					"state":          "disabled",
					"name":           "foo",
					"slug":           "bar",
				}),
			))

			deployment := Deployment{
				Identifier: "fake-deployment-id",
			}
			err := engine.Get(context.TODO(), &deployment)
			Expect(err).ToNot(HaveOccurred())

			Expect(deployment.APIIdentifier).To(Equal("fake-api-id"))
			Expect(deployment.State).To(Equal("disabled"))
			Expect(deployment.Name).To(Equal("foo"))
			Expect(deployment.Slug).To(Equal("bar"))
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/frontier/v1/deployment.json/fake-deployment-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &Deployment{Identifier: "fake-deployment-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
