package types

import (
	"context"
	"net/http"
	"net/url"
)

// Object is the interface all objects to be retrieved by the generic API client are required to implement.
//
// On top of implementing this interface, an Object is always implemented as a struct, the pointer to it is what is passed to the generic API client.
// These Object structs need to have a member with `anxcloud:"identifier"` tag on it. Only strings and "github.com/satori/go.uuid".UUID identifiers
// are supported for now (look in pkg/api/object.go_getObjectIdentifier for code specifying the allowed types).
type Object interface {
	// Returns the URL to retrieve resources of the given type from or an error.
	// The request URL is formed of `client.BaseURL() + first return value of this function`, requests for a single object get
	// the object identifier appended to the path, a / added as needed. APIs using other URL schemes need to implement RequestFilterHook.
	EndpointURL(ctx context.Context, op Operation, options Options) (*url.URL, error)
}

// IdentifiedObject is the same as Object and is intended as a doc-helper. Objects are IdentifiedObjects when their
// identifying attribute (commonly something like "Identifier") is set.
type IdentifiedObject Object

// FilterObject is the same as Object and is intended as a doc-helper, telling the user only it's type and it implementing an interface is important for the generic API client.
// Some of the attributes on this Object can be used for filtering or searching, depending on the specific API.
type FilterObject Object

// RequestFilterHook is an interface Objects can optionally implement to modify requests before they are sent to the engine.
type RequestFilterHook interface {
	// FilterAPIRequest is called for every API request involving this object. Instead of the original request, the one returned from this function is sent to the engine.
	FilterAPIRequest(op Operation, options Options, req *http.Request) (*http.Request, error)
}

// RequestBodyHook is an interface Objects can optionally implement to customize request bodies based on the object given by the user of this library.
type RequestBodyHook interface {
	// FilterAPIRequestBody returns the object to be sent as request body. Not implementing this interface is equivalent to returning the received object from this function.
	FilterAPIRequestBody(op Operation, options Options) (interface{}, error)
}

// ResponseFilterHook is an interface Objects can optionally implement to modify response given by the engine before they are decoded.
type ResponseFilterHook interface {
	// FilterAPIResponse is called after a response from the engine regarding this object is received. Instead of the original response, the one returned by this function is decoded.
	FilterAPIResponse(op Operation, options Options, res *http.Response) (*http.Response, error)
}
