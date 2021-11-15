package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const context_test_object_baseurl = "/v1/context_test_object"

type context_test_object struct {
	Test string `anxcloud:"identifier"`
}

func (o *context_test_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	switch o.Test {
	case "Operation":
		Expect(OperationFromContext(ctx)).To(Equal(op))
	case "Options":
		Expect(OptionsFromContext(ctx)).To(Equal(opts))
	case "URL":
		u, err := URLFromContext(ctx)
		Expect(err).To(MatchError(ErrContextKeyNotSet))
		Expect(u).To(BeZero())
	default:
		Fail(fmt.Sprintf("Unknown property to test: %v", o.Test))
	}

	return url.Parse(context_test_object_baseurl)
}

func (o *context_test_object) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	switch o.Test {
	case "Operation":
		Expect(OperationFromContext(ctx)).To(Equal(types.OperationGet))
	case "Options":
		Expect(OptionsFromContext(ctx)).NotTo(BeNil())
		Expect(OptionsFromContext(ctx)).To(BeAssignableToTypeOf(&types.GetOptions{}))
	case "URL":
		compare, err := url.Parse(context_test_object_baseurl)
		Expect(err).NotTo(HaveOccurred())

		Expect(URLFromContext(ctx)).To(Equal(*compare))
	default:
		Fail(fmt.Sprintf("Unknown property to test: %v", o.Test))
	}

	d := json.NewDecoder(data)
	d.DisallowUnknownFields()
	return d.Decode(o)
}

var _ = Describe("context passed to Object methods", func() {
	var server *ghttp.Server
	var api API
	var ctx context.Context

	JustBeforeEach(func() {
		ctx = context.TODO()

		server = ghttp.NewServer()
		a, err := NewAPI(WithClientOptions(
			client.IgnoreMissingToken(),
			client.BaseURL(server.URL()),
		))

		Expect(err).NotTo(HaveOccurred())
		api = a
	})

	It("has operation in context for every method call", func() {
		o := context_test_object{"Operation"}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("%v/%v", context_test_object_baseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has options in context for every method call", func() {
		o := context_test_object{"Options"}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("%v/%v", context_test_object_baseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has URL in context for every method call except EndpointURL", func() {
		o := context_test_object{"URL"}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("%v/%v", context_test_object_baseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("context key retriever functions", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.TODO()
	})

	Context("with no attributes set", func() {
		It("returns ErrContextKeyNotSet error for URL", func() {
			u, err := URLFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(u).To(BeZero())
		})

		It("returns ErrContextKeyNotSet error for Operation", func() {
			o, err := OperationFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(o).To(BeZero())
		})

		It("returns ErrContextKeyNotSet error for Options", func() {
			o, err := OptionsFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(o).To(BeNil())
		})
	})
})
