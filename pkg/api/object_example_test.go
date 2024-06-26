package api

import (
	"context"
	"encoding/json"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"

	// the following are for testing only, you don't need them
	"fmt"
	"net/http"
	"net/http/httptest"

	"go.anx.io/go-anxcloud/pkg/client"
)

// ExampleObject is an API Object we define as example how to make something an Object.
//
// Objects must have tags for json encoding/decoding and exactly one must be tagged as anxcloud:"identifier".
type ExampleObject struct {
	Identifier          string `json:"identifier" anxcloud:"identifier"`
	SecondaryIdentifier string `json:"secondary_identifier"`
	Name                string `json:"name"`
}

// This is the most-basic implementation for EndpointURL, only returning the URL. This is the case for resources
// that support all operations and have the default URL mapping:
// * Create:  POST    url
// * List:    GET     url
// * Get:     GET     url/identifier
// * Update:  PUT     url/identifier
// * Destroy: DELETE  url/identifier
//
// Some objects don't support all operations, you'd then have to check the passed `op`eration and return
// ErrOperationNotSupported for unsupported operations.
//
// Sometimes URLs for Objects done match this schema. As long as the last part of the URL is the identifier for
// operations on specific objects, you can switch-case on the operation and return the correct URLs. The
// identifier is appended by default for Get, Update and Destroy operations. You can implement the interface
// types.RequestFilterHook to have full control over the requests done for your object.
func (o *ExampleObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get operation from context: %w", err)
	}

	endpointSuffix := "foo.json"
	if op == types.OperationGet && o.Identifier == "" && o.SecondaryIdentifier != "" {
		endpointSuffix = "foo/by-secondary-identifier.json"
	}

	return url.Parse(fmt.Sprintf("/example/v1/%s", endpointSuffix))
}

// GetIdentifier returns the identifier of an object.
// The Makefiles `generate` target creates this method for objects that have a anxcloud:"identifier" tagged field.
// Objects with multiple identifying fields can have this method overwritten to return an operation dependent identifier.
// The provided context might not always have the operation set. In this case the method must return the primary identifier.
//
// For example: `corev1.Location` has the (primary) `Identifier` attribute, which is required by many API endpoints.
// Additionally locations are identified by the human readable `Code` attribute (e.g. ANX04).
//
// Handling of secondary identifiers usually requires additional logic in the objects EndpointURL method.
func (o *ExampleObject) GetIdentifier(ctx context.Context) (string, error) {
	op, err := types.OperationFromContext(ctx)
	if o.Identifier != "" || err != nil {
		return o.Identifier, nil
	}

	if op == types.OperationGet {
		return o.SecondaryIdentifier, nil
	}

	return "", nil
}

// This is a more complex example, supporting to List with a filter
type ExampleFilterableObject struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
	Mode       string `json:"mode"`
}

// GetIdentifier returns the objects identifier
func (o *ExampleFilterableObject) GetIdentifier(context.Context) (string, error) {
	return o.Identifier, nil
}

// This is an example for the EndpointURL method for an Object that can use a filter for List operations.
//
// The API in this case expects a query argument called `filter` with a URL-encoded query string in it,
// so for filtering for name=foo and mode=tcp the full URL might look like this:
// `/filter_example/v1?filter=name%3Dfoo%26mode%3Dtcp`.
func (o *ExampleFilterableObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	// we can ignore the error since the URL is hard-coded known as valid
	u, _ := url.Parse("/filter_example/v1")

	if op, err := types.OperationFromContext(ctx); err == nil && op == types.OperationList {
		filter := url.Values{}

		if o.Name != "" {
			filter.Add("name", o.Name)
		}

		if o.Mode != "" {
			filter.Add("mode", o.Mode)
		}

		if filters := filter.Encode(); filters != "" {
			query := u.Query()
			query.Add("filter", filters)
			u.RawQuery = query.Encode()
		}
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// code below is not part of this example but makes it appear in the docs at all and uses it as test. //
////////////////////////////////////////////////////////////////////////////////////////////////////////

type ExampleObjectMockHandler struct {
	filtered []ExampleFilterableObject
}

func (h *ExampleObjectMockHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json; charset=utf-8")

	switch req.URL.Path {
	case "/example/v1/foo.json":
		o := ExampleObject{}
		_ = json.NewDecoder(req.Body).Decode(&o)

		o.Identifier = "some random identifier"
		_ = json.NewEncoder(res).Encode(o)
	case "/filter_example/v1":
		if page := req.URL.Query().Get("page"); page != "1" && page != "" {
			break
		}

		nameFilter := ""
		modeFilter := ""

		if f := req.URL.Query().Get("filter"); f != "" {
			filters, _ := url.ParseQuery(f)

			nameFilter = filters.Get("name")
			modeFilter = filters.Get("mode")
		}

		ret := make([]ExampleFilterableObject, 0, len(h.filtered))

		for _, o := range h.filtered {
			ok := true

			if nameFilter != "" && o.Name != nameFilter {
				ok = false
			}

			if modeFilter != "" && o.Mode != modeFilter {
				ok = false
			}

			if ok {
				ret = append(ret, o)
			}
		}

		_ = json.NewEncoder(res).Encode(ret)
	}
}

func Example_implementObject() {
	mock := ExampleObjectMockHandler{
		filtered: []ExampleFilterableObject{
			{Name: "hello TCP 1", Mode: "tcp", Identifier: "random identifier 1"},
			{Name: "hello UDP 1", Mode: "udp", Identifier: "random identifier 2"},
			{Name: "hello TCP 2", Mode: "tcp", Identifier: "random identifier 3"},
			{Name: "hello UDP 2", Mode: "udp", Identifier: "random identifier 4"},
		},
	}

	server := httptest.NewServer(&mock)

	api, err := NewAPI(
		WithClientOptions(
			client.IgnoreMissingToken(),
			client.BaseURL(server.URL),
		),
	)

	if err != nil {
		fmt.Printf("Error creating API instance: %v\n", err)
		return
	}

	ctx := context.TODO()

	// trying to create an ExampleObject on the API
	o := ExampleObject{Name: "hello world"}
	if err := api.Create(ctx, &o); err != nil {
		fmt.Printf("Error creating object on API: %v\n", err)
		return
	}

	fmt.Printf("Object created, identifier '%v'\n", o.Identifier)

	// trying to list ExampleFilterableObjects on the API, filtered on mode=tcp
	fo := ExampleFilterableObject{Mode: "tcp"}
	var fopi types.PageInfo
	if err := api.List(ctx, &fo, Paged(1, 1, &fopi)); err != nil {
		fmt.Printf("Error listing objects on API: %v\n", err)
		return
	}

	var fos []ExampleFilterableObject
	for fopi.Next(&fos) {
		for _, fo := range fos {
			fmt.Printf("Retrieved object with mode '%v' named '%v'\n", fo.Mode, fo.Name)
		}
	}

	// Output:
	// Object created, identifier 'some random identifier'
	// Retrieved object with mode 'tcp' named 'hello TCP 1'
	// Retrieved object with mode 'tcp' named 'hello TCP 2'
}

func ExampleWithRequestOptions() {
	api, err := NewAPI(
		WithRequestOptions(
			// automatically assign tags to newly created resources
			AutoTag("foo", "bar"),
		),
	)

	if err != nil {
		panic(fmt.Errorf("Error creating API instance: %v\n", err))
	}

	// create resource and automatically apply 'foo' & 'bar' tags
	if err := api.Create(context.TODO(), &ExampleObject{Name: "foo"}); err != nil {
		panic(err)
	}
}
