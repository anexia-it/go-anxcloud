package api

import (
	"context"
	"encoding/json"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type paginationTestObject struct {
	Message string `json:"message"`
}

func (p paginationTestObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return nil, ErrOperationNotSupported
}

func (p paginationTestObject) GetIdentifier(context.Context) (string, error) {
	return p.Message, nil
}

var _ = Describe("PageInfo implementation pageIter", func() {
	var responses []json.RawMessage
	var afterErrorResponse json.RawMessage

	var pi types.PageInfo
	var piCreateError error

	JustBeforeEach(func() {
		returnedErrors := 0

		expectedPage := 2
		fetcher := func(page uint) (json.RawMessage, error) {
			Expect(page).To(BeEquivalentTo(expectedPage))

			if page > uint(len(responses)) {
				return json.RawMessage("[]"), nil
			}

			if responses[page-1] == nil {
				returnedErrors++

				switch returnedErrors {
				case 1:
					return nil, HTTPError{statusCode: 500, message: "Server error"}
				case 2:
					return json.RawMessage(`-öasjfn.ksdjfbksdnmf, sdf`), nil
				case 3:
					expectedPage++
					return afterErrorResponse, nil
				}
			} else {
				expectedPage++
			}

			return responses[page-1], nil
		}

		var responseBody json.RawMessage = nil

		if len(responses) > 0 {
			responseBody = responses[0]
		}

		opts := types.ListOptions{
			EntriesPerPage: 2,
		}

		pi, piCreateError = newPageIter(context.TODO(), nil, responseBody, opts, fetcher, false)
	})

	AssertCommonBehavior := func() {
		It("creates the iterator without error", func() {
			Expect(piCreateError).NotTo(HaveOccurred())
		})

		It("barks when argument given to Next is not as expected", func() {
			var out []string
			ok := pi.Next(&out)
			err := pi.Error()

			Expect(ok).To(BeFalse())
			Expect(err).To(MatchError(ErrTypeNotSupported))
			Expect(pi.CurrentPage()).To(BeEquivalentTo(0))
		})

		It("iterates through the pages until first error", func() {
			var out []json.RawMessage

			expectedPage := 1
			var page uint
			for pi.Next(&out) {
				page = pi.CurrentPage()
				Expect(page).To(BeEquivalentTo(expectedPage))
				expectedPage++
			}

			Expect(page).To(BeEquivalentTo(1))

			// pi still reports page 1 because that's the data in out
			Expect(pi.CurrentPage()).To(BeEquivalentTo(1))

			// This is the data of the first page. With the fetcher returning an error, this data
			// is not overwritten even though we already are on second page. We can use this behavior
			// for testing, but since it's not obvious, there is a comment here :)
			Expect(out).To(HaveLen(3))

			err := pi.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Server error"))
		})

		It("continues after error is cleared", func() {
			ok := true

			expectedPage := 1
			expectedError := 0

			for ok {
				var out []paginationTestObject
				ok = pi.Next(&out)

				if err := pi.Error(); !ok && err != nil {
					expectedError++

					switch expectedError {
					case 1:
						Expect(err.Error()).To(ContainSubstring("Server error"))
					case 2:
						Expect(err).To(MatchError(ErrPageResponseNotSupported))
					}

					Expect(pi.Next(&out)).To(BeFalse())

					pi.ResetError()
					ok = pi.Error() == nil
					continue
				}

				Expect(pi.CurrentPage()).To(BeEquivalentTo(expectedPage))
				expectedPage++
			}

			Expect(pi.Error()).NotTo(HaveOccurred())
			Expect(expectedError).To(Equal(2))
		})
	}

	// Just a plain json array with the results is returned.
	Context("with raw-array pages", func() {
		BeforeEach(func() {
			responses = []json.RawMessage{
				json.RawMessage(`[ { "message": "foo" }, { "message": "bar" }, { "message": "why not give one more than limit? :o)" } ]`),
				nil,
				json.RawMessage(`[ { "message": "baz" } ]`),
			}

			afterErrorResponse = json.RawMessage(`[ { "message": "server ok again" } ]`)
		})

		It("returns 0 for TotalPages and TotalItems and 2 for ItemsPerPage", func() {
			Expect(pi.TotalPages()).To(BeEquivalentTo(0))
			Expect(pi.TotalItems()).To(BeEquivalentTo(0))
			Expect(pi.ItemsPerPage()).To(BeEquivalentTo(2))
		})

		AssertCommonBehavior()
	})

	// A json object containing the current page, total page and item counts and a data array is returned.
	Context("with metadata page responses", func() {
		BeforeEach(func() {
			responses = []json.RawMessage{
				json.RawMessage(`{ "page": 1, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "foo" }, { "message": "bar" }, { "message": "why not give one more than limit? :o)" } ] }`),
				nil,
				json.RawMessage(`{ "page": 3, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "baz" } ] }`),
			}

			afterErrorResponse = json.RawMessage(`{ "page": 2, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "server ok again" } ] }`)
		})

		It("returns correct values for TotalPages, TotalItems and ItemsPerPage", func() {
			Expect(pi.TotalPages()).To(BeEquivalentTo(3))
			Expect(pi.TotalItems()).To(BeEquivalentTo(4))
			Expect(pi.ItemsPerPage()).To(BeEquivalentTo(2))
		})

		AssertCommonBehavior()
	})

	// A json object containing "state", a "messages" array and a "data" object containing the current page, total page and item counts and a data array is returned.
	// Since the actual data is at `data.data` we call this data.data response from now on.
	Context("with data.data page responses", func() {
		BeforeEach(func() {
			responses = []json.RawMessage{
				json.RawMessage(`{ "state": "success", "messages": [], "data": { "page": 1, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "foo" }, { "message": "bar" }, { "message": "why not give one more than limit? :o)" } ] } }`),
				nil,
				json.RawMessage(`{ "state": "success", "messages": [], "data": { "page": 3, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "baz" } ] } }`),
			}

			afterErrorResponse = json.RawMessage(`{ "state": "success", "messages": [], "data": { "page": 2, "total_pages": 3, "total_items": 4, "limit": 2, "data": [ { "message": "server ok again" } ] } }`)
		})

		It("returns correct values for TotalPages, TotalItems and ItemsPerPage", func() {
			Expect(pi.TotalPages()).To(BeEquivalentTo(3))
			Expect(pi.TotalItems()).To(BeEquivalentTo(4))
			Expect(pi.ItemsPerPage()).To(BeEquivalentTo(2))
		})

		AssertCommonBehavior()
	})

	Context("with unknown page response format", func() {
		BeforeEach(func() {
			responses = []json.RawMessage{
				json.RawMessage(`{ "alien_pages": 42, "alien_data": [] }`),
			}
		})

		It("returns an error on creating the iterator", func() {
			Expect(piCreateError).To(MatchError(ErrPageResponseNotSupported))
		})
	})
})
