package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vspherev1 "go.anx.io/go-anxcloud/pkg/apis/vsphere/v1"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("Template API bindings", func() {
	Context("BuildNumber method", func() {
		It("can parse valid Build identifier", func() {
			t := vspherev1.Template{Build: "b04"}
			Expect(t.BuildNumber()).To(Equal(4))
		})

		DescribeTable("invalid build identifier", func(build string) {
			t := vspherev1.Template{Build: build}
			_, err := t.BuildNumber()
			Expect(err).To(MatchError(vspherev1.ErrFailedToParseTemplateBuildNumber))
		},
			Entry("with empty build string", ""),
			Entry("without build digits", "b"),
			Entry("with unknown build prefix", "c123"),
			Entry("with characters between build digits", "b1c23"),
		)
	})

	Context("with mocked server", func() {
		var a api.API

		location := corev1.Location{Identifier: "mock-location-id"}

		locationWithIncompatibleJSONSchema := corev1.Location{Identifier: "location-with-incompatible-json-schema"}
		locationWithBrokenJSONResponse := corev1.Location{Identifier: "location-with-broken-json-response"}
		locationWithPlaintextResponse := corev1.Location{Identifier: "location-with-plaintext-response"}

		BeforeEach(func() {
			srv := ghttp.NewServer()

			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/templates", location.Identifier), "limit=1000&page=1"),
				ghttp.RespondWithJSONEncoded(200, mockedTemplateList()),
			))

			srv.RouteToHandler("GET", fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/templates", locationWithIncompatibleJSONSchema.Identifier),
				ghttp.RespondWithJSONEncoded(200, []map[string]interface{}{{"id": 123}}))

			srv.RouteToHandler("GET", fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/templates", locationWithBrokenJSONResponse.Identifier),
				ghttp.RespondWith(200, "not json", http.Header{"content-type": []string{"application/json"}}))

			srv.RouteToHandler("GET", fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/templates", locationWithPlaintextResponse.Identifier),
				ghttp.RespondWith(200, "not json", http.Header{"content-type": []string{"text/plain"}}),
			)

			var err error
			a, err = api.NewAPI(api.WithClientOptions(
				client.BaseURL(srv.URL()),
				client.IgnoreMissingToken(),
			))
			Expect(err).ToNot(HaveOccurred())
		})

		It("can List templates", func() {
			var channel types.ObjectChannel
			err := a.List(context.TODO(), &vspherev1.Template{Type: vspherev1.TypeTemplate, Location: location}, api.ObjectChannel(&channel))
			Expect(err).ToNot(HaveOccurred())

			templateCount := 0
			for res := range channel {
				t := vspherev1.Template{}
				err := res(&t)
				Expect(err).NotTo(HaveOccurred())
				templateCount++
			}
			Expect(templateCount).To(Equal(10))
		})

		It("can Get templates", func() {
			tpl := vspherev1.Template{Identifier: "26a47eee-dc9a-4eea-b67a-8fb1baa2fcc0", Type: vspherev1.TypeTemplate, Location: location}
			err := a.Get(context.TODO(), &tpl)
			Expect(err).ToNot(HaveOccurred())
			Expect(tpl.Name).To(Equal("Flatcar Linux Stable"))
		})

		It("returns api.ErrNotFound when template with id was not found", func() {
			err := a.Get(context.TODO(), &vspherev1.Template{Identifier: "this-id-does-not-exist", Location: location, Type: vspherev1.TypeTemplate})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("covers non-json responses", func() {
			err := a.Get(context.TODO(), &vspherev1.Template{Identifier: "random-identifier", Location: locationWithPlaintextResponse, Type: vspherev1.TypeTemplate})
			Expect(err).To(MatchError(api.ErrUnsupportedResponseFormat))
		})

		It("handles broken json in response", func() {
			err := a.Get(context.TODO(), &vspherev1.Template{Identifier: "random-identifier", Location: locationWithBrokenJSONResponse, Type: vspherev1.TypeTemplate})
			var e *json.SyntaxError
			Expect(errors.As(err, &e)).To(BeTrue())
		})

		Context("FindNamedTemplate helper", func() {
			It("handles non-json responses", func() {
				_, err := vspherev1.FindNamedTemplate(context.TODO(), a, "Debian 11", vspherev1.LatestTemplateBuild, locationWithPlaintextResponse)
				Expect(err).To(MatchError(api.ErrUnsupportedResponseFormat))
			})

			It("handles wrong json schemas in response", func() {
				_, err := vspherev1.FindNamedTemplate(context.TODO(), a, "Debian 11", vspherev1.LatestTemplateBuild, locationWithIncompatibleJSONSchema)
				var e *json.UnmarshalTypeError
				Expect(errors.As(err, &e)).To(BeTrue())
			})
		})

		DescribeTable("Template->EndpointURL returns ErrOperationNotSupported when it was called with an unsupported operation", func(op types.Operation) {
			tpl := &vspherev1.Template{Identifier: "this-id-does-not-exist", Location: location, Type: vspherev1.TypeTemplate}
			_, err := tpl.EndpointURL(types.ContextWithOperation(context.TODO(), op))
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
		},
			Entry("with operation Create", types.OperationCreate),
			Entry("with operation Destroy", types.OperationDestroy),
			Entry("with operation Update", types.OperationUpdate),
		)

		DescribeTable("find named template", func(name, build, expectedID string) {
			template, err := vspherev1.FindNamedTemplate(context.TODO(), a, name, build, location)
			Expect(err).ToNot(HaveOccurred())
			Expect(template.Identifier).To(Equal(expectedID))
		},
			Entry("latest with empty build", "Debian 11", "", "ec547552-d453-42e6-987d-51abe703c439"),
			Entry("latest with latest build", "Flatcar Linux Stable", vspherev1.LatestTemplateBuild, "26a47eee-dc9a-4eea-b67a-8fb1baa2fcc0"),
			Entry("latest with specified build", "Flatcar Linux Stable", "b74", "26a47eee-dc9a-4eea-b67a-8fb1baa2fcc0"),
			Entry("not latest build", "Windows 2022", "b06", "cb16dc94-ec55-4e9a-a1a3-b76a91bbe274"),
			Entry("with non-standard build id", "Debian 11", "possibly-valid-build-id", "9d863fd9-d0d3-4959-b226-e73192f3e43d"),
		)

		DescribeTable("find named template errors", func(name, build string) {
			_, err := vspherev1.FindNamedTemplate(context.TODO(), a, name, build, location)
			Expect(err).To(MatchError(vspherev1.ErrTemplateNotFound))
		},
			Entry("non-existing template name with build id", "FooOS 22.05", "b01"),
			Entry("non-existing template name without build id", "FooOS 22.05", ""),
			Entry("existing template name with non-existing build id", "Debian 11", "non-existing-build-id"),
		)
	})
})

func mockedTemplateList() []vspherev1.Template {
	return []vspherev1.Template{
		{Identifier: "e9325be9-25b9-468e-851e-56b5c0367e5a", Name: "Ubuntu 21.04", Build: "b72"},
		{Identifier: "b21b8b77-30e3-478a-9b6d-1f61d29e9f9a", Name: "Flatcar Linux Stable", Build: "b73"},
		{Identifier: "ec547552-d453-42e6-987d-51abe703c439", Name: "Debian 11", Build: "b18"},
		{Identifier: "26a47eee-dc9a-4eea-b67a-8fb1baa2fcc0", Name: "Flatcar Linux Stable", Build: "b74"},
		{Identifier: "cb16dc94-ec55-4e9a-a1a3-b76a91bbe274", Name: "Windows 2022", Build: "b06"},
		{Identifier: "fc3a63c6-6f4e-4193-b368-ebe9e08b4302", Name: "Debian 10", Build: "b80"},
		{Identifier: "844ac596-5f62-4ed2-936e-b99ffe0d4f88", Name: "Flatcar Linux Stable", Build: "b72"},
		{Identifier: "c3d4f0a6-978a-49fb-a952-7361bf531e4f", Name: "Debian 9", Build: "b92"},
		{Identifier: "086c5f99-1be6-46ec-8374-cdc23cedd6a4", Name: "Windows 2022", Build: "b12"},
		{Identifier: "9d863fd9-d0d3-4959-b226-e73192f3e43d", Name: "Debian 11", Build: "possibly-valid-build-id"},
	}
}
