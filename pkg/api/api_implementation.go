package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"

	"github.com/go-logr/logr"

	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1helper "go.anx.io/go-anxcloud/pkg/apis/core/v1/helper"
	"go.anx.io/go-anxcloud/pkg/client"
)

const (
	// ListChannelDefaultPageSize specifies the default page size for List operations returning the data via channel.
	ListChannelDefaultPageSize = 10
)

// defaultAPI is the type for our generic implementation of the API interface.
type defaultAPI struct {
	client client.Client

	// If the user wants to override the logging for this API, it can be done by providing the [WithLogger] method.
	// Said method is responsible for persisting the correct logger on this field.
	//
	// However, for calling, this field should *never* be called. Instead, use the Logger(context.Context) method
	// on this struct to get a valid logger for further usage.
	logger logr.Logger

	clientOptions  []client.Option
	requestOptions []types.Option
}

// Logger returns the logger for the given API in the following order:
//  1. If set, the logger on the struct is returned.
//  2. If the struct does not have a logger, the logger from the context is returned.
//  3. In all other cases, the logger of the underlying client is returned, which falls back to the discard logger.
func (a defaultAPI) Logger(ctx context.Context) logr.Logger {
	if !a.logger.IsZero() {
		return a.logger
	}

	if ctxLogger, err := logr.FromContext(ctx); err == nil {
		return ctxLogger
	}

	// If the WithLogger method wasn't called, we use the logger of the underlying client.
	// Because all of this is a rather hacky workaround for SYSENG-1746, the interface is not strongly typed.
	//
	// FIXME: Remove the duplicate logger handling.
	c, ok := a.client.(interface{ Logger() logr.Logger })
	if !ok {
		panic("The underlying API client does not have a public Logger() method.")
	}

	return c.Logger()
}

// NewAPIOption is the type for giving options to the NewAPI function.
type NewAPIOption func(*defaultAPI)

// WithClientOptions configures the API to pass the given client.Option to the client when creating it.
func WithClientOptions(o ...client.Option) NewAPIOption {
	return func(a *defaultAPI) {
		a.clientOptions = append(a.clientOptions, o...)
	}
}

// WithRequestOptions configures default options applied to requests
func WithRequestOptions(opts ...types.Option) NewAPIOption {
	return func(a *defaultAPI) {
		a.requestOptions = opts
	}
}

// WithLogger configures the API to use the given logger. It is recommended to pass a named logger.
// If you don't pass an existing client, the logger you give here will be given to the client (with
// added name "client").
func WithLogger(l logr.Logger) NewAPIOption {
	return func(a *defaultAPI) {
		a.logger = l
		a.clientOptions = append(a.clientOptions, client.Logger(l.WithName("client")))
	}
}

// NewAPI creates a new API client which implements the API interface.
func NewAPI(opts ...NewAPIOption) (API, error) {
	api := defaultAPI{
		clientOptions: []client.Option{
			client.ParseEngineErrors(false),
		},
	}

	for _, opt := range opts {
		opt(&api)
	}

	if api.client == nil {
		if c, err := client.New(api.clientOptions...); err == nil {
			api.client = c
		} else {
			return nil, err
		}
	}

	return api, nil
}

// Get the identified object from the engine.
func (a defaultAPI) Get(ctx context.Context, o types.IdentifiedObject, opts ...types.GetOption) error {
	options := types.GetOptions{}
	var err error
	for _, opt := range resolveRequestOptions(a.requestOptions, opts) {
		err = errors.Join(err, opt.ApplyToGet(&options))
	}
	if err != nil {
		return fmt.Errorf("apply request options: %w", err)
	}

	return a.do(ctx, o, o, &options, types.OperationGet)
}

// Create the given object on the engine.
func (a defaultAPI) Create(ctx context.Context, o types.Object, opts ...types.CreateOption) error {
	options := types.CreateOptions{}
	var err error
	for _, opt := range resolveRequestOptions(a.requestOptions, opts) {
		err = errors.Join(err, opt.ApplyToCreate(&options))
	}
	if err != nil {
		return fmt.Errorf("apply request options: %w", err)
	}

	if err := a.do(ctx, o, o, &options, types.OperationCreate); err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	return a.handlePostCreateOptions(ctx, o, options)
}

// handlePostCreateOptions executes configured Create options
// which should be handled after the object was successfully created
func (a defaultAPI) handlePostCreateOptions(ctx context.Context, o types.IdentifiedObject, options types.CreateOptions) error {
	if options.AutoTags != nil {
		if err := corev1helper.TaggerImplementation.Tag(ctx, a, o, options.AutoTags...); err != nil {
			return newErrTaggingFailed(err)
		}
	}

	return nil
}

// Update the object on the engine.
func (a defaultAPI) Update(ctx context.Context, o types.IdentifiedObject, opts ...types.UpdateOption) error {
	options := types.UpdateOptions{}
	var err error
	for _, opt := range resolveRequestOptions(a.requestOptions, opts) {
		err = errors.Join(err, opt.ApplyToUpdate(&options))
	}
	if err != nil {
		return fmt.Errorf("apply request options: %w", err)
	}

	return a.do(ctx, o, o, &options, types.OperationUpdate)
}

// Destroy the identified object.
func (a defaultAPI) Destroy(ctx context.Context, o types.IdentifiedObject, opts ...types.DestroyOption) error {
	options := types.DestroyOptions{}
	var err error
	for _, opt := range resolveRequestOptions(a.requestOptions, opts) {
		err = errors.Join(err, opt.ApplyToDestroy(&options))
	}
	if err != nil {
		return fmt.Errorf("apply request options: %w", err)
	}

	return a.do(ctx, o, o, &options, types.OperationDestroy)
}

// List objects matching the info given in the object.
func (a defaultAPI) List(ctx context.Context, o types.FilterObject, opts ...types.ListOption) error {
	options := types.ListOptions{}
	var err error
	for _, opt := range resolveRequestOptions(a.requestOptions, opts) {
		err = errors.Join(err, opt.ApplyToList(&options))
	}
	if err != nil {
		return fmt.Errorf("apply request options: %w", err)
	}

	ctx, err = a.contextPrepare(ctx, o, types.OperationList, &options)

	if err != nil {
		return err
	}

	req, err := a.makeRequest(ctx, o, nil, types.OperationList)
	if err != nil {
		return err
	}
	ctx = req.Context() // makeRequest extends the context

	var channelPageIterator types.PageInfo
	if options.ObjectChannel != nil && !options.Paged {
		options.Paged = true
		options.Page = 1
		options.EntriesPerPage = ListChannelDefaultPageSize
		options.PageInfo = &channelPageIterator
	} else if options.ObjectChannel != nil && options.PageInfo != nil {
		return ErrCannotListChannelAndPaged
	}

	singlePageMode := false

	if psh, ok := o.(types.PaginationSupportHook); ok {
		v, err := psh.HasPagination(ctx)
		if err != nil {
			return err
		}
		singlePageMode = !v
	}

	if options.Paged {
		if options.Page == 0 {
			a.Logger(ctx).V(1).Info("List called requesting page 0, fixing to page 1")
			options.Page = 1
		}

		if !singlePageMode {
			addPaginationQueryParameters(req, options)
		}
	}

	result := json.RawMessage{}
	err = a.doRequest(req, o, &result)
	if err != nil {
		return err
	}

	if options.Paged {
		fetcher := func(page uint) (json.RawMessage, error) {
			req := req.Clone(ctx)

			if !singlePageMode {
				query := req.URL.Query()
				query.Set("page", strconv.FormatUint(uint64(page), 10))

				req.URL.RawQuery = query.Encode()
			}

			var response json.RawMessage
			err := a.doRequest(req, o, &response)
			if err != nil {
				return nil, err
			}

			return response, nil
		}

		iter, err := newPageIter(ctx, a, result, options, fetcher, singlePageMode)
		if err != nil {
			return err
		}

		*options.PageInfo = iter
	}

	if options.ObjectChannel != nil {
		c := make(chan types.ObjectRetriever)
		*options.ObjectChannel = c

		objectRetrieved := make(chan bool)
		go func(pi types.PageInfo) {
			var pageData []json.RawMessage

		outer:
			for pi.Next(&pageData) {
				for _, o := range pageData {
					// since we are in a goroutine, we might already be in the next iteration of this loop
					// at the time the receiving end of this channel calls the closure. Having a loop-body
					// scoped variables makes the data for the closure perfectly identified.
					closureData := o
					c <- func(out types.Object) error {
						err := decodeResponse(ctx, "application/json", bytes.NewBuffer(closureData), out)
						if err != nil {
							return err
						}

						if options.FullObjects {
							if err := a.Get(ctx, out); err != nil {
								return err
							}
						}

						select {
						case <-ctx.Done():
						case objectRetrieved <- true:
						}

						return nil
					}

					select {
					case <-ctx.Done():
						break outer
					case <-objectRetrieved:
					}
				}
			}

			close(c)
		}(channelPageIterator)
	}

	return nil
}

func (a defaultAPI) makeRequest(ctx context.Context, obj types.Object, body interface{}, op types.Operation) (*http.Request, error) {
	singleObjectOperation := op == types.OperationGet || op == types.OperationUpdate || op == types.OperationDestroy

	// We do this right on top to checks if the Object has a suitible identifier set (when required by the operation type).
	identifier, err := obj.GetIdentifier(ctx)
	if err != nil {
		return nil, err
	} else if singleObjectOperation && identifier == "" {
		return nil, ErrUnidentifiedObject
	}

	resourceURL, err := obj.EndpointURL(ctx)
	if err != nil {
		return nil, err
	}

	ctx = types.ContextWithURL(ctx, *resourceURL)

	baseURL, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("error parsing client's BaseURL: %w", err)
	}

	fullResourcePath := path.Join(baseURL.Path, resourceURL.Path)

	if singleObjectOperation {
		fullResourcePath = path.Join(fullResourcePath, identifier)
	}

	fullQuery := baseURL.Query()
	resourceQuery := resourceURL.Query()
	for key, vals := range resourceQuery {
		for _, val := range vals {
			fullQuery.Add(key, val)
		}
	}

	fullURL := url.URL{
		Scheme: baseURL.Scheme,
		// Opaque URLs are not supported by us
		User:     baseURL.User,
		Host:     baseURL.Host,
		Path:     fullResourcePath,
		RawQuery: fullQuery.Encode(),
		// Fragment is never sent to a server
	}

	if obj, ok := obj.(types.FilterRequestURLHook); ok {
		filteredURL, err := obj.FilterRequestURL(ctx, &fullURL)
		if err != nil {
			return nil, err
		}
		fullURL = *filteredURL
	}

	var method string
	hasRequestBody := false

	switch op {
	case types.OperationGet:
		fallthrough
	case types.OperationList:
		method = "GET"

	case types.OperationCreate:
		method = "POST"
		hasRequestBody = true
	case types.OperationUpdate:
		method = "PUT"
		hasRequestBody = true
	case types.OperationDestroy:
		method = "DELETE"
	default:
		return nil, ErrOperationNotSupported
	}

	var bodyReader io.Reader = nil

	if hasRequestBody {
		buffer := bytes.Buffer{}

		var requestBody interface{} = body

		if filterRequestBody, ok := obj.(types.RequestBodyHook); ok {
			rb, err := filterRequestBody.FilterAPIRequestBody(ctx)

			if err != nil {
				return nil, err
			}

			requestBody = rb
		}

		if err := json.NewEncoder(&buffer).Encode(requestBody); err != nil {
			return nil, err
		}

		bodyReader = &buffer
	}

	request, err := http.NewRequestWithContext(ctx, method, fullURL.String(), bodyReader)
	if err != nil {
		// currently unreachable. http.NewRequestWithContext() returns an error in the following cases:
		// * the passed method is invalid (we have a hardcoded list of methods we use some lines above)
		// * ctx is nil (we check that in prepareContext() already)
		// * the URL cannot be parsed (we check that already some lines above)
		// makes it non-testable right now, but, I don't care since it's because of all errors already handled.
		// -- Mara @LittleFox94 Grosch, 2021-10-27
		return nil, err
	}

	if hasRequestBody {
		request.Header.Add("Content-Type", "application/json; charset=utf-8")
	}

	if filterRequest, ok := obj.(types.RequestFilterHook); ok {
		request, err = filterRequest.FilterAPIRequest(ctx, request)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

func (a defaultAPI) doRequest(req *http.Request, obj types.Object, body interface{}) error {
	ctx := req.Context()
	response, err := a.client.Do(req)

	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	if filterResponse, ok := obj.(types.ResponseFilterHook); ok {
		response, err = filterResponse.FilterAPIResponse(ctx, response)
	}

	if err != nil {
		return fmt.Errorf("Object returned an error from FilterAPIResponse: %w", err)
	}

	defer response.Body.Close()

	if err := ErrorFromResponse(req, response); err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		mediaType, err := getResponseType(response)
		if err != nil {
			return err
		}

		return decodeResponse(ctx, mediaType, response.Body, body)
	}

	return nil
}

// contextPrepare attaches the given operation op, the options opts and the logger of the API to a newly constructed context and returns that.
// If ctx is nil, [ErrContextRequired] is returned.
func (a defaultAPI) contextPrepare(ctx context.Context, o types.Object, op types.Operation, opts types.Options) (context.Context, error) {
	if ctx == nil {
		return nil, ErrContextRequired
	}

	ctx = types.ContextWithOperation(ctx, op)
	ctx = types.ContextWithOptions(ctx, opts)

	objectType := reflect.TypeOf(o)
	for objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}

	logger := a.Logger(ctx)
	return logr.NewContext(ctx, logger.WithValues("operation", op, "resource", objectType)), nil
}

func (a defaultAPI) do(ctx context.Context, obj types.Object, body interface{}, opts types.Options, op types.Operation) error {
	var err error
	ctx, err = a.contextPrepare(ctx, obj, op, opts)

	if err != nil {
		return err
	}

	request, err := a.makeRequest(ctx, obj, body, op)
	if err != nil {
		return err
	}

	return a.doRequest(request, obj, body)
}

func getResponseType(res *http.Response) (string, error) {
	knownTypes := []string{"application/json"}

	if contentType := res.Header.Get("content-type"); contentType != "" {
		mt, _, err := mime.ParseMediaType(contentType)

		if err != nil {
			return "", fmt.Errorf("error parsing Content-Type header in Engine response: %w (was: '%v')", err, contentType)
		}

		for _, kt := range knownTypes {
			if kt == mt {
				return mt, nil
			}
		}

		return "", fmt.Errorf("%w: unknown mime-type %v", ErrUnsupportedResponseFormat, mt)
	}

	return "application/json", nil
}

func addPaginationQueryParameters(req *http.Request, opts types.ListOptions) {
	query := req.URL.Query()
	query.Add("page", strconv.FormatUint(uint64(opts.Page), 10))
	query.Add("limit", strconv.FormatUint(uint64(opts.EntriesPerPage), 10))

	req.URL.RawQuery = query.Encode()
}

func decodeResponse(ctx context.Context, mediaType string, data io.Reader, res interface{}) error {
	if mediaType == "application/json" {
		if decodeResponse, ok := res.(types.ResponseDecodeHook); ok {
			if err := decodeResponse.DecodeAPIResponse(ctx, data); err != nil {
				return err
			}

			return nil
		}

		return json.NewDecoder(data).Decode(res)
	}

	return fmt.Errorf("%w: no idea how to handle media type %v", ErrUnsupportedResponseFormat, mediaType)
}

func resolveRequestOptions[T any](commonOptions []types.Option, requestOptions []T) []T {
	return append(filterOptions[T](commonOptions), requestOptions...)
}

func filterOptions[T any](opts []types.Option) []T {
	ret := make([]T, 0, len(opts))
	for _, v := range opts {
		if v, ok := v.(T); ok {
			ret = append(ret, v)
		}
	}

	return ret
}
