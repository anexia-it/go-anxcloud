package api

import (
	"context"
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("getResponseType function", func() {
	It("returns the mime type for valid data", func() {
		rec := httptest.NewRecorder()
		rec.Header().Add("Content-Type", "application/json; charset=utf-8")
		rec.WriteHeader(500)

		ret, err := getResponseType(rec.Result())

		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal("application/json"))
	})

	It("returns an error for invalid mime type data", func() {
		rec := httptest.NewRecorder()
		rec.Header().Add("Content-Type", "foo/bar; foo")
		rec.WriteHeader(500)

		ret, err := getResponseType(rec.Result())

		Expect(err).To(MatchError(mime.ErrInvalidMediaParameter))
		Expect(ret).To(Equal(""))
	})

	It("returns an error for valid but unknown mime type", func() {
		rec := httptest.NewRecorder()
		rec.Header().Add("Content-Type", "application/pdf")
		rec.WriteHeader(500)

		ret, err := getResponseType(rec.Result())

		Expect(err).To(HaveOccurred())
		Expect(err).NotTo(MatchError(mime.ErrInvalidMediaParameter))
		Expect(err.Error()).To(ContainSubstring("application/pdf"))
		Expect(ret).To(Equal(""))
	})

	It("returns the JSON mime type when no Content-Type header is present", func() {
		rec := httptest.NewRecorder()
		rec.WriteHeader(500)

		ret, err := getResponseType(rec.Result())

		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal("application/json"))
	})
})

type api_test_anyop_option string

func (o api_test_anyop_option) ApplyToGet(opts *types.GetOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o api_test_anyop_option) ApplyToList(opts *types.ListOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o api_test_anyop_option) ApplyToCreate(opts *types.CreateOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o api_test_anyop_option) ApplyToUpdate(opts *types.UpdateOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o api_test_anyop_option) ApplyToDestroy(opts *types.DestroyOptions) {
	_ = opts.Set("api_test_option", o, false)
}

type api_test_object struct {
	Val string `json:"value" anxcloud:"identifier"`
}

var api_test_error = errors.New("we shall fail")

func (o *api_test_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	if o.Val == "failing" {
		return nil, api_test_error
	} else if o.Val == "option-check" {
		expected_option_value := ctx.Value(api_test_error)

		if v, err := opts.Get("api_test_option"); err != nil {
			return nil, err
		} else if v != expected_option_value {
			return nil, api_test_error
		}
	}

	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("Hello from api_test_object!")

	u, _ := url.ParseRequestURI("/resource/v1")
	return u, nil
}

func (o *api_test_object) FilterAPIRequest(op types.Operation, opts types.Options, req *http.Request) (*http.Request, error) {
	if o.Val == "failing_filter_request" {
		return nil, api_test_error
	}

	return req, nil
}

func (o *api_test_object) FilterAPIResponse(op types.Operation, opts types.Options, res *http.Response) (*http.Response, error) {
	if o.Val == "failing_filter_response" {
		return nil, api_test_error
	}

	return res, nil
}

func (o *api_test_object) FilterAPIRequestBody(op types.Operation, opts types.Options) (interface{}, error) {
	if o.Val == "failing_filter_request_body" {
		return nil, api_test_error
	} else if o.Val == "function_filter_request_body" {
		return func() {}, nil
	}

	return o, nil
}

type api_test_error_roundtripper bool

func (rt api_test_error_roundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, api_test_error
}

type api_test_nonstruct_object bool

func (o *api_test_nonstruct_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_nonpointer_object bool

func (o api_test_nonpointer_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_noident_object struct {
	Value string `json:"value"`
}

func (o api_test_noident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_invalidident_object struct {
	Value int `json:"value" anxcloud:"identifier"`
}

func (o api_test_invalidident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

var _ = Describe("creating API with different options", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	It("barks when creating a client without token and without ignoring the missing token", func() {
		_, err := NewAPI()
		Expect(err).To(HaveOccurred())
	})

	It("barks when making a request while using a client with unparsable BaseURL", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL("as.lfdna,smdnasd:::"), // a keysmash, added ::: to have it a really unparsable URL
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error parsing client's BaseURL"))
	})

	It("attaches the given logger to context", func() {
		server.SetAllowUnhandledRequests(true)

		log := strings.Builder{}

		logger := funcr.New(
			func(prefix, args string) {
				_, _ = log.WriteString(prefix + "\t" + args + "\n")
			},
			funcr.Options{
				Verbosity: 3,
			})

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("Hello from api_test_object!"))
	})

	It("uses a logger already on the context", func() {
		server.SetAllowUnhandledRequests(true)

		log := strings.Builder{}

		logger := funcr.New(
			func(prefix, args string) {
				_, _ = log.WriteString(prefix + "\t" + args + "\n")
			},
			funcr.Options{
				Verbosity: 3,
			})

		ctx := logr.NewContext(context.TODO(), logger)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		err = api.Create(ctx, &o)
		Expect(err).To(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("Hello from api_test_object!"))
	})

	It("handles the Object returning an error on EndpointURL", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"failing"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))

		err = api.List(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))
	})

	It("handles the Object returning an empty identifier on EndpointURL for operations requiring an IdentifiedObject", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{""}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(ErrUnidentifiedObject))
	})

	It("handles the Object returning an error on FilterAPIRequest", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"failing_filter_request"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))
	})

	It("handles the Object returning an error on FilterAPIRequestBody", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"failing_filter_request_body"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))
	})

	It("handles the Object returning a request body that cannot be encoded in json on FilterAPIRequestBody", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"function_filter_request_body"}
		err = api.Create(context.TODO(), &o)

		var e *json.UnsupportedTypeError
		Expect(errors.As(err, &e)).To(BeTrue())
	})

	It("handles the Object returning an error on FilterAPIResponse", func() {
		server.SetAllowUnhandledRequests(true)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"failing_filter_response"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))
	})

	It("handles the Engine returning a weird response content-type", func() {
		server.AppendHandlers(ghttp.RespondWith(200, nil, http.Header{"Content-Type": []string{"application/octet-stream"}}))

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(ErrUnsupportedResponseFormat))
	})

	It("handles the Engine returning bad responses for List requests", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `foo no json`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWithJSONEncoded(200, []map[string]string{{"value": "hello world"}}),
			ghttp.RespondWith(200, `foo no json`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"foo": "hello world"}),
		)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}

		var pi types.PageInfo
		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))

		var e *json.SyntaxError
		Expect(errors.As(err, &e)).To(BeTrue())

		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))
		Expect(err).NotTo(HaveOccurred())

		var os []api_test_object
		ok := pi.Next(&os)
		Expect(pi.Error()).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		ok = pi.Next(&os)
		Expect(ok).To(BeFalse())
		Expect(errors.As(pi.Error(), &e)).To(BeTrue())

		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))
		Expect(err).To(MatchError(ErrPageResponseNotSupported))
	})

	It("handles users trying to list with page iterator and channel simultaneously", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}

		var pi types.PageInfo
		oc := make(types.ObjectChannel)
		err = api.List(context.TODO(), &o, Paged(1, 2, &pi), AsObjectChannel(&oc))
		Expect(err).To(MatchError(ErrCannotListChannelAndPaged))
	})

	It("handles users trying to list page 0", func() {
		server.AppendHandlers(ghttp.RespondWithJSONEncoded(200, []string{}))

		log := strings.Builder{}

		logger := funcr.New(
			func(prefix, args string) {
				_, _ = log.WriteString(prefix + "\t" + args + "\n")
			},
			funcr.Options{
				Verbosity: 3,
			})

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}

		var pi types.PageInfo
		err = api.List(context.TODO(), &o, Paged(0, 2, &pi))
		Expect(err).NotTo(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("requesting page 0, fixing to page 1"))
	})

	It("handles http.Client.Do() returning an error", func() {
		hc := http.Client{
			Transport: api_test_error_roundtripper(false),
		}

		api, err := NewAPI(
			WithClientOptions(
				client.WithClient(&hc),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(api_test_error))
	})

	It("handles not being given a context", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		// the two nolint comments are for passing nil context, which is the behavior we want to test here.

		o := api_test_object{"identifier"}
		err = api.Get(nil, &o) //nolint:golint,staticcheck
		Expect(err).To(MatchError(ErrContextRequired))

		err = api.List(nil, &o) //nolint:golint,staticcheck
		Expect(err).To(MatchError(ErrContextRequired))
	})

	It("handles bogus operations", func() {
		api, err := NewAPI(
			WithClientOptions(
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"identifier"}
		req, err := api.(defaultAPI).makeRequest(context.TODO(), &o, &o, &types.ListOptions{}, types.Operation("bogus operation"))
		Expect(err).To(MatchError(ErrOperationNotSupported))
		Expect(req).To(BeNil())
	})

	It("consumes the given options for all operations", func() {
		opt := api_test_anyop_option("hello world")
		ctx := context.WithValue(context.TODO(), api_test_error, opt)

		server.AppendHandlers(
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, []map[string]string{{"value": "option-check"}}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{}),
		)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := api_test_object{"option-check"}

		err = api.Create(ctx, &o, opt)
		Expect(err).NotTo(HaveOccurred())

		err = api.Get(ctx, &o, opt)
		Expect(err).NotTo(HaveOccurred())

		err = api.List(ctx, &o, opt)
		Expect(err).NotTo(HaveOccurred())

		err = api.Update(ctx, &o, opt)
		Expect(err).NotTo(HaveOccurred())

		err = api.Destroy(ctx, &o, opt)
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("getObjectIdentifier function", func() {
	It("errors out on invalid Object types", func() {
		nso := api_test_nonstruct_object(false)
		identifier, err := getObjectIdentifier(&nso, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented as structs"))
		Expect(identifier).To(BeEmpty())

		npo := api_test_nonpointer_object(false)
		identifier, err = getObjectIdentifier(npo, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented on a pointer to struct"))
		Expect(identifier).To(BeEmpty())

		nio := api_test_noident_object{"invalid"}
		identifier, err = getObjectIdentifier(&nio, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("lacks identifier field"))
		Expect(identifier).To(BeEmpty())

		iio := api_test_invalidident_object{32}
		identifier, err = getObjectIdentifier(&iio, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("identifier field has an unsupported type"))
		Expect(identifier).To(BeEmpty())
	})

	Context("when doing an operation on a specific object", func() {
		It("errors out on valid Object type but empty identifier", func() {
			o := api_test_object{""}
			identifier, err := getObjectIdentifier(&o, true)
			Expect(err).To(MatchError(ErrUnidentifiedObject))
			Expect(identifier).To(BeEmpty())
		})

		It("returns the correct identifier", func() {
			o := api_test_object{"test"}
			identifier, err := getObjectIdentifier(&o, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(identifier).To(Equal("test"))
		})
	})
})

func TestAPIUnits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "api unit test suite")
}
