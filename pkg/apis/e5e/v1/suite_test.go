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
)

func TestE5EAPIBindings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E5E API Bindings Suite")
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

	Context("Application", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/e5e/v1/application.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Application{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/e5e/v1/application.json"),
				ghttp.VerifyJSON(`{
					"name": "foo"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Application{
				Name: "foo",
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/e5e/v1/application.json/fake-api-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &Application{Identifier: "fake-api-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Function", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/e5e/v1/function.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Function{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/e5e/v1/function.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"application_identifier": "bar",
					"runtime": "fake-runtime",
					"entrypoint": "fake-entrypoint",
					"storage_backend": "fake-backend",
					"storage_backend_meta": {
						"git_url": "fake-url",
						"git_branch": "fake-branch",
						"git_private_key": "fake-private-key",
						"git_username": "fake-username",
						"git_password": "fake-password",
						"s3_endpoint": "fake-endpoint",
						"s3_bucket_name": "fake-bucket-name",
						"s3_object_path": "fake-object-path",
						"s3_access_key": "fake-access-key",
						"s3_secret_key": "fake-secret-key",
						"archive_file": {
							"content": "fake-content",
							"name": "fake-name"
						}
					},
					"keep_alive": 10,
					"quota_storage": 20,
					"quota_memory": 30,
					"quota_cpu": 40,
					"quota_timeout": 50,
					"quota_concurrency": 60,
					"worker_type": "standard"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))

			err := engine.Create(context.TODO(), &Function{
				Name:                  "foo",
				ApplicationIdentifier: "bar",
				Runtime:               "fake-runtime",
				Entrypoint:            "fake-entrypoint",
				StorageBackend:        "fake-backend",
				StorageBackendMeta: &StorageBackendMeta{
					StorageBackendMetaS3: &StorageBackendMetaS3{
						Endpoint:   "fake-endpoint",
						BucketName: "fake-bucket-name",
						ObjectPath: "fake-object-path",
						AccessKey:  "fake-access-key",
						SecretKey:  "fake-secret-key",
					},
					StorageBackendMetaGit: &StorageBackendMetaGit{
						URL:        "fake-url",
						Branch:     "fake-branch",
						PrivateKey: "fake-private-key",
						Username:   "fake-username",
						Password:   "fake-password",
					},
					StorageBackendMetaArchive: &StorageBackendMetaArchive{
						Content: "fake-content",
						Name:    "fake-name",
					},
				},
				KeepAlive:        10,
				QuotaStorage:     20,
				QuotaMemory:      30,
				QuotaCPU:         40,
				QuotaTimeout:     50,
				QuotaConcurrency: 60,
				WorkerType:       WorkerTypeStandard,
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not decode Destroy responses", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/e5e/v1/function.json/fake-api-id"),
				ghttp.RespondWith(http.StatusOK, "not-json"),
			))

			err := engine.Destroy(context.TODO(), &Function{Identifier: "fake-api-id"})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
