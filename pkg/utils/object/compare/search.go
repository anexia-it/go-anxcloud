package compare

import (
	"reflect"
)

// Search iterates through the array haystack until it finds an Object matching needle by the given
// (nested) attribute names (using the Compare function in this package).
func Search(needle, haystack interface{}, compareAttributes ...string) (int, error) {
	haystackValue := reflect.ValueOf(haystack)
	for i := 0; i < haystackValue.Len(); i++ {
		haystackElem := haystackValue.Index(i).Interface()
		compare, err := Compare(needle, haystackElem, compareAttributes...)

		if err == nil && len(compare) == 0 {
			return i, nil
		} else if err != nil {
			return -1, err
		}
	}

	return -1, nil
}
