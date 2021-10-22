package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("client", func() {
	var c Client
	var s *Server

	AfterEach(func() {
		if s != nil {
			Expect(len(s.ReceivedRequests())).To(Equal(1))
		}
	})

	Context("when configured with default values", func() {
		BeforeEach(func() {
			s = NewServer()

			var err error
			c, err = New(
				IgnoreMissingToken(),
				BaseURL(s.URL()),
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
	})

	Context("when configured not to parse errors", func() {
		BeforeEach(func() {
			s = NewServer()

			var err error
			c, err = New(
				IgnoreMissingToken(),
				ParseEngineErrors(false),
				BaseURL(s.URL()),
			)

			Expect(err).NotTo(HaveOccurred())
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
			s = NewServer()

			var err error
			c, err = New(
				TokenFromString(dummyToken),
				ParseEngineErrors(false),
				BaseURL(s.URL()),
			)

			Expect(err).NotTo(HaveOccurred())
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
			s = NewServer()

			var err error
			c, err = New(
				UserAgent(testUA),
				IgnoreMissingToken(),
				ParseEngineErrors(false),
				BaseURL(s.URL()),
			)

			Expect(err).NotTo(HaveOccurred())
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
})

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "client suite")
}
