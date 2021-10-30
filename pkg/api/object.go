package api

import (
	"fmt"
	"reflect"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

func getObjectIdentifier(obj types.Object, singleObjectOperation bool) (string, error) {
	objectType := reflect.TypeOf(obj)

	if objectType.Kind() != reflect.Ptr {
		return "", fmt.Errorf("%w: the Object interface must be implemented on a pointer to struct", ErrTypeNotSupported)
	} else if objectType.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("%w: Objects must be implemented as structs", ErrTypeNotSupported)
	}

	objectStructType := objectType.Elem()
	numFields := objectStructType.NumField()

	for i := 0; i < numFields; i++ {
		field := objectStructType.Field(i)

		if val, ok := field.Tag.Lookup("anxcloud"); ok {
			if val == "identifier" {
				identifierValue := reflect.ValueOf(obj).Elem().Field(i)

				// We check on the value to have a type-independent zero check, in case we later allow other
				// types for identifier. A int identifier is zero with value 0, which encoded to string "0",
				// so a later identifier == "" check would not work.
				if singleObjectOperation && identifierValue.IsZero() {
					return "", ErrUnidentifiedObject
				}

				// TODO: maybe we need to support other types, too - any scalar, string-convertable type should do just fine
				switch identifierValue.Kind() {
				case reflect.String:
					return identifierValue.String(), nil
				default:
					return "", fmt.Errorf("%w: Objects identifier field has an unsupported type (type %v has an identifier of type %v)", ErrTypeNotSupported, objectStructType, field.Type)
				}
			}
		}
	}

	return "", fmt.Errorf("%w: Object lacks identifier field (type %v does not have a field with `anxcloud:\"identifier\"` tag)", ErrTypeNotSupported, objectStructType)
}
