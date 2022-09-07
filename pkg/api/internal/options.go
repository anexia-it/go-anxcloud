package internal

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
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
	Channel *types.ObjectChannel
}

// ApplyToList applies the AsObjectChannel option to all the ListOptions.
func (aoc ObjectChannelOption) ApplyToList(o *types.ListOptions) {
	o.ObjectChannel = aoc.Channel
}

// FullObjectsOption configures if the List operation shall make a Get operation for each object before
// returning it to the caller.
type FullObjectsOption bool

// ApplyToList applies the FullObjectsOption option to all the ListOptions.
func (foo FullObjectsOption) ApplyToList(o *types.ListOptions) {
	o.FullObjects = bool(foo)
}

// AutoTagOption configures the Create operation to automatically tag objects after creation
type AutoTagOption []string

// ApplyToCreate applies the AutoTagOption to the ListOptions
func (ato AutoTagOption) ApplyToCreate(o *types.CreateOptions) {
	o.AutoTags = ato
}
