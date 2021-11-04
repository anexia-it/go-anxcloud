package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/go-logr/logr"
)

// maxPageFetchRetry is the maximum number of retries to fetch a single page.
// When that retry count is reached, an error is set that cannot be cleared with ResetError().
const maxPageFetchRetry = 10

type pageFetcher func(page uint) (json.RawMessage, error)

type pageIter struct {
	currentPage  uint
	totalPages   uint
	totalItems   uint
	itemsPerPage uint

	err             error
	errRetryCounter uint

	pageFetcher pageFetcher

	ctx context.Context
}

// CurrentPage returns the page number the last Next call processed.
func (p *pageIter) CurrentPage() uint {
	return p.currentPage
}

// TotalPages returns the total number of pages. Note: not all APIs support this and will then return 0.
func (p *pageIter) TotalPages() uint {
	return p.totalPages
}

// TotalItems returns the total number of items. Note: not all APIs support this and will then return 0.
func (p *pageIter) TotalItems() uint {
	return p.totalItems
}

// ItemsPerPage returns the maximum number of entries per page, corresponding to the Limit parameter given
// to the Paged attribute.
func (p *pageIter) ItemsPerPage() uint {
	return p.itemsPerPage
}

// Next retrieves the next page of objects to process. On the first call, it gives the exact same
// objects as api.List() returned to allow iterating over all pages easily. It returns true when it
// has received another page of objects and false on completion or error. Errors can be retrieved by
// calling PageInfo.Error().
func (p *pageIter) Next(objects interface{}) bool {
	if p.err != nil {
		return false
	}

	val := reflect.ValueOf(objects)

	isPointer := val.Kind() == reflect.Ptr
	isArrayOrSlice := false
	isObjects := false
	isRawMessages := false

	wrongType := val.Type()

	if isPointer {
		kind := val.Elem().Kind()
		isArrayOrSlice = kind == reflect.Slice || kind == reflect.Array
	}

	if isArrayOrSlice {
		objectType := reflect.TypeOf((*types.Object)(nil)).Elem()

		elementType := val.Elem().Type().Elem()
		ptrToElementType := reflect.PtrTo(elementType)

		isObjects = elementType.Implements(objectType) || ptrToElementType.Implements(objectType)
		isRawMessages = elementType == reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	}

	// the check for isObjects || isRawMessages isn't actually required, but is kept to prevent users decoding their
	// page of objects into something completely different by accident. I currently don't see a valid reason to do
	// that, but if one comes up, this can probably be removed. -- Mara @LittleFox94 Grosch, 2021-10-16
	// json.RawMessage is allowed for retrieving objects via channel, where the page is decoded into an array of
	// json.RawMessage and every entry of that is decoded into the target object as soon as it is needed.
	if !isPointer || !isArrayOrSlice || (!isObjects && !isRawMessages) {
		p.err = fmt.Errorf("%w: the argument given to PageInfo.Next() must be a pointer to []T where T or *T implements types.Object or T is json.RawMessage; expected *[]T, you gave %v", ErrTypeNotSupported, wrongType)
		return false
	}

	pageData, err := p.pageFetcher(p.currentPage + 1)
	if err != nil {
		p.errRetryCounter++
		p.err = err
		return false
	}

	_, _, _, _, data, err := decodePaginationResponseBody(pageData, types.ListOptions{Page: p.currentPage + 1, EntriesPerPage: p.itemsPerPage})
	if err != nil {
		p.errRetryCounter++
		p.err = err
		return false
	}

	d := json.NewDecoder(bytes.NewBuffer(data))
	d.DisallowUnknownFields()

	if err := d.Decode(objects); err != nil {
		p.errRetryCounter++
		p.err = err
		return false
	}

	log := logr.FromContextOrDiscard(p.ctx)

	retrievedElements := uint(val.Elem().Len())
	if retrievedElements > p.itemsPerPage && p.itemsPerPage > 0 {
		log.Info("Retrieved more elements in one Next() than wanted", "wanted", p.itemsPerPage, "retrieved", retrievedElements)
	} else {
		log.V(1).Info("Retrieved elements from engine", "limit", p.itemsPerPage, "retrieved", retrievedElements)
	}

	p.errRetryCounter = 0
	p.currentPage++

	return retrievedElements > 0
}

// Returns error. An iteration over all pages has successfully completed when Next() returns false and
// Error() returns nil. You should check for errors after Next() returns false to differentiate between
// "all pages done" and "error retrieving page".
func (p *pageIter) Error() error {
	return p.err
}

// ResetError clears any stored error to resume the iterator. If the retry counter for the current page exceeded
// a package-defined maximum, the error cannot be cleared and Error() will return it after ResetError() was called.
// you have to check for this.
func (p *pageIter) ResetError() {
	if p.errRetryCounter < maxPageFetchRetry {
		p.err = nil
	}
}

func newPageIter(ctx context.Context, responseBody json.RawMessage, opts types.ListOptions, fetcher pageFetcher) (types.PageInfo, error) {
	if logger, err := logr.FromContext(ctx); err == nil {
		ctx = logr.NewContext(ctx, logger.WithName("pagination"))
	}

	ret := pageIter{
		ctx: ctx,
	}

	currentPage, limit, totalPages, totalItems, data, err := decodePaginationResponseBody(responseBody, opts)
	if err != nil {
		return nil, err
	}

	if currentPage == 1 {
		currentPage = 0
	}

	ret.currentPage = currentPage
	ret.itemsPerPage = limit
	ret.totalPages = totalPages
	ret.totalItems = totalItems

	// first pageFetcher is returning the data we got with the initial request, after this is fetched, we
	// use the pageFetcher provided as argument
	ret.pageFetcher = func(page uint) (json.RawMessage, error) {
		ret.pageFetcher = fetcher
		return data, nil
	}

	return &ret, nil
}

func decodePaginationResponseBody(data json.RawMessage, opts types.ListOptions) (page, limit, totalPages, totalItems uint, ret json.RawMessage, err error) {
	page = 0
	limit = 0
	totalPages = 0
	totalItems = 0
	ret = json.RawMessage{}

	// TODO(LittleFox94): this is not the same for every API and we need a way to override this or
	// find the X ways it's done and have options for that. Currently we support those two types and
	// "plain data array".

	type dataResponse struct {
		CurrentPage    uint `json:"page"`
		TotalPages     uint `json:"total_pages"`
		TotalItems     uint `json:"total_items"`
		EntriesPerPage uint `json:"limit"`

		Data json.RawMessage `json:"data"`
	}

	type dataDataResponse struct {
		State    string       `json:"state"`
		Messages []string     `json:"messages"`
		Data     dataResponse `json:"data"`
	}

	// First dataData then data is important since we switch over the index of the decoded message,
	// set data from dataData and fallthrough.
	// The entries have to be pointers, else every entry matches every data - since it is an interface{} then.
	responseTypes := []interface{}{&dataDataResponse{}, &dataResponse{}, &[]json.RawMessage{}}
	actualResponse := -1

	for i, response := range responseTypes {
		decoder := json.NewDecoder(bytes.NewBuffer(data))

		// in case we receive a completely different response we have to prevent it being decodable into one
		// of the supported formats by accident.
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&response); err == nil {
			actualResponse = i
			break
		}
	}

	if actualResponse == -1 {
		err = ErrPageResponseNotSupported
		return
	}

	switch actualResponse {
	case 0:
		responseTypes[1] = &responseTypes[0].(*dataDataResponse).Data
		fallthrough
	case 1:
		data := responseTypes[1].(*dataResponse)
		page = data.CurrentPage - 1 // Next increments the page before trying to retrieve it
		limit = data.EntriesPerPage
		totalPages = data.TotalPages
		totalItems = data.TotalItems
		ret = data.Data
	case 2:
		page = opts.Page
		limit = opts.EntriesPerPage

		ret, err = json.Marshal(responseTypes[2])
	}

	return page, limit, totalPages, totalItems, ret, err
}
