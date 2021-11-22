package internal

import (
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

// PagedOption is an option valid for List operations to retrieve objects in a paged fashion (instead of all at once).
type PagedOption struct {
	// Page to retrieve
	Page uint

	// Entries per page
	Limit uint

	// Additional output about the current page, includes a way to iterate through all pages.
	Info *types.PageInfo
}

// ApplyToList applies the Paged option to all the ListOptions.
func (p PagedOption) ApplyToList(o *types.ListOptions) {
	o.Paged = true
	o.Page = p.Page
	o.EntriesPerPage = p.Limit
	o.PageInfo = p.Info
}

// ObjectChannelOption configures the List operation to return the objects via the given channel.
type ObjectChannelOption struct {
	Channel types.ObjectChannelCloser
}

// ApplyToList applies the ObjectChannelOption to all the ListOptions.
func (oc ObjectChannelOption) ApplyToList(o *types.ListOptions) {
	o.Channel = oc.Channel
}
