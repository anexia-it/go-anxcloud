package mock

import (
	"errors"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var (
	// ErrPageSizeCannotBeZero is returned when newMockPageIter is getting called with itemsPerPage = 0 to prevent devision by zero panic
	ErrPageSizeCannotBeZero = errors.New("mockPageIter itemsPerPage cannot be zero")
)

type mockPageIter struct {
	page         uint
	itemsPerPage uint
	items        []types.Object
	err          error
}

// CurrentPage returns the page number the last Next call processed.
func (i *mockPageIter) CurrentPage() uint {
	return i.page
}

// TotalPages returns the total number of pages.
func (i *mockPageIter) TotalPages() uint {
	return 1 + (i.TotalItems()-1)/i.ItemsPerPage()
}

// TotalItems returns the total number of items.
func (i *mockPageIter) TotalItems() uint {
	return uint(len(i.items))
}

// ItemsPerPage returns the number of items per page.
func (i *mockPageIter) ItemsPerPage() uint {
	return i.itemsPerPage
}

// Next retrieves the next page of objects to process.
func (i *mockPageIter) Next(objects interface{}) bool {
	if i.Error() != nil || i.page > i.TotalPages() {
		return false
	}

	sliceFrom := uintMin((i.CurrentPage()-1)*i.ItemsPerPage(), i.TotalItems())
	sliceTo := uintMin(i.CurrentPage()*i.ItemsPerPage(), i.TotalItems())

	*objects.(*[]types.Object) = i.items[sliceFrom:sliceTo]
	i.page++
	return true
}

// Returns error.
func (i *mockPageIter) Error() error {
	return i.err
}

// ResetError clears any stored error to resume the iterator.
func (i *mockPageIter) ResetError() {
	i.err = nil
}

func newMockPageIter(items []types.Object, itemsPerPage, page uint) (types.PageInfo, error) {
	if itemsPerPage == 0 {
		return nil, ErrPageSizeCannotBeZero
	}
	return &mockPageIter{
		items:        items,
		page:         page,
		itemsPerPage: itemsPerPage,
	}, nil
}
