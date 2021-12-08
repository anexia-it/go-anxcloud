package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("logRequest and logResponse", func() {
	var fullLog strings.Builder

	var verbosity int
	var logger logr.Logger

	JustBeforeEach(func() {
		fullLog = strings.Builder{}

		logger = funcr.New(
			func(prefix, args string) {
				message := fmt.Sprintf("%s\t%s\n", prefix, args)
				fullLog.WriteString(message)
			},
			funcr.Options{
				Verbosity: verbosity,
			},
		)
	})

	Context("when configured for verbosity 3", func() {
		BeforeEach(func() {
			verbosity = 3
		})

		It("redacts the Authorization headers contents in request", func() {
			req := httptest.NewRequest("GET", "/foo", nil)
			req.Header.Add("Authorization", "auth_token")

			logRequest(req, logger)

			Expect(fullLog.String()).NotTo(ContainSubstring("auth_token"))
			Expect(fullLog.String()).To(ContainSubstring("REDACTED"))
		})

		It("redacts Set-Cookie header contents in response", func() {
			w := httptest.NewRecorder()

			cookie := http.Cookie{
				Name:  "Session-Cookie",
				Value: "session_id",
			}

			http.SetCookie(w, &cookie)

			written, err := w.Write([]byte("OK"))
			Expect(err).NotTo(HaveOccurred())
			Expect(written).To(Equal(2))

			logResponse(w.Result(), logger)

			Expect(fullLog.String()).NotTo(ContainSubstring("session_id"))
			Expect(fullLog.String()).To(ContainSubstring("REDACTED"))
		})

		It("does not mangle the request body", func() {
			req := httptest.NewRequest("GET", "/foo", bytes.NewBuffer([]byte("OK")))

			logRequest(req, logger)

			body, err := ioutil.ReadAll(req.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal([]byte("OK")))
		})

		It("does not crash on nil requests", func() {
			Expect(func() {
				logRequest(nil, logger)
			}).NotTo(Panic())
		})

		It("does not crash on nil responses", func() {
			Expect(func() {
				logResponse(nil, logger)
			}).NotTo(Panic())
		})
	})

	Context("when configured for verbosity below 2", func() {
		BeforeEach(func() {
			verbosity = 2
		})

		// technically this is testing only code our package ..
		It("detects the logger as disabled", func() {
			Expect(logger.V(3).Enabled()).To(BeFalse())
		})

		It("does not log anything", func() {
			req := httptest.NewRequest("GET", "/foo", bytes.NewBuffer([]byte("OK")))
			logRequest(req, logger)

			Expect(fullLog.String()).To(BeEmpty())
		})
	})
})

var _ = Describe("stringifyHeaders", func() {
	var headers string

	BeforeEach(func() {
		h := make(http.Header, 4)
		h.Add("Foxes", "are cool")

		// this should be sorted to the start
		h.Add("000", "1234")

		// add values for Foo in wrongly sorted order to test values are sorted
		h.Add("Foo", "Baz")
		h.Add("Foo", "Bar")

		headers = stringifyHeaders(h)
	})

	It("surrounds header values with single quotes", func() {
		Expect(headers).To(ContainSubstring("Foxes: 'are cool'"))
	})

	It("sorts headers by their name", func() {
		Expect(headers).To(ContainSubstring("000: '1234', Foo:"))
	})

	It("sorts header contents by their value", func() {
		Expect(headers).To(ContainSubstring("Foo: ['Bar', 'Baz']"))
	})

	It("stringifies headers fully as expected", func() {
		Expect(headers).To(Equal("000: '1234', Foo: ['Bar', 'Baz'], Foxes: 'are cool'"))
	})
})
