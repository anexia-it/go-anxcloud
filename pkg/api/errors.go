package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	// ErrOperationNotSupported is returned when requesting an operation on a resource it does not support.
	ErrOperationNotSupported = errors.New("requested operation is not supported by the resource type")

	// ErrUnsupportedResponseFormat is set when the engine responds in a format we don't understand, for example unknown Content-Types.
	ErrUnsupportedResponseFormat = errors.New("response format is not supported")

	// ErrPageResponseNotSupported is returned when trying to parse a paged response and the format of the response body is not (yet) supported.
	ErrPageResponseNotSupported = fmt.Errorf("paged response invalid: %w", ErrUnsupportedResponseFormat)

	// ErrCannotListChannelAndPaged is returned when the user List()ing with the AsObjectChannel() and Paged() options set and didn't gave nil for Paged() PageInfo output argument.
	ErrCannotListChannelAndPaged = errors.New("list with Paged and ObjectChannel is only valid when not retrieving the PageInfo iterator via Paged option")

	// ErrContextRequired is returned when a nil context was passed as argument.
	ErrContextRequired = errors.New("no context given")
)

// RateLimitError occurs after a [http.TooManyRequests] status code got returned by the engine.
type RateLimitError struct {
	// RetryAfter contains the point in time at which the request can be retried again.
	//
	// Note: Right now, it's always set to be 30 minutes into the future, because we do not have a Retry-After header yet.
	// This is tracked internally in ENGSUP-9027, after which this value will contain the correct value.
	RetryAfter time.Time
}

func (e RateLimitError) Error() string {
	d := time.Until(e.RetryAfter)
	return fmt.Sprintf("rate limited by the engine, retry in %s", d)
}

// IsRateLimitError returns true if err is of type [RateLimitError]
// or of type [HTTPError] with the status code of [http.StatusTooManyRequests].
//
// This is provided as a convenience helper, because errors.Is against the
// RateLimitError would compare the RetryAfter value of the error as well.
func IsRateLimitError(err error) bool {
	switch e := err.(type) {
	case RateLimitError:
		return true
	case HTTPError:
		return e.statusCode == http.StatusTooManyRequests
	default:
		return false
	}
}

// EngineError is the base type for all errors returned by the engine.
//
// Ideally all errors returned by the API are transformed into EngineErrors, making HTTPError obsolete, as this
// would completely decouple communicating with the Engine from using HTTP.
type EngineError struct {
	message string
	wrapped error
}

var (
	// ErrNotFound is returned when the given identified object does not exist in the engine. Take a look at IgnoreNotFound(), too.
	ErrNotFound EngineError = EngineError{message: "requested resource does not exist on the engine"}

	// ErrAccessDenied is returned when the used authentication credential is not authorized to do the requested operation.
	ErrAccessDenied EngineError = EngineError{message: "access to requested resource was denied by the engine"}
)

// IgnoreNotFound is a helper to handle ErrNotFound differently than other errors with less code.
func IgnoreNotFound(err error) error {
	if errors.Is(err, ErrNotFound) {
		return nil
	}

	return err
}

// Error returns the message of the EngineError, implementing the `error` interface.
func (e EngineError) Error() string {
	return e.message
}

// Unwrap returns the wrapped error of the EngineError, making it compatible with `errors.Is/As/Unwrap`.
func (e EngineError) Unwrap() error {
	return e.wrapped
}

// HTTPError is an not-specially-implemented EngineError for a given status code. Ideally this is not used
// because every returned error is mapped to an ErrSomething package variable, decoupling error handling from
// the transport protocol.
type HTTPError struct {
	message    string
	wrapped    error
	statusCode int
	url        *url.URL
	method     string
}

// newHTTPError creates a new HTTP error, taking the information from the given request and response. It
// can optionally wrap an error and have a custom message.
func newHTTPError(req *http.Request, res *http.Response, wrapped error, message *string) HTTPError {
	var msg string

	if message != nil {
		msg = *message
	} else {
		msg = fmt.Sprintf("Engine returned an error: %v (%v)", res.Status, res.StatusCode)
	}

	e := HTTPError{
		message:    msg,
		wrapped:    wrapped,
		statusCode: res.StatusCode,
		url:        req.URL,
		method:     req.Method,
	}

	return e
}

// StatusCode returns the HTTP status code of the HTTPError.
func (e HTTPError) StatusCode() int {
	return e.statusCode
}

// Unwrap returns the error which caused this one.
func (e HTTPError) Unwrap() error {
	return e.wrapped
}

// Error returns the error message.
func (e HTTPError) Error() string {
	return e.message
}

// ErrorFromResponse creates a new HTTPError from the given response.
func ErrorFromResponse(req *http.Request, res *http.Response) error {
	var specificError error

	switch res.StatusCode {
	case 403:
		specificError = ErrAccessDenied
	case 404:
		specificError = ErrNotFound
	case 429:
		// TODO: Set this to a concrete value extracted from the Retry-After header.
		specificError = RateLimitError{RetryAfter: time.Now().Add(time.Minute * 30)}
	}

	// We check for higher than 300 because redirects should be handled already
	if res.StatusCode > 300 || specificError != nil {
		return newHTTPError(req, res, specificError, nil)
	}

	return nil
}

// NewHTTPError creates a new HTTPError instance with the given values, which is mostly useful for mock-testing.
func NewHTTPError(status int, method string, url *url.URL, wrapped error) error {
	return HTTPError{
		message:    http.StatusText(status),
		wrapped:    wrapped,
		statusCode: status,
		url:        url,
		method:     method,
	}
}

// ErrTaggingFailed is returned when resource tagging failed
type ErrTaggingFailed struct {
	wrapped error
}

// Unwrap returns the error which caused this one.
func (e ErrTaggingFailed) Unwrap() error {
	return e.wrapped
}

// Error returns the error message.
func (e ErrTaggingFailed) Error() string {
	return fmt.Sprintf("failed tagging resource: %s", e.wrapped)
}

func newErrTaggingFailed(wrapped error) error {
	return &ErrTaggingFailed{wrapped}
}
