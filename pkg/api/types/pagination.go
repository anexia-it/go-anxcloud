package types

// PageInfo contains information about the currently retrieved page and also a way to
// iterate over this and following pages.
type PageInfo interface {
	// CurrentPage returns the 1-based number of the last page processed in Next.
	CurrentPage() uint

	// TotalPages returns the total number of pages if supported by the given API, 0 otherwise.
	TotalPages() uint

	// TotalItems returns the total number of items if supported by the given API, 0 otherwise.
	TotalItems() uint

	// ItemsPerPage returns the desired number of items retrieved per page from the Engine.
	ItemsPerPage() uint

	// Next retrieves the data for the next page and stores the decoded values in the given pointer
	// to array of Object or json.RawMessage.
	Next(objects interface{}) bool

	// Error returns the error preventing Next to continue.
	Error() error

	// ResetError can be used to clear errors, making Next able to continue. Some errors cannot be
	// cleared and Error will still return them after ResetError, you have to check this.
	ResetError()
}
