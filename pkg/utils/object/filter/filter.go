// Package filter implements a helper for Objects supporting filtered List operations.
package filter

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// NewHelper creates a new instance of a filter.Helper from the given object.
//
// The Helper will be set up for the fields in the given object (which commonly is a generic client Object, but
// has to be a struct or pointer to one) that have a `anxcloud:"filterable"` tag, retrieving their values,
// allowing easy access and building the filter query out of them automatically.
//
// The tag has an optional second field, `anxcloud:"filterable,foo", allowing to rename the field in the query
// to "foo". If no name is given, the name given in the encoding/json tag is used, if that is not given either,
// the name of the field is used.
//
// References to generic client Objects are resolved to their identifier, making the filter not set when the
// identifier of the referenced Object is empty.
func NewHelper(o interface{}) (Helper, error) {
	helper := filterHelper{
		values: make(map[string]interface{}),
		fields: make(map[string]bool),
	}

	err := helper.parseObject(o)
	if err != nil {
		return nil, err
	}

	return helper, nil
}

type filterHelper struct {
	values map[string]interface{}
	fields map[string]bool
}

// Get returns the value and if it was set for a given named field.
func (f filterHelper) Get(field string) (interface{}, bool, error) {
	if _, ok := f.fields[field]; !ok {
		return nil, false, fmt.Errorf("%w: field %q is not configured as filterable", ErrUnknownField, field)
	}

	v, ok := f.values[field]
	return v, ok, nil
}

// BuildQuery returns the query parameters to set for filtering.
func (f filterHelper) BuildQuery() url.Values {
	values := make(url.Values)

	for field, value := range f.values {
		// we also support numbers and stuff - this should work with all we have left
		values.Set(field, fmt.Sprintf("%v", value))
	}

	return values
}

func (f *filterHelper) parseObject(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Type().Kind() != reflect.Struct {
		return fmt.Errorf("%w: filter.Helper only works with structs or pointers to them", api.ErrTypeNotSupported)
	}

	numFields := val.NumField()
	for i := 0; i < numFields; i++ {
		field := val.Type().Field(i)
		fieldName := field.Name

		tag, ok := field.Tag.Lookup("anxcloud")
		if !ok {
			continue
		}

		tagParts := strings.Split(tag, ",")

		if len(tagParts) == 0 || tagParts[0] != "filterable" {
			continue
		}

		if len(tagParts) >= 2 && tagParts[1] != "" {
			fieldName = tagParts[1]
		} else {
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				parts := strings.Split(jsonTag, ",")

				if parts[0] != "" || len(parts) > 1 {
					fieldName = parts[0]
				}
			}
		}

		f.fields[fieldName] = true

		fieldValue := val.Field(i)
		fieldType := fieldValue.Type()
		fieldKind := fieldType.Kind()

		if fieldKind == reflect.Slice || fieldKind == reflect.Array {
			// only filter on first entry of an array
			if len(tagParts) >= 3 && tagParts[2] == "single" {
				if fieldValue.Len() == 1 {
					fieldValue = fieldValue.Index(0)
				} else if fieldValue.Len() == 0 {
					fieldValue = reflect.New(fieldType.Elem()).Elem()
				} else {
					return fmt.Errorf("%w: only a single value can be filtered for %q", types.ErrInvalidFilter, fieldName)
				}

				fieldType = fieldValue.Type()
				fieldKind = fieldType.Kind()
			}
		}

		if isSupportedPrimitive(fieldKind) {
			if !fieldValue.IsZero() {
				f.values[fieldName] = fieldValue.Interface()
			}
		} else if fieldKind == reflect.Ptr || fieldKind == reflect.Struct {
			if fieldKind == reflect.Ptr {
				if fieldValue.IsNil() || fieldValue.IsZero() {
					continue
				}

				if isSupportedPrimitive(fieldType.Elem().Kind()) {
					f.values[fieldName] = fieldValue.Elem().Interface()
				}
			}

			maybeObject := fieldValue

			if fieldKind == reflect.Struct {
				maybeObject = maybeObject.Addr()
			}

			if object, ok := maybeObject.Interface().(types.Object); ok {
				identifier, err := api.GetObjectIdentifier(object, false)
				if err != nil {
					return fmt.Errorf("Object referenced in %v: %w", fieldName, err)
				}

				if identifier != "" {
					f.values[fieldName] = identifier
				}
			}
		}
	}

	return nil
}

func isSupportedPrimitive(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}
