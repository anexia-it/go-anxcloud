package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"
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

type apiTestAnyopOption string

func (o apiTestAnyopOption) ApplyToGet(opts *types.GetOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o apiTestAnyopOption) ApplyToList(opts *types.ListOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o apiTestAnyopOption) ApplyToCreate(opts *types.CreateOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o apiTestAnyopOption) ApplyToUpdate(opts *types.UpdateOptions) {
	_ = opts.Set("api_test_option", o, false)
}

func (o apiTestAnyopOption) ApplyToDestroy(opts *types.DestroyOptions) {
	_ = opts.Set("api_test_option", o, false)
}

type apiTestObject struct {
	Val string `json:"value" anxcloud:"identifier"`
}

var errAPITest = errors.New("we shall fail")

func (o *apiTestObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	if o.Val == "failing" {
		return nil, errAPITest
	} else if o.Val == "option-check" {
		expectedOptionValue := ctx.Value(errAPITest)

		opts, err := types.OptionsFromContext(ctx)
		if err != nil {
			return nil, err
		}

		if v, err := opts.Get("api_test_option"); err != nil {
			return nil, err
		} else if v != expectedOptionValue {
			return nil, errAPITest
		}
	}

	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("Hello from apiTestObject!")

	u, _ := url.ParseRequestURI("/resource/v1")
	return u, nil
}

func (o *apiTestObject) GetIdentifier(context.Context) (string, error) {
	if o.Val == "failing-identifier" {
		return "", errAPITest
	}
	return o.Val, nil
}

func (o *apiTestObject) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	if o.Val == "failing_filter_request" {
		return nil, errAPITest
	}

	return req, nil
}

func (o *apiTestObject) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	if o.Val == "failing_filter_response" {
		return nil, errAPITest
	}

	return res, nil
}

func (o *apiTestObject) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	if o.Val == "failing_filter_request_body" {
		return nil, errAPITest
	} else if o.Val == "function_filter_request_body" {
		return func() {}, nil
	}

	return o, nil
}

func (o *apiTestObject) HasPagination(ctx context.Context) (bool, error) {
	if o.Val == "failing_has_pagination" {
		return false, errAPITest
	}

	return o.Val != "no_pagination", nil
}

func (o *apiTestObject) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	if o.Val == "failing_decode_response" {
		return errAPITest
	} else if o.Val == "success_decode_response" {
		o.Val = "Decode hook called!"
		return nil
	}

	return json.NewDecoder(data).Decode(o)
}

type apiTestErrorRoundtripper bool

func (rt apiTestErrorRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errAPITest
}

var _ = Describe("decodeResponse function", func() {
	var ctx context.Context

	JustBeforeEach(func() {
		ctx = types.ContextWithOptions(
			types.ContextWithOperation(
				types.ContextWithURL(
					context.TODO(),
					url.URL{Path: "/"},
				),
				types.OperationGet,
			),
			&types.GetOptions{},
		)
	})

	It("fails on media types other than application/json", func() {
		var out json.RawMessage
		err := decodeResponse(ctx, "foo/bar", &bytes.Buffer{}, &out)
		Expect(err).To(MatchError(ErrUnsupportedResponseFormat))
	})

	It("decodes json message into []json.RawMessage", func() {
		var out []json.RawMessage
		err := decodeResponse(ctx, "application/json", bytes.NewReader([]byte(`[{},{}]`)), &out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(HaveLen(2))
	})

	It("decodes json message using the Objects response decode hook", func() {
		obj := apiTestObject{"success_decode_response"}
		err := decodeResponse(ctx, "application/json", bytes.NewReader([]byte(`{"value": "decode hook not called :("}`)), &obj)
		Expect(err).NotTo(HaveOccurred())
		Expect(obj.Val).To(Equal("Decode hook called!"))
	})
})

var _ = Describe("using an API object", func() {
	var server *ghttp.Server

	logger := stdr.New(log.New(GinkgoWriter, "", log.Ltime|log.Lshortfile|log.Lmsgprefix))
	stdr.SetVerbosity(3)

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	It("barks when creating a client without token and without ignoring the missing token", func() {
		_, err := NewAPI()
		Expect(err).To(HaveOccurred())
	})

	It("barks when making a request while using a client with unparsable BaseURL", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL("as.lfdna,smdnasd:::"), // a keysmash, added ::: to have it a really unparsable URL
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}
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

		o := apiTestObject{"identifier"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("Hello from apiTestObject!"))
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

		o := apiTestObject{"identifier"}
		err = api.Create(ctx, &o)
		Expect(err).To(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("Hello from apiTestObject!"))
	})

	It("handles the Object returning an error on EndpointURL", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))

		err = api.List(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object returning an error on GetIdentifier", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing-identifier"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object returning an empty identifier on EndpointURL for operations requiring an IdentifiedObject", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{""}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(ErrUnidentifiedObject))
	})

	It("handles the Object returning an error on FilterAPIRequest", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_filter_request"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object returning an error on FilterAPIRequestBody", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_filter_request_body"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object returning a request body that cannot be encoded in json on FilterAPIRequestBody", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"function_filter_request_body"}
		err = api.Create(context.TODO(), &o)

		var e *json.UnsupportedTypeError
		Expect(errors.As(err, &e)).To(BeTrue())
	})

	It("handles the Object returning an error on FilterAPIResponse", func() {
		server.SetAllowUnhandledRequests(true)

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_filter_response"}
		err = api.Create(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object returning an error on DecodeAPIResponse", func() {
		server.AppendHandlers(
			ghttp.RespondWithJSONEncoded(200, apiTestObject{"indentifier"}),
		)

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_decode_response"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Engine returning a weird response content-type", func() {
		server.AppendHandlers(ghttp.RespondWith(200, "randomgarbage", http.Header{"Content-Type": []string{"application/octet-stream"}}))

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(ErrUnsupportedResponseFormat))
	})

	It("does not crash on Engine responses without body", func() {
		server.AppendHandlers(ghttp.RespondWith(204, nil))

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}
		err = api.Destroy(context.TODO(), &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("handles the Engine returning bad responses for List requests", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `foo no json`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWithJSONEncoded(200, []map[string]string{{"value": "hello world"}}),
			ghttp.RespondWith(200, `foo no json`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"foo": "hello world"}),
		)

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}

		var pi types.PageInfo
		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))

		var e *json.SyntaxError
		Expect(errors.As(err, &e)).To(BeTrue())

		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))
		Expect(err).NotTo(HaveOccurred())

		var os []apiTestObject
		ok := pi.Next(&os)
		Expect(pi.Error()).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		ok = pi.Next(&os)
		Expect(ok).To(BeFalse())
		Expect(errors.As(pi.Error(), &e)).To(BeTrue())

		err = api.List(context.TODO(), &o, Paged(1, 1, &pi))
		Expect(err).To(MatchError(ErrPageResponseNotSupported))
	})

	Context("configured to use a mock server", func() {
		type response struct {
			status int
			path   string
			query  string
			data   interface{}
		}

		var responses []response

		var api API

		JustBeforeEach(func() {
			for _, r := range responses {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", r.path, r.query),
					ghttp.RespondWithJSONEncoded(r.status, r.data),
				))
			}

			var err error
			api, err = NewAPI(
				WithLogger(logger),
				WithClientOptions(
					client.BaseURL(server.URL()),
					client.IgnoreMissingToken(),
				),
			)
			Expect(err).NotTo(HaveOccurred())
		})

		commonCheck := func(fullObjects bool) {
			namePrefix := ""

			if fullObjects {
				namePrefix = "full "
			}

			It("returns correct data with List operation used with pagination", func() {
				o := apiTestObject{}

				var pi types.PageInfo
				// we use the same page size as for channel to make testing easier and to not have to
				// provide the Paged option when using a channel
				err := api.List(context.TODO(), &o, Paged(1, ListChannelDefaultPageSize, &pi), FullObjects(fullObjects))
				Expect(err).NotTo(HaveOccurred())
				Expect(pi.CurrentPage()).To(BeEquivalentTo(0))

				var objects []apiTestObject
				for pi.Next(&objects) {
					Expect(objects).To(HaveLen(2))

					switch pi.CurrentPage() {
					case 1:
						Expect(objects).To(BeEquivalentTo([]apiTestObject{
							{namePrefix + "foo 1"},
							{namePrefix + "foo 2"},
						}))
					case 2:
						Expect(objects).To(BeEquivalentTo([]apiTestObject{
							{namePrefix + "foo 3"},
							{namePrefix + "foo 4"},
						}))
					case 3:
						Expect(objects).To(BeEquivalentTo([]apiTestObject{}))
					default:
						Fail("Unexpected current page")
					}
				}

				Expect(pi.Error()).NotTo(HaveOccurred())
				Expect(pi.CurrentPage()).To(BeEquivalentTo(3))
			})

			It("returns correct data with List operation used with a channel", func() {
				o := apiTestObject{}

				channel := make(types.ObjectChannel)
				err := api.List(context.TODO(), &o, ObjectChannel(&channel), FullObjects(fullObjects))
				Expect(err).NotTo(HaveOccurred())

				i := 0
				for retriever := range channel {
					i++

					err := retriever(&o)
					Expect(err).NotTo(HaveOccurred())

					switch i {
					case 1:
						Expect(o.Val).To(Equal(namePrefix + "foo 1"))
					case 2:
						Expect(o.Val).To(Equal(namePrefix + "foo 2"))
					case 3:
						Expect(o.Val).To(Equal(namePrefix + "foo 3"))
					case 4:
						Expect(o.Val).To(Equal(namePrefix + "foo 4"))
					default:
						Fail("Unexpected number of objects retrieved")
					}
				}
			})
		}

		Context("FullObjects disabled", func() {
			BeforeEach(func() {
				responses = []response{
					{200, "/resource/v1", "page=1&limit=10", []apiTestObject{{"foo 1"}, {"foo 2"}}},
					{200, "/resource/v1", "page=2&limit=10", []apiTestObject{{"foo 3"}, {"foo 4"}}},
					{200, "/resource/v1", "page=3&limit=10", []apiTestObject{}},
				}
			})

			commonCheck(false)
		})

		Context("FullObjects enabled", func() {
			Context("requests all succeeding", func() {
				BeforeEach(func() {
					responses = []response{
						{200, "/resource/v1", "page=1&limit=10", []apiTestObject{{"foo 1"}, {"foo 2"}}},

						{200, "/resource/v1/foo 1", "", apiTestObject{"full foo 1"}},
						{200, "/resource/v1/foo 2", "", apiTestObject{"full foo 2"}},

						{200, "/resource/v1", "page=2&limit=10", []apiTestObject{{"foo 3"}, {"foo 4"}}},

						{200, "/resource/v1/foo 3", "", apiTestObject{"full foo 3"}},
						{200, "/resource/v1/foo 4", "", apiTestObject{"full foo 4"}},

						{200, "/resource/v1", "page=3&limit=10", []apiTestObject{}},
					}
				})

				commonCheck(true)
			})

			Context("list request succeeding but get request failing", func() {
				BeforeEach(func() {
					responses = []response{
						{200, "/resource/v1", "page=1&limit=10", []apiTestObject{{"foo 1"}, {"foo 2"}}},

						{400, "/resource/v1/foo 1", "", map[string]string{"error": "something went wrong"}}, // a very realistic error, sadly.
					}
				})

				It("returns the error via page iterator", func() {
					o := apiTestObject{}

					var pi types.PageInfo
					err := api.List(context.TODO(), &o, Paged(1, 10, &pi), FullObjects(true))
					Expect(err).NotTo(HaveOccurred())
					Expect(pi.CurrentPage()).To(BeEquivalentTo(0))

					var objects []apiTestObject
					Expect(pi.Next(&objects)).To(BeFalse())
					Expect(pi.Error()).To(HaveOccurred())
					Expect(pi.CurrentPage()).To(BeEquivalentTo(0))
				})

				It("returns the error via channel", func() {
					o := apiTestObject{}
					ctx, cancel := context.WithCancel(context.TODO())

					var c types.ObjectChannel
					err := api.List(ctx, &o, ObjectChannel(&c), FullObjects(true))
					Expect(err).NotTo(HaveOccurred())

					retriever := <-c
					Expect(retriever(&o)).To(HaveOccurred())
					cancel()
				})
			})

			Context("list request succeeding but decoding fails", func() {
				BeforeEach(func() {
					responses = []response{
						{200, "/resource/v1", "page=1&limit=10", []apiTestObject{{"foo 1"}, {"foo 2"}}},

						{200, "/resource/v1/foo 1", "", "foo"},
					}
				})

				It("returns the error via page iterator", func() {
					o := apiTestObject{}

					var pi types.PageInfo
					err := api.List(context.TODO(), &o, Paged(1, 10, &pi), FullObjects(true))
					Expect(err).NotTo(HaveOccurred())
					Expect(pi.CurrentPage()).To(BeEquivalentTo(0))

					var objects []apiTestObject
					Expect(pi.Next(&objects)).To(BeFalse())
					Expect(pi.Error()).To(HaveOccurred())
					Expect(pi.CurrentPage()).To(BeEquivalentTo(0))
				})

				It("returns the error via channel", func() {
					o := apiTestObject{}
					ctx, cancel := context.WithCancel(context.TODO())

					var c types.ObjectChannel
					err := api.List(ctx, &o, ObjectChannel(&c), FullObjects(true))
					Expect(err).NotTo(HaveOccurred())

					retriever := <-c
					Expect(retriever(&o)).To(HaveOccurred())
					cancel()
				})
			})
		})
	})

	It("handles users trying to list with page iterator and channel simultaneously", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}

		var pi types.PageInfo
		oc := make(types.ObjectChannel)
		err = api.List(context.TODO(), &o, Paged(1, 2, &pi), ObjectChannel(&oc))
		Expect(err).To(MatchError(ErrCannotListChannelAndPaged))
	})

	It("handles users trying to list page 0", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.RespondWithJSONEncoded(200, []string{}),
				func(res http.ResponseWriter, req *http.Request) {
					Expect(req.URL.Query().Get("page")).To(Equal("1"))
				},
			),
		)

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

		o := apiTestObject{"identifier"}

		var pi types.PageInfo
		err = api.List(context.TODO(), &o, Paged(0, 2, &pi))
		Expect(err).NotTo(HaveOccurred())

		Expect(log.String()).To(ContainSubstring("requesting page 0, fixing to page 1"))
	})

	It("handles the Object returning an error on HasPagination", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_has_pagination"}
		err = api.List(context.TODO(), &o, Paged(1, 10, nil))
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles the Object not being able to be Listed with pagination", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"}]`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"}]`, http.Header{"Content-Type": []string{"application/json"}}),
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"}]`, http.Header{"Content-Type": []string{"application/json"}}),
		)

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"no_pagination"}

		var pi types.PageInfo
		err = api.List(context.TODO(), &o, Paged(1, 10, &pi))
		Expect(err).NotTo(HaveOccurred())

		var os []apiTestObject
		Expect(pi.Next(&os)).To(BeTrue())
		Expect(os).To(HaveLen(2))

		Expect(pi.Next(&os)).To(BeFalse())

		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	It("can list objects with channel", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"}]`, http.Header{"Content-Type": []string{"application/json"}}),
		)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"no_pagination"}

		var ch types.ObjectChannel
		err = api.List(context.TODO(), &o, ObjectChannel(&ch))
		Expect(err).NotTo(HaveOccurred())

		i := 0
		for retriever := range ch {
			err = retriever(&o)
			Expect(err).NotTo(HaveOccurred())

			switch i {
			case 0:
				Expect(o.Val).To(Equal("foo"))
			case 1:
				Expect(o.Val).To(Equal("bar"))
			default:
				Fail("unexpected number of objects")
			}

			i++
		}

		Expect(i).To(Equal(2))
		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	It("listing objects with channel handles decode errors", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"}]`, http.Header{"Content-Type": []string{"application/json"}}),
		)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"failing_decode_response"}

		ctx, cancel := context.WithCancel(context.TODO())
		var ch types.ObjectChannel
		err = api.List(ctx, &o, ObjectChannel(&ch))
		Expect(err).NotTo(HaveOccurred())

		retriever := <-ch
		cancel()

		err = retriever(&o)
		Expect(err).To(MatchError(errAPITest))

		Eventually(ch).Should(BeClosed())
		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	It("can abort listing objects with channel", func() {
		server.AppendHandlers(
			ghttp.RespondWith(200, `[{"value":"foo"},{"value":"bar"},{"value":"baz"},{"value":"bla"}]`, http.Header{"Content-Type": []string{"application/json"}}),
		)

		api, err := NewAPI(
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"no_pagination"}

		ctx, cancel := context.WithCancel(context.TODO())
		var ch types.ObjectChannel
		err = api.List(ctx, &o, ObjectChannel(&ch))
		Expect(err).NotTo(HaveOccurred())

		// to prevent having another retriever pushed to the channel we have to
		// cancel before running the retriever
		retriever := <-ch
		cancel()

		err = retriever(&o)
		Expect(err).NotTo(HaveOccurred())
		Expect(o.Val).To(Equal("foo"))

		Eventually(ch).Should(BeClosed())
		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	It("handles http.Client.Do() returning an error", func() {
		hc := http.Client{
			Transport: apiTestErrorRoundtripper(false),
		}

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.WithClient(&hc),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}
		err = api.Get(context.TODO(), &o)
		Expect(err).To(MatchError(errAPITest))
	})

	It("handles not being given a context", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		// the two nolint comments are for passing nil context, which is the behavior we want to test here.

		o := apiTestObject{"identifier"}
		err = api.Get(nil, &o) //nolint:golint,staticcheck
		Expect(err).To(MatchError(ErrContextRequired))

		err = api.List(nil, &o) //nolint:golint,staticcheck
		Expect(err).To(MatchError(ErrContextRequired))
	})

	It("handles bogus operations", func() {
		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"identifier"}
		req, err := api.(defaultAPI).makeRequest(context.TODO(), &o, &o, types.Operation("bogus operation"))
		Expect(err).To(MatchError(ErrOperationNotSupported))
		Expect(req).To(BeNil())
	})

	It("consumes the given options for all operations", func() {
		opt := apiTestAnyopOption("hello world")
		ctx := context.WithValue(context.TODO(), errAPITest, opt)

		server.AppendHandlers(
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, []map[string]string{{"value": "option-check"}}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{"value": "option-check"}),
			ghttp.RespondWithJSONEncoded(200, map[string]string{}),
		)

		api, err := NewAPI(
			WithLogger(logger),
			WithClientOptions(
				client.BaseURL(server.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).NotTo(HaveOccurred())

		o := apiTestObject{"option-check"}

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

const contextTestObjectBaseurl = "/v1/context_test_object"

type contextTestObject struct {
	Test string `anxcloud:"identifier"`

	endpointURLCalled    bool
	filterRequestCalled  bool
	filterResponseCalled bool
	requestBodyCalled    bool
	responseBodyCalled   bool
}

func (o *contextTestObject) GetIdentifier(context.Context) (string, error) {
	return o.Test, nil
}

func (o contextTestObject) checkContext(ctx context.Context, hasURL bool) {
	switch o.Test {
	case "Hooks":
		// nothing to do
	case "Operation":
		Expect(types.OperationFromContext(ctx)).To(Equal(types.OperationUpdate))
	case "Options":
		Expect(types.OptionsFromContext(ctx)).To(BeAssignableToTypeOf(&types.UpdateOptions{}))
	case "URL":
		u, err := types.URLFromContext(ctx)

		if hasURL {
			Expect(err).NotTo(HaveOccurred())
			Expect(u).NotTo(BeZero())
		} else {
			Expect(err).To(MatchError(types.ErrContextKeyNotSet))
			Expect(u).To(BeZero())
		}
	default:
		Fail(fmt.Sprintf("Unknown property to test: %v", o.Test))
	}
}

func (o *contextTestObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	o.checkContext(ctx, false)
	o.endpointURLCalled = true

	return url.Parse(contextTestObjectBaseurl)
}

func (o *contextTestObject) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	o.checkContext(ctx, true)
	o.filterRequestCalled = true

	return req, nil
}

func (o *contextTestObject) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	o.checkContext(ctx, true)
	o.filterResponseCalled = true

	return res, nil
}

func (o *contextTestObject) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	o.checkContext(ctx, true)
	o.requestBodyCalled = true

	return o, nil
}

func (o *contextTestObject) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	o.checkContext(ctx, true)
	o.responseBodyCalled = true
	return json.NewDecoder(data).Decode(o)
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

	It("has all hooks called on it", func() {
		o := contextTestObject{"Hooks", false, false, false, false, false}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("%v/%v", contextTestObjectBaseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Update(ctx, &o)
		Expect(err).NotTo(HaveOccurred())

		Expect(o.endpointURLCalled).To(BeTrue())
		Expect(o.filterRequestCalled).To(BeTrue())
		Expect(o.filterResponseCalled).To(BeTrue())
		Expect(o.requestBodyCalled).To(BeTrue())
		Expect(o.responseBodyCalled).To(BeTrue())
	})

	It("has operation in context for every method call", func() {
		o := contextTestObject{"Operation", false, false, false, false, false}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("%v/%v", contextTestObjectBaseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Update(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has options in context for every method call", func() {
		o := contextTestObject{"Options", false, false, false, false, false}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("%v/%v", contextTestObjectBaseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Update(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has URL in context for every method call except EndpointURL", func() {
		o := contextTestObject{"URL", false, false, false, false, false}

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("%v/%v", contextTestObjectBaseurl, o.Test)),
			ghttp.RespondWithJSONEncoded(200, o),
		))

		err := api.Update(ctx, &o)
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestAPIUnits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "api unit test suite")
}
