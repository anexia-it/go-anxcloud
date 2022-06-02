// Package filter implements a helper for Objects supporting filtered List operations.
package filter

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

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

func parseFilterableTag(field reflect.StructField) (isFilterable bool, filterName string, option string) {
	filterName = field.Name

	tag, ok := field.Tag.Lookup("anxcloud")
	if !ok {
		return false, "", ""
	}

	tagParts := strings.Split(tag, ",")

	if len(tagParts) == 0 || tagParts[0] != "filterable" {
		return false, "", ""
	}

	isFilterable = true

	if len(tagParts) >= 2 && tagParts[1] != "" {
		filterName = tagParts[1]
	} else {
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(jsonTag, ",")

			if parts[0] != "" || len(parts) > 1 {
				filterName = parts[0]
			}
		}
	}

	if len(tagParts) >= 3 {
		option = tagParts[2]
	}

	return isFilterable, filterName, option
}

func parseOptionSingle(fieldValue reflect.Value) (reflect.Value, error) {
	fieldType := fieldValue.Type()
	fieldKind := fieldType.Kind()

	if fieldKind == reflect.Slice || fieldKind == reflect.Array {
		if fieldValue.Len() == 1 {
			return fieldValue.Index(0), nil
		} else if fieldValue.Len() == 0 {
			return reflect.New(fieldType.Elem()).Elem(), nil
		} else {
			return reflect.Value{}, fmt.Errorf("%w: only a single value can be filtered", types.ErrInvalidFilter)
		}
	} else {
		return reflect.Value{}, fmt.Errorf("%w: option 'single' can only be used on array or slice attributes", types.ErrInvalidFilter)
	}
}

func extractFilterValue(fieldValue reflect.Value) (reflect.Value, error) {
	fieldType := fieldValue.Type()
	fieldKind := fieldType.Kind()

	if fieldKind == reflect.Ptr {
		fieldType = fieldType.Elem()
		fieldKind = fieldType.Kind()

		if fieldValue.IsNil() || fieldValue.IsZero() {
			return reflect.Zero(fieldType), nil
		}

		fieldValue = fieldValue.Elem()
	}

	if fieldKind == reflect.Struct {
		if object, ok := fieldValue.Addr().Interface().(types.Object); ok {
			identifier, err := types.GetObjectIdentifier(object, false)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("Object referenced: %w", err)
			}

			fieldValue = reflect.ValueOf(identifier)
		}
	}

	return fieldValue, nil
}

func (f *filterHelper) parseField(fieldValue reflect.Value, field reflect.StructField) error {
	filterable, filterName, option := parseFilterableTag(field)

	if !filterable {
		return nil
	}

	f.fields[filterName] = true

	if option == "single" {
		if val, err := parseOptionSingle(fieldValue); err != nil {
			return fmt.Errorf("field %q: %w", filterName, err)
		} else {
			fieldValue = val
		}
	}

	fieldValue, err := extractFilterValue(fieldValue)
	if err != nil {
		return err
	}

	if isSupportedPrimitive(fieldValue.Type().Kind()) && !fieldValue.IsZero() {
		f.values[filterName] = fieldValue.Interface()
	}

	return nil
}

func (f *filterHelper) parseObject(v interface{}) error {
	val := reflect.ValueOf(v)
	valType := val.Type()

	if valType.Kind() == reflect.Ptr {
		val = val.Elem()
		valType = val.Type()
	}

	if valType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: filter.Helper only works with structs or pointers to them", types.ErrTypeNotSupported)
	}

	numFields := val.NumField()
	for i := 0; i < numFields; i++ {
		if err := f.parseField(val.Field(i), valType.Field(i)); err != nil {
			return err
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
