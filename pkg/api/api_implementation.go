package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
)

const (
	// ListChannelDefaultPageSize specifies the default page size for List operations returning the data via channel.
	ListChannelDefaultPageSize = 10
)

// defaultAPI is the type for our generic implementation of the API interface.
type defaultAPI struct {
	client client.Client
	logger *logr.Logger

	clientOptions []client.Option
}

// NewAPIOption is the type for giving options to the NewAPI function.
type NewAPIOption func(*defaultAPI)

// WithClientOptions configures the API to pass the given client.Option to the client when creating it.
func WithClientOptions(o ...client.Option) NewAPIOption {
	return func(a *defaultAPI) {
		a.clientOptions = append(a.clientOptions, o...)
	}
}

// WithLogger configures the API to use the given logger. It is recommended to pass a named logger.
// If you don't pass an existing client, the logger you give here will given to the client (with
// added name "client").
func WithLogger(l logr.Logger) NewAPIOption {
	return func(a *defaultAPI) {
		a.logger = &l
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
	for _, opt := range opts {
		opt.ApplyToGet(&options)
	}

	return a.do(ctx, o, o, &options, types.OperationGet)
}

// Create the given object on the engine.
func (a defaultAPI) Create(ctx context.Context, o types.Object, opts ...types.CreateOption) error {
	options := types.CreateOptions{}
	for _, opt := range opts {
		opt.ApplyToCreate(&options)
	}

	return a.do(ctx, o, o, &options, types.OperationCreate)
}

// Update the object on the engine.
func (a defaultAPI) Update(ctx context.Context, o types.IdentifiedObject, opts ...types.UpdateOption) error {
	options := types.UpdateOptions{}
	for _, opt := range opts {
		opt.ApplyToUpdate(&options)
	}

	return a.do(ctx, o, o, &options, types.OperationUpdate)
}

// Destroy the identified object.
func (a defaultAPI) Destroy(ctx context.Context, o types.IdentifiedObject, opts ...types.DestroyOption) error {
	options := types.DestroyOptions{}
	for _, opt := range opts {
		opt.ApplyToDestroy(&options)
	}

	return a.do(ctx, o, o, &options, types.OperationDestroy)
}

// List objects matching the info given in the object.
func (a defaultAPI) List(ctx context.Context, o types.FilterObject, opts ...types.ListOption) error {
	options := types.ListOptions{}
	for _, opt := range opts {
		opt.ApplyToList(&options)
	}

	var err error
	ctx, err = a.contextPrepare(ctx, o, types.OperationList)

	if err != nil {
		return err
	}

	req, err := a.makeRequest(ctx, o, nil, &options, types.OperationList)
	if err != nil {
		return err
	}

	var channelPageIterator types.PageInfo
	if options.ObjectChannel != nil && !options.Paged {
		options.Paged = true
		options.Page = 1
		options.EntriesPerPage = ListChannelDefaultPageSize
		options.PageInfo = &channelPageIterator
	} else if options.ObjectChannel != nil && options.PageInfo != nil {
		return ErrCannotListChannelAndPaged
	}

	if options.Paged {
		addPaginationQueryParameters(ctx, req, options)
	}

	result := json.RawMessage{}
	err = a.doRequest(req, o, &result, &options, types.OperationList)
	if err != nil {
		return err
	}

	if options.Paged {
		fetcher := func(page uint) (json.RawMessage, error) {
			req := req.Clone(ctx)

			query := req.URL.Query()
			query.Set("page", strconv.FormatUint(uint64(page), 10))

			req.URL.RawQuery = query.Encode()

			var response json.RawMessage
			err := a.doRequest(req, o, &response, &options, types.OperationList)
			if err != nil {
				return nil, err
			}

			return response, nil
		}

		iter, err := newPageIter(ctx, result, options, fetcher)
		if err != nil {
			return err
		}

		*options.PageInfo = iter
	}

	if options.ObjectChannel != nil {
		go func(c types.ObjectChannel, pi types.PageInfo) {
			var pageData []json.RawMessage

			for pi.Next(&pageData) {
				for _, o := range pageData {
					c <- func(out types.Object) error {
						return json.Unmarshal(o, out)
					}
				}
			}

			close(c)
		}(*options.ObjectChannel, channelPageIterator)
	}

	return nil
}

func (a defaultAPI) makeRequest(ctx context.Context, obj types.Object, body interface{}, opts types.Options, op types.Operation) (*http.Request, error) {
	singleObjectOperation := op == types.OperationGet || op == types.OperationUpdate || op == types.OperationDestroy

	// We do this right on top because this checks if the Object has a correct type which is more strictly defined than just the interface.
	// In a perfect world this would be a compile-time check.
	identifier, err := getObjectIdentifier(obj, singleObjectOperation)
	if err != nil {
		return nil, err
	}

	resourceURL, err := obj.EndpointURL(ctx, op, opts)
	if err != nil {
		return nil, err
	}

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
			rb, err := filterRequestBody.FilterAPIRequestBody(op, opts)

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
		request, err = filterRequest.FilterAPIRequest(op, opts, request)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

func (a defaultAPI) doRequest(req *http.Request, obj types.Object, body interface{}, opts types.Options, op types.Operation) error {
	response, err := a.client.Do(req)

	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	if filterResponse, ok := obj.(types.ResponseFilterHook); ok {
		response, err = filterResponse.FilterAPIResponse(op, opts, response)
	}

	if err != nil {
		return fmt.Errorf("Object returned an error from FilterAPIResponse: %w", err)
	}

	defer response.Body.Close()

	if err := errorFromResponse(req, response); err != nil {
		return err
	}

	if mediaType, err := getResponseType(response); err == nil {
		if mediaType == "application/json" {
			return json.NewDecoder(response.Body).Decode(body)
		}

		// unreachable, getResponseType() already checks for supported types
		return ErrUnsupportedResponseFormat
	} else {
		return err
	}
}

func (a defaultAPI) contextPrepare(ctx context.Context, o types.Object, op types.Operation) (context.Context, error) {
	if ctx == nil {
		return nil, ErrContextRequired
	}

	objectType := reflect.TypeOf(o)
	for objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}

	logger := logr.Discard()

	// Checking if we have a logger on the context and attach one if we don't.
	if l, err := logr.FromContext(ctx); err != nil && a.logger != nil {
		logger = *a.logger
	} else if err == nil {
		// TODO(LittleFox94): derive a named one from this?
		logger = l
	}

	return logr.NewContext(ctx, logger.WithValues("operation", op, "resource", objectType)), nil
}

func (a defaultAPI) do(ctx context.Context, obj types.Object, body interface{}, opts types.Options, op types.Operation) error {
	var err error
	ctx, err = a.contextPrepare(ctx, obj, types.OperationList)

	if err != nil {
		return err
	}

	request, err := a.makeRequest(ctx, obj, body, opts, op)
	if err != nil {
		return err
	}

	return a.doRequest(request, obj, body, opts, op)
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

func addPaginationQueryParameters(ctx context.Context, req *http.Request, opts types.ListOptions) {
	if opts.Page == 0 {
		log := logr.FromContextOrDiscard(ctx)
		log.V(1).Info("List called requesting page 0, fixing to page 1")

		opts.Page = 1
	}

	query := req.URL.Query()
	query.Add("page", strconv.FormatUint(uint64(opts.Page), 10))
	query.Add("limit", strconv.FormatUint(uint64(opts.EntriesPerPage), 10))

	req.URL.RawQuery = query.Encode()
}