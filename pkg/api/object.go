package api

import (
	"fmt"
	"reflect"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	uuid "github.com/satori/go.uuid"
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
	numFields := t.NumField()

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
					return ret, nil
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

				for t, vf := range allowedIdentifierTypes {
					if identifierValue.Type() == t {
						return vf(identifierValue.Interface()), nil
					}
				}

				return "", fmt.Errorf("%w: Objects identifier field has an unsupported type (type %v has an identifier of type %v)", ErrTypeNotSupported, t, field.Type)
			}
		}
	}

	return "", fmt.Errorf("%w: Object lacks identifier field (type %v does not have a field with `anxcloud:\"identifier\"` tag)", ErrTypeNotSupported, t)
}
