package mock

import (
	"errors"
	"reflect"
	"strings"

	"github.com/mitchellh/copystructure"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/test"
)

var (
	// ErrMergeTypeMissmatch is returned when the merge function is called with two different types
	ErrMergeTypeMissmatch = errors.New("merge expected two objects of same type")
)

func uintMin(x, y uint) uint {
	if x > y {
		return y
	}
	return x
}

func merge(dest, src interface{}) (interface{}, error) {
	if reflect.TypeOf(dest) != reflect.TypeOf(src) {
		return nil, ErrMergeTypeMissmatch
	}

	destClone, err := copystructure.Copy(dest)
	if err != nil {
		return nil, err
	}

	_mergeStructs(
		reflect.Indirect(reflect.ValueOf(destClone)),
		reflect.Indirect(reflect.ValueOf(src)),
		reflect.TypeOf(dest).Elem(),
	)

	return destClone, nil
}

func _mergeStructs(dest, src reflect.Value, t reflect.Type) {
	for i := 0; i < dest.NumField(); i++ {
		if (t.Field(i).Type.Kind() == reflect.Ptr && src.Field(i).IsNil()) ||
			(t.Field(i).Type.Kind() != reflect.Ptr && strings.Contains(t.Field(i).Tag.Get("json"), ",omitempty") && src.Field(i).IsZero()) {
			// skip nil pointer and omitempty
			continue
		}

		if t.Field(i).Type.Kind() == reflect.Struct {
			_mergeStructs(dest.Field(i), src.Field(i), t.Field(i).Type)
		} else {
			dest.Field(i).Set(src.Field(i))
		}
	}
}

func makeObjectIdentifiable(o types.Object) string {
	identifier := test.RandomIdentifier()
	v := reflect.Indirect(reflect.ValueOf(o))
	t := reflect.TypeOf(o).Elem()

	if t.Kind() != reflect.Struct {
		return identifier
	}

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup("anxcloud"); ok &&
			tag == "identifier" &&
			t.Field(i).Type.Kind() == reflect.String {
			if v.Field(i).String() == "" {
				v.Field(i).SetString(identifier)
			} else {
				// might be already set (DNS Zone has zone-name as identifier)
				identifier = v.Field(i).String()
			}
		}
	}
	return identifier
}
