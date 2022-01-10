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
	})

	It("returns an error for invalid error responses", func() {
		resW := httptest.NewRecorder()
		resW.WriteHeader(500)
		resW.Header().Add("Content-Type", "application/json; charset=utf-8")
		_, _ = resW.Write([]byte(`not really json`))

		req := httptest.NewRequest("GET", "https://engine.anexia.com/api/foo/bar", nil)
		res := resW.Result()

		err := parseEngineError(req, res)
		Expect(err).To(HaveOccurred())

		jsonSyntaxErr := &json.SyntaxError{}
		Expect(errors.As(err, &jsonSyntaxErr)).To(BeTrue())
	})
})
