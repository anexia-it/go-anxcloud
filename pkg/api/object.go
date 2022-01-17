package api

import (
	"errors"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func getObjectIdentifier(obj types.Object, singleObjectOperation bool) (string, error) {
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
					reflect.TypeOf(""):       func(v interface{}) string { return v.(string) },
					reflect.TypeOf(uuid.Nil): func(v interface{}) string { return v.(uuid.UUID).String() },
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
