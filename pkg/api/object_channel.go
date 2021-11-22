package api

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/go-logr/logr"
)

const (
	// ListChannelDefaultPageSize specifies the default page size for List operations returning the data via channel.
	ListChannelDefaultPageSize = 10
)

type objectChannel struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	api API
	log logr.Logger

	objects      chan types.ObjectRetriever
	abort        chan bool
	pageIterator types.PageInfo
}

// NewObjectChannel creates a new object implementing types.ObjectChannelCloser which can be used after passed to API.List()
// with the ObjectChannel option.
func NewObjectChannel() types.ObjectChannelCloser {
	return &objectChannel{
		objects: make(chan types.ObjectRetriever),
		abort:   make(chan bool),
	}
}

// Channel starts the object listing logic and returns the channel objects can be retrieved from.
func (oc *objectChannel) Channel() (types.ObjectChannel, error) {
	if oc.pageIterator == nil {
		return nil, ErrObjectChannelNotReady
	}

	pageRetrieveWaiter := make(chan bool)
	go oc.run(pageRetrieveWaiter)
	pageRetrieveWaiter <- true

	return types.ObjectChannel(oc.objects), nil
}

// Close signals to abort listing objects and cancels any currently ongoing API requests.
func (oc *objectChannel) Close() error {
	if oc.pageIterator == nil {
		return ErrObjectChannelNotReady
	}

	oc.abort <- true
	oc.ctxCancel()
	return nil
}

func (oc *objectChannel) prepare(ctx context.Context, api API, opts *types.ListOptions) error {
	// configure the List operation to be paged with default settings if no other settings are configured
	// by the user already
	if !opts.Paged {
		opts.Paged = true
		opts.Page = 1
		opts.EntriesPerPage = ListChannelDefaultPageSize
	} else if opts.PageInfo != nil {
		// user tried to retrieve a page iterator and object channel - since we only create a single page
		// iterator used by the channel, returning this to the user wouldn't have the intended effects, so
		// we deny that.
		return ErrCannotListChannelAndPaged
	}

	oc.ctx, oc.ctxCancel = context.WithCancel(ctx)
	oc.api = api
	oc.log = logr.FromContextOrDiscard(ctx).WithName("ObjectChannel").V(2)

	opts.PageInfo = &oc.pageIterator

	return nil
}

func (oc *objectChannel) run(pageRetrieveWaiter chan bool) {
	var pageData []json.RawMessage

cancelChannel:
	for oc.pageIterator.Next(&pageData) {
		for _, o := range pageData {
			select {
			case <-oc.abort:
				oc.log.Info("User requested aborting the channel")
				break cancelChannel
			case <-pageRetrieveWaiter:
				oc.log.Info("Object retrieved by user")
			}

			oc.log.Info("Pushing retriever on channel")

			// since we are in a goroutine, we might already be in the next iteration of this loop
			// at the time the receiving end of this channel calls the closure. Having a loop-body
			// scoped variables makes the data for the closure perfectly identified.
			closureData := o
			oc.objects <- func(out types.Object) error {
				err := decodeResponse(oc.ctx, "application/json", bytes.NewBuffer(closureData), out)
				if err != nil {
					return err
				}

				pageRetrieveWaiter <- true
				return nil
			}
		}

		oc.log.Info("Retrieving next page")
	}

	// we have to wait for last retriever used before closing the channel
	<-pageRetrieveWaiter
	close(oc.objects)
}
