package types

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
	EndpointURL(ctx context.Context) (*url.URL, error)
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
	FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error)
}

// RequestBodyHook is an interface Objects can optionally implement to customize request bodies based on the object given by the user of this library.
type RequestBodyHook interface {
	// FilterAPIRequestBody returns the object to be sent as request body. Not implementing this interface is equivalent to returning the received object from this function.
	FilterAPIRequestBody(ctx context.Context) (interface{}, error)
}

// ResponseFilterHook is an interface Objects can optionally implement to modify response given by the engine before they are decoded.
type ResponseFilterHook interface {
	// FilterAPIResponse is called after a response from the engine regarding this object is received. Instead of the original response, the one returned by this function is decoded.
	FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error)
}

// PaginationSupportHook is an interface Objects can optionally implement to enable or disable pagination support for List operations.
type PaginationSupportHook interface {
	// Returns if the API supports pagination for List operations. Mind optionally supported filters in EndpointURL, which may go to different API endpoints which independently might
	// or might not support pagination.
	HasPagination(ctx context.Context) (bool, error)
}

// ResponseDecodeHook is an interface Objects can optionally implement to change the API response decode behavior
type ResponseDecodeHook interface {
	// Decodes the API response into the Object this function is called on
	DecodeAPIResponse(ctx context.Context, data io.Reader) error
}

// FilterRequestURLHook is an interface Objects can optionally implement to modify the request URL
type FilterRequestURLHook interface {
	// FilterRequestURL returns the modified URL
	FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error)
}

// GetObjectIdentifier extracts the identifier of the given object, returning an error if no identifier field
// is found or singleObjectOperation is true and an identifier field is found, but empty.
func GetObjectIdentifier(obj Object, singleObjectOperation bool) (string, error) {
	objectType := reflect.TypeOf(obj)

	if objectType.Kind() != reflect.Ptr {
		return "", fmt.Errorf("%w: the Object interface must be implemented on a pointer to struct", ErrTypeNotSupported)
	} else if objectType.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("%w: Objects must be implemented as structs", ErrTypeNotSupported)
	}

	objectStructType := objectType.Elem()
	return findIdentifierInStruct(objectStructType, reflect.ValueOf(obj).Elem(), singleObjectOperation)
}

func findIdentifierInStruct(t reflect.Type, v reflect.Value, singleObjectOp bool) (string, error) {
	// we also use this to track if we found an identifier already
	var returnIdentifier *string

	numFields := t.NumField()

fields:
	for i := 0; i < numFields; i++ {
		field := t.Field(i)

		if field.Anonymous {
			embeddedType := field.Type
			embeddedValue := v.Field(i)

			for embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
				embeddedValue = embeddedValue.Elem()
			}

			if embeddedType.Kind() == reflect.Struct {
				if ret, err := findIdentifierInStruct(embeddedType, embeddedValue, singleObjectOp); err == nil {
					if returnIdentifier == nil {
						returnIdentifier = &ret
					} else {
						return "", fmt.Errorf("%w (type %v has multiple fields tagged as identifier)", ErrObjectWithMultipleIdentifier, t)
					}
				} else if errors.Is(err, ErrObjectWithMultipleIdentifier) || errors.Is(err, ErrObjectIdentifierTypeNotSupported) {
					return "", err
				}
			}

			continue
		}

		if val, ok := field.Tag.Lookup("anxcloud"); ok {
			if val == "identifier" {
				identifierValue := v.Field(i)

				// We check on the value to have a type-independent zero check, in case we later allow other
				// types for identifier. A int identifier is zero with value 0, which encoded to string "0",
				// so a later identifier == "" check would not work.
				if singleObjectOp && identifierValue.IsZero() {
					return "", ErrUnidentifiedObject
				}

				allowedIdentifierTypes := map[reflect.Type]func(interface{}) string{
					reflect.TypeOf(""): func(v interface{}) string { return v.(string) },
				}

				for ft, vf := range allowedIdentifierTypes {
					if identifierValue.Type() == ft {
						if returnIdentifier == nil {
							val := vf(identifierValue.Interface())
							returnIdentifier = &val

							continue fields
						} else {
							return "", fmt.Errorf("%w (type %v has multiple fields tagged as identifier)", ErrObjectWithMultipleIdentifier, t)
						}
					}
				}

				return "", fmt.Errorf("%w (type %v has an identifier of type %v)", ErrObjectIdentifierTypeNotSupported, t, field.Type)
			}
		}
	}

	if returnIdentifier != nil {
		return *returnIdentifier, nil
	}

	return "", fmt.Errorf("%w (type %v does not have a field with `anxcloud:\"identifier\"` tag)", ErrObjectWithoutIdentifier, t)
}
