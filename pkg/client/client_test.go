package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/go-logr/stdr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("client", func() {
	var c Client
	var s *Server
	var numRequests int

	var options []Option

	AfterEach(func() {
		if s != nil {
			Expect(len(s.ReceivedRequests())).To(Equal(numRequests))
		}
	})

	BeforeEach(func() {
		numRequests = 1
		options = make([]Option, 0, 2)
		options = append(options, IgnoreMissingToken())
	})

	JustBeforeEach(func() {
		s = NewServer()

		options = append(options,
			BaseURL(s.URL()),
		)

		var err error
		c, err = New(
			options...,
		)

		Expect(err).NotTo(HaveOccurred())
	})

	It("handles request and response correctly", func() {
		data := url.Values{}
		data.Set("foo", "bar")
		encData := data.Encode()

		s.AppendHandlers(CombineHandlers(
			VerifyFormKV("foo", "bar"),
			VerifyHeaderKV("Authorization", "sensible-value"),
			VerifyHeaderKV("Content-Length", strconv.Itoa(len(encData))),
			RespondWith(http.StatusOK, "bar"),
		))

		req, err := http.NewRequest(http.MethodPost, s.URL(), strings.NewReader(encData))
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "sensible-value")

		response, err := c.Do(req)
		Expect(err).NotTo(HaveOccurred())

		Expect(response.Body).NotTo(BeNil())

		body, err := ioutil.ReadAll(response.Body)

		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(Equal("bar"))
	})

	It("handles the request and correctly parses the error response", func() {
		data := url.Values{}
		data.Set("foo", "bar")
		encData := data.Encode()

		s.AppendHandlers(CombineHandlers(
			VerifyFormKV("foo", "bar"),
			VerifyHeaderKV("Authorization", "sensible-value"),
			VerifyHeaderKV("Content-Length", strconv.Itoa(len(encData))),
			VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
			RespondWithJSONEncoded(http.StatusBadRequest, map[string]string{
				"msg": "error message",
			}),
		))

		req, err := http.NewRequest(http.MethodPost, s.URL(), strings.NewReader(encData))
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "sensible-value")

		response, err := c.Do(req)

		Expect(err).To(HaveOccurred())
		Expect(response).NotTo(BeNil())

		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("uses the auto-generated User-Agent in requests", func() {
		expectedUA := fmt.Sprintf("go-anxcloud/%s (%s)", "snapshot", runtime.GOOS)
		s.AppendHandlers(CombineHandlers(
			VerifyHeaderKV("User-Agent", expectedUA),
		))

		req, err := http.NewRequest("GET", s.URL(), nil)
		Expect(err).NotTo(HaveOccurred())

		_, err = c.Do(req)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when configured not to parse errors", func() {
		BeforeEach(func() {
			options = append(options, ParseEngineErrors(false))
		})

		It("handles request and response correctly, without returning an error for http error responses", func() {
			data := url.Values{}
			data.Set("foo", "bar")
			encData := data.Encode()

			s.AppendHandlers(CombineHandlers(
				VerifyFormKV("foo", "bar"),
				VerifyHeaderKV("Authorization", "sensible-value"),
				VerifyHeaderKV("Content-Length", strconv.Itoa(len(encData))),
				VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				RespondWithJSONEncoded(http.StatusBadRequest, map[string]string{
					"msg": "error message",
				}),
			))

			req, err := http.NewRequest(http.MethodPost, s.URL(), strings.NewReader(encData))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization", "sensible-value")

			response, err := c.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())

			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

			bodyData := map[string]string{}

			err = json.NewDecoder(response.Body).Decode(&bodyData)
			Expect(err).NotTo(HaveOccurred())

			Expect(bodyData["msg"]).To(Equal("error message"))
		})
	})

	Context("when configured with an authorization token", func() {
		const dummyToken = "ie7dois8Ooquoo1ieB9kae8Od9ooshee3nejuach4inae3gai0Re0Shaipeihail" //nolint:gosec // Not a real token.

		BeforeEach(func() {
			options = append(options, TokenFromString(dummyToken))
		})

		It("sends the authorization token with the request", func() {
			s.AppendHandlers(CombineHandlers(
				VerifyHeaderKV("Authorization", fmt.Sprintf("Token %v", dummyToken)),
			))

			req, err := http.NewRequest("GET", s.URL(), nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = c.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when configured with an explicit User-Agent", func() {
		const testUA = "Firefox 42.0" //nolint:gosec // Not a real token.

		BeforeEach(func() {
			options = append(options, UserAgent(testUA))
		})

		It("sends the correct User-Agent with requests", func() {
			s.AppendHandlers(CombineHandlers(
				VerifyHeaderKV("User-Agent", testUA),
			))

			req, err := http.NewRequest("GET", s.URL(), nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = c.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("request and response logging", func() {
		var logMessages strings.Builder
		BeforeEach(func() {
			logMessages = strings.Builder{}
		})

		commonLogTest := func(shouldContain bool) {
			It("logs request and response", func() {
				s.AppendHandlers(VerifyRequest("GET", "/"))

				req, err := http.NewRequest("GET", s.URL(), nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = c.Do(req)
				Expect(err).NotTo(HaveOccurred())

				checks := And(
					ContainSubstring(`Sending request to Engine`),
					ContainSubstring(`"method"="GET"`),

					ContainSubstring(`Received response from Engine`),
					ContainSubstring(`"body"=""`),
					ContainSubstring(`"method"="GET"`),
					ContainSubstring(`"statusCode"=200`),
				)

				if shouldContain {
					Expect(logMessages.String()).To(checks, "Log should contain request and response logs")
				} else {
					Expect(logMessages.String()).NotTo(checks, "Log should not contain request and response logs")
				}
			})
		}

		Context("when configured with the LogWriter option", func() {
			BeforeEach(func() {
				options = append(options, LogWriter(&logMessages))
			})

			It("logs the deprecation warning to the deprecated logger", func() {
				numRequests = 0

				log := strings.Builder{}

				opts := clientOptions{}
				err := LogWriter(&log)(&opts)
				Expect(err).NotTo(HaveOccurred())

				Expect(log.String()).To(ContainSubstring("The LogWriter option of github.com/anexia-it/go-anxcloud/pkg/client is deprecated."))
			})

			commonLogTest(true)
		})

		Context("when configured with the Logger option", func() {
			BeforeEach(func() {
				logger := stdr.New(
					log.New(&logMessages, "", log.Lmsgprefix|log.Lshortfile),
				)

				options = append(options, Logger(logger))
			})

			commonLogTest(false)

			Context("log verbosity set to LogVerbosityRequests", func() {
				BeforeEach(func() {
					stdr.SetVerbosity(LogVerbosityRequests)
				})

				commonLogTest(true)
			})
		})
	})

	Context("not configured to IgnoreMissingToken", func() {
		It("returns an error on creating the instance", func() {
			client, err := New(BaseURL(s.URL()))
			Expect(err).To(MatchError(ErrConfiguration))
			Expect(err.Error()).To(ContainSubstring("token not set"))
			Expect(client).To(BeNil())

			numRequests = 0
		})
	})
})

var _ = Describe("AuthFromEnv option", func() {
	var unset bool
	var option Option

	JustBeforeEach(func() {
		option = AuthFromEnv(unset)
	})

	commonTest := func() {
		It("loads the token from environment", func() {
			os.Setenv("ANEXIA_TOKEN", "foo bar baz")

			opts := clientOptions{}
			err := option(&opts)

			Expect(opts.token).To(Equal("foo bar baz"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("errors when environment variable is not set", func() {
			os.Unsetenv("ANEXIA_TOKEN")

			opts := clientOptions{}
			err := option(&opts)

			Expect(err).To(MatchError(ErrEnvMissing))
		})
	}

	Context("configured to unset", func() {
		BeforeEach(func() {
			unset = true
		})

		commonTest()

		It("unsets the environment variable", func() {
			os.Setenv("ANEXIA_TOKEN", "foo bar baz")
			opts := clientOptions{}
			err := option(&opts)
			Expect(err).NotTo(HaveOccurred())

			_, present := os.LookupEnv("ANEXIA_TOKEN")
			Expect(present).To(BeFalse())
		})
	})

	Context("configured not to unset", func() {
		BeforeEach(func() {
			unset = false
		})

		commonTest()
	})
})

func TestClientSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "client test suite")
}
