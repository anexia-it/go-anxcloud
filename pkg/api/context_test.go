package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

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
		Expect(u).To(BeNil())
	default:
		Fail(fmt.Sprintf("Unknown property to test: %v", o.Test))
	}

	return url.Parse("/")
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
			ghttp.VerifyRequest("GET", fmt.Sprintf("/%v", o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has options in context for every method call", func() {
		o := context_test_object{"Options"}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/%v", o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has URL in context for every method call except EndpointURL", func() {
		o := context_test_object{"Options"}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/%v", o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Get(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})
})
