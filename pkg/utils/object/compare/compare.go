package compare

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrInvalidType = errors.New("invalid type given")

	ErrDifferentTypes = fmt.Errorf("%w: cannot compare structs of different types", ErrInvalidType)

	ErrKeyNotFound = fmt.Errorf("key not found in struct")
)

// Difference gives the (nested) attribute name and the values of the attribute of the two objects differing.
type Difference struct {
	Key string
	A   interface{}
	B   interface{}
}

// Compare compares two objects by the given (dot-nested) attribute names.
func Compare(a, b interface{}, attributes ...string) ([]Difference, error) {
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)

	if typeA.Kind() == reflect.Ptr {
		typeA = typeA.Elem()
	}

	if typeB.Kind() == reflect.Ptr {
		typeB = typeB.Elem()
	}

	if typeA.Kind() != reflect.Struct || typeB.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: can only compare structs and pointers to structs", ErrInvalidType)
	}

	if typeA != typeB {
		return nil, ErrDifferentTypes
	}

	attributeIndexes := make([][]int, 0, len(attributes))

	for _, key := range attributes {
		keyParts := strings.Split(key, ".")
		indexes := make([]int, 0, len(keyParts))

		for _, keyPart := range keyParts {
			if field, ok := typeA.FieldByName(keyPart); !ok {
				return nil, fmt.Errorf("%w: key %v", ErrKeyNotFound, keyPart)
			} else {
				indexes = append(indexes, field.Index...)
			}
		}

		attributeIndexes = append(attributeIndexes, indexes)
	}

	ret := make([]Difference, 0)

	for i, fieldKey := range attributeIndexes {
		valA, nilIdxA := fieldByIndex(reflect.ValueOf(a), fieldKey)
		valB, nilIdxB := fieldByIndex(reflect.ValueOf(b), fieldKey)

		equal := true

		if nilIdxA != nilIdxB {
			equal = false
		}

		if equal && nilIdxA == -1 && !reflect.DeepEqual(valA, valB) {
			equal = false
		}

		if !equal {
			ret = append(ret, Difference{
				Key: attributes[i],
				A:   valA,
				B:   valB,
			})
		}
	}

	return ret, nil
}

// fieldByIndex works like reflect.Value.FieldByIndex but does not panic if any value on the idxs-path is nil,
// instead returning nil and the index in idxs being nil. It also returns fieldValue.Interfac() instead of fieldValue
func fieldByIndex(v reflect.Value, idxs []int) (retVal interface{}, nilIdx int) {
	current := v
	for i, idx := range idxs {
		if current.Kind() == reflect.Ptr && current.IsNil() {
			return nil, i
		}

		current = reflect.Indirect(current).Field(idx)
	}

	return current.Interface(), -1
}
