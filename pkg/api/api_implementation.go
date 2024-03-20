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
	"slices"
	"strconv"
	"strings"

	"github.com/go-logr/logr"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
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
			log := logr.FromContextOrDiscard(ctx)
			log.V(1).Info("List called requesting page 0, fixing to page 1")

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

		requestBody, err = flattenObject(ctx, requestBody)
		if err != nil {
			return nil, err
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

	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get operation from request context: %w", err)
	}
	if isGenericServiceResource(obj) && op == types.OperationDestroy {
		response.Body.Close()
		response.Body = io.NopCloser(bytes.NewReader([]byte("{}")))
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

func isGenericServiceResource(o any) bool {
	ot := reflect.TypeOf(o)
	for ; ot.Kind() == reflect.Ptr; ot = ot.Elem() {
		// resolve pointer
	}

	if ot.Kind() != reflect.Struct {
		return false
	}

	gsBaseType := reflect.TypeOf(gs.GenericService{})
	for i := 0; i < ot.NumField(); i++ {
		if sf := ot.Field(i); sf.Type == gsBaseType && sf.Anonymous {
			return true
		}
	}

	return false
}

func flattenObject(ctx context.Context, o any) (any, error) {
	type IdentifiedObject interface {
		GetIdentifier(context.Context) (string, error)
	}

	type flatmeta struct {
		multi    bool
		ptr      bool
		embedVal reflect.Value
	}

	origType := reflect.TypeOf(o)
	origInstance := reflect.ValueOf(o)

	if origType.Kind() == reflect.Pointer {
		if origInstance.IsNil() {
			return o, nil
		}

		origType = origType.Elem()
		origInstance = origInstance.Elem()
	}

	if origType.Kind() != reflect.Struct {
		// not a pointer to a struct? -> can't flatten -> return as is
		return o, nil
	}

	identifiedObjectType := reflect.TypeOf((*IdentifiedObject)(nil)).Elem()

	var newFields []reflect.StructField
	var flatMetaMap = make(map[int]flatmeta)

	// iterate all struct fields to construct new type
	// if field is not exported
	// -> omit from new type
	// if field is anonymous
	// -> flatten recursive
	// if field contains `flatten` in anxencode tag
	// -> change field type to (pointer - if original type is also a pointer) string
	// -> store some meta in flatMetaMap (is pointer? is multi?)
	for i := 0; i < origType.NumField(); i++ {
		field := origType.Field(i)

		if !field.IsExported() {
			// not exported? -> omit in new type
			continue
		} else if field.Anonymous {
			// handle embedded field recursively
			flattened, err := flattenObject(ctx, origInstance.Field(i).Interface())
			if err != nil {
				return nil, err
			}

			field.Type = reflect.TypeOf(flattened)
			val := reflect.ValueOf(flattened)
			flatMetaMap[i] = flatmeta{embedVal: val}
		} else {
			isPtr := field.Type.Kind() == reflect.Pointer

			fieldType := field.Type
			if isPtr {
				fieldType = fieldType.Elem()
			}

			if encodeOpts, ok := field.Tag.Lookup("anxencode"); ok && slices.Contains(strings.Split(encodeOpts, ","), "flatten") {
				flatMetaMap[i] = flatmeta{ptr: isPtr, multi: fieldType.Kind() == reflect.Slice}
				if !fieldType.Implements(identifiedObjectType) &&
					!(fieldType.Kind() == reflect.Slice && fieldType.Elem().Implements(identifiedObjectType)) {
					return nil, fmt.Errorf("%q does not implement IdentifiedObject interface", field.Name)
				}

				newType := reflect.TypeOf("")
				if isPtr {
					newType = reflect.PointerTo(newType)
				}

				field.Type = newType
			}
		}

		newFields = append(newFields, field)
	}

	newType := reflect.StructOf(newFields)
	newInstance := reflect.New(newType).Elem()

	reflectGetIdentifier := func(v reflect.Value) (string, error) {
		ret := v.MethodByName("GetIdentifier").Call([]reflect.Value{reflect.ValueOf(ctx)})
		if err := ret[1].Interface(); err != nil {
			return "", err.(error)
		}

		return ret[0].String(), nil
	}

	skipped := 0
	// set values in instance of new type
	for oi := 0; oi < origInstance.NumField(); oi++ {
		ni := oi - skipped
		if !origType.Field(oi).IsExported() {
			skipped += 1
			continue
		}

		if meta, ok := flatMetaMap[oi]; ok {
			newField := newInstance.Field(ni)

			if newType.Field(ni).Anonymous {
				newField.Set(meta.embedVal)
				continue
			}

			// original field value was null? -> nothing to do
			if meta.ptr && origInstance.Field(oi).IsNil() {
				continue
			}

			origFieldVal := origInstance.Field(oi)
			if meta.ptr {
				origFieldVal = origFieldVal.Elem()
			}

			var newFieldValString string

			if !meta.multi {
				id, err := reflectGetIdentifier(origFieldVal)
				if err != nil {
					return nil, err
				}

				if meta.ptr && id == "" {
					return nil, types.ErrUnidentifiedObject
				}

				newFieldValString = id
			} else {
				// is multi? -> join ids with comma into a single string
				ids := make([]string, 0, origFieldVal.Len())
				for j := 0; j < origFieldVal.Len(); j++ {
					id, err := reflectGetIdentifier(origFieldVal.Index(j))
					if err != nil {
						return nil, err
					}

					if id == "" {
						return nil, types.ErrUnidentifiedObject
					}

					ids = append(ids, id)
				}

				newFieldValString = strings.Join(ids, ",")
			}

			var newVal reflect.Value
			if !meta.ptr {
				newVal = reflect.ValueOf(newFieldValString)
			} else {
				newVal = reflect.ValueOf(&newFieldValString)
			}

			newField.Set(newVal)
		} else {
			newInstance.Field(ni).Set(origInstance.Field(oi))
		}
	}

	return newInstance.Interface(), nil
}
