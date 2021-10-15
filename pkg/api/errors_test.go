package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExampleIgnoreNotFound() {
	api := newExampleAPI()

	backend := backend.Backend{Identifier: "non-existing identifier"}
	if err := api.Get(context.TODO(), &backend); IgnoreNotFound(err) != nil {
		fmt.Printf("Error retrieving backend from engine: %v\n", err)
	} else if err != nil {
		fmt.Printf("Requested backend does not exist\n")
	} else {
		fmt.Printf("Retrieved backend with name '%v'\n", backend.Name)
	}

	// Output:
	// Requested backend does not exist
}

var _ = Describe("HTTPError", func() {
	Context("when creating a HTTPError without custom message and without wrapping an error", func() {
		var err error

		BeforeEach(func() {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			rec.WriteHeader(500)

			err = newHTTPError(req, rec.Result(), nil, nil)

			he := HTTPError{
				message:    "Engine returned an error: 500 Internal Server Error (500)",
				statusCode: 500,
				url: &url.URL{
					Path: "/",
				},
				method: "GET",
			}

			Expect(err).To(MatchError(he))
		})

		It("it returns the status code", func() {
			var he HTTPError
			Expect(errors.As(err, &he)).To(BeTrue())
			Expect(he.StatusCode()).To(Equal(500))
		})

		It("it returns the expected message", func() {
			var he HTTPError
			Expect(errors.As(err, &he)).To(BeTrue())
			Expect(he.Error()).To(Equal("Engine returned an error: 500 Internal Server Error (500)"))
		})

		It("it does not wrap an EngineError", func() {
			var ee EngineError
			Expect(errors.As(err, &ee)).To(BeFalse())
		})
	})

	Context("when creating a HTTPError with wrapping an EngineError", func() {
		var err error

		BeforeEach(func() {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			rec.WriteHeader(500)

			err = newHTTPError(req, rec.Result(), ErrNotFound, nil)
		})

		It("it wraps the given EngineError error", func() {
			Expect(err).To(MatchError(ErrNotFound))
		})
	})

	Context("when creating a HTTPError with a custom error message", func() {
		var err error

		BeforeEach(func() {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			rec.WriteHeader(500)

			msg := "Random message for testing"
			err = newHTTPError(req, rec.Result(), nil, &msg)
		})

		It("it returns the correct message", func() {
			Expect(err.Error()).To(Equal("Random message for testing"))
		})
	})
})

var _ = Describe("errorFromResponse function", func() {
	req := httptest.NewRequest("GET", "/", nil)

	var statusCode int
	var res *http.Response

	JustBeforeEach(func() {
		rec := httptest.NewRecorder()
		rec.WriteHeader(statusCode)
		res = rec.Result()
	})

	Context("for status code 404", func() {
		BeforeEach(func() {
			statusCode = 404
		})

		It("returns ErrNotFound as expected", func() {
			err := errorFromResponse(req, res)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ErrNotFound))
		})
	})

	Context("for status code 403", func() {
		BeforeEach(func() {
			statusCode = 403
		})

		It("returns ErrAccessDenied as expected", func() {
			err := errorFromResponse(req, res)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ErrAccessDenied))
		})
	})

	Context("for status code 500", func() {
		BeforeEach(func() {
			statusCode = 500
		})

		It("returns a matching HTTPError as expected", func() {
			err := errorFromResponse(req, res)
			Expect(err).To(HaveOccurred())

			var he HTTPError
			ok := errors.As(err, &he)
			Expect(ok).To(BeTrue())

			Expect(he.StatusCode()).To(Equal(500))
		})
	})
})

var _ = Describe("EngineError", func() {
	It("returns the correct message", func() {
		Expect(ErrNotFound.Error()).To(Equal("requested resource does not exist on the engine"))
	})

	Context("when created without a wrapping error", func() {
		It("does not return a wrapped error", func() {
			Expect(ErrNotFound.Unwrap()).To(BeNil())
		})
	})
})
