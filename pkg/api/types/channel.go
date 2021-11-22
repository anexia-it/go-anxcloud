package types

// ObjectRetriever retrieves an object and decodes it to the correct go type.
type ObjectRetriever func(Object) error

// ObjectChannel streams objects.
type ObjectChannel <-chan ObjectRetriever

// ObjectChannelCloser is used to List objects via a channel that can be closed by the retriever, canceling
// the list operation.
// The user is required to call Close() when not reading the channel until it finishes, it will leak goroutines
// otherwise.
type ObjectChannelCloser interface {
	Channel() (ObjectChannel, error)
	Close() error
}
