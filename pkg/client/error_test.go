package client

import (
	"encoding/json"
	"errors"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseError", func() {
	It("parses the valid test error response", func() {
		resW := httptest.NewRecorder()
		resW.WriteHeader(401)
		resW.Header().Add("Content-Type", "application/json; charset=utf-8")
		_, _ = resW.Write([]byte(`{ "error": { "code": 401, "message": "Something went wrong. Please contact support.", "validation": {} } }`))

		req := httptest.NewRequest("GET", "https://engine.anexia.com/api/foo/bar", nil)
		res := resW.Result()

		err := parseEngineError(req, res)
		Expect(err).To(HaveOccurred())

		resError := &ResponseError{}
		Expect(errors.As(err, &resError)).To(BeTrue())

		Expect(err.Error()).To(ContainSubstring("error from api:"))
		Expect(err.Error()).To(ContainSubstring("Code:401"))
		Expect(err.Error()).To(ContainSubstring("Message:Something went wrong. Please contact support."))
		statusError, ok := err.(interface{ StatusCode() int })
		Expect(ok).To(BeTrue())
		Expect(statusError.StatusCode()).To(Equal(401))
	})

	It("keeps the status for an empty error response", func() {
		resW := httptest.NewRecorder()
		resW.WriteHeader(404)

		req := httptest.NewRequest("GET", "https://engine.anexia.com/api/foo/bar", nil)
		res := resW.Result()

		err := parseEngineError(req, res)
		Expect(err).To(HaveOccurred())

		resError := &ResponseError{}
		Expect(errors.As(err, &resError)).To(BeTrue())
		statusError, ok := err.(interface{ StatusCode() int })
		Expect(ok).To(BeTrue())
		Expect(statusError.StatusCode()).To(Equal(404))
		Expect(resError.ErrorData.Code).To(Equal(404))
		Expect(resError.ErrorData.Message).To(Equal("404 Not Found"))
	})

	It("keeps the status for an HTML error response", func() {
		resW := httptest.NewRecorder()
		resW.WriteHeader(403)
		_, _ = resW.Write([]byte(`<html><body>Forbidden</body></html>`))

		req := httptest.NewRequest("GET", "https://engine.anexia.com/api/foo/bar", nil)
		res := resW.Result()

		err := parseEngineError(req, res)
		Expect(err).To(HaveOccurred())

		resError := &ResponseError{}
		Expect(errors.As(err, &resError)).To(BeTrue())
		statusError, ok := err.(interface{ StatusCode() int })
		Expect(ok).To(BeTrue())
		Expect(statusError.StatusCode()).To(Equal(403))
		Expect(resError.ErrorData.Code).To(Equal(403))
		Expect(resError.ErrorData.Message).To(Equal("403 Forbidden"))

		jsonSyntaxErr := &json.SyntaxError{}
		Expect(errors.As(err, &jsonSyntaxErr)).To(BeTrue())
	})

	It("uses the response status when the error code is omitted", func() {
		resW := httptest.NewRecorder()
		resW.WriteHeader(500)
		_, _ = resW.Write([]byte(`{"error":{"message":"Internal Server Error"}}`))

		req := httptest.NewRequest("GET", "https://engine.anexia.com/api/foo/bar", nil)
		err := parseEngineError(req, resW.Result())

		resError := &ResponseError{}
		Expect(errors.As(err, &resError)).To(BeTrue())
		Expect(resError.ErrorData.Code).To(Equal(500))
	})
})
