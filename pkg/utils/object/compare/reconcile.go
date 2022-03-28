package compare

import (
	"fmt"
	"reflect"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// Reconcile fills the given arrays of Objects to create and destroy, by adding the Objects in target but not in existing to create
// and adding the Objects in existing but not in target to destroy. Objects are compared with the Compare function in this package,
// using the provided compareAttributes.
//
// Objects in target with a matching Object in existing are updated, filling e.g. the identifiers in target.
func Reconcile(target, existing interface{}, create, destroy *[]types.Object, compareAttributes ...string) error {
	targetArray := reflect.ValueOf(target)
	existingArray := reflect.ValueOf(existing)

	if targetArray.Kind() != reflect.Slice ||
		existingArray.Kind() != reflect.Slice {
		return fmt.Errorf("%w: target and existing have to be arrays", ErrInvalidType)
	}

	for i := 0; i < targetArray.Len(); i++ {
		targetElement := targetArray.Index(i)

		if index, err := Search(targetElement.Interface(), existing, compareAttributes...); err == nil && index == -1 {
			*create = append(*create, targetElement.Addr().Interface().(types.Object))
		} else if err != nil {
			return fmt.Errorf("error comparing existing and target resource: %w", err)
		} else {
			// update entry in target array with the existing Object
			targetElement.Set(existingArray.Index(index))
		}
	}

	for i := 0; i < existingArray.Len(); i++ {
		existingElement := existingArray.Index(i)

		if index, err := Search(existingElement.Interface(), target, compareAttributes...); err == nil && index == -1 {
			*destroy = append(*destroy, existingElement.Addr().Interface().(types.Object))
		} else if err != nil {
			// Not reached as the same error would also have happened in the loop above already and since Search
			// only returns the error when target contains elements, the loop above is definitely executed, too.
			return fmt.Errorf("error comparing existing and target resource: %w", err)
		}
	}

	return nil
}
