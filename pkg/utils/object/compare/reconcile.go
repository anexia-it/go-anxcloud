package compare

import (
	"fmt"
	"reflect"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

type reconcileObjectRetriever interface {
	Len() int
	Index(int) reflect.Value
	Object(int) types.Object
}

// Reconcile fills the given arrays of Objects to create and destroy, by adding the Objects in target but not in existing to create
// and adding the Objects in existing but not in target to destroy. Objects are compared with the Compare function in this package,
// using the provided compareAttributes.
//
// Objects in target with a matching Object in existing are updated, filling e.g. the identifiers in target.
func Reconcile(target, existing interface{}, create, destroy *[]types.Object, compareAttributes ...string) error {
	targetArray, existingArray, err := reconcileValidateInputArrays(target, existing)
	if err != nil {
		return err
	}

	for i := 0; i < targetArray.Len(); i++ {
		targetElement := targetArray.Index(i)

		if index, err := Search(targetElement.Interface(), existing, compareAttributes...); err == nil && index == -1 {
			*create = append(*create, targetArray.Object(i))
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
			*destroy = append(*destroy, existingArray.Object(i))
		} else if err != nil {
			// Not reached as the same error would also have happened in the loop above already and since Search
			// only returns the error when target contains elements, the loop above is definitely executed, too.
			return fmt.Errorf("error comparing existing and target resource: %w", err)
		}
	}

	return nil
}

func reconcileValidateInputArrays(target, existing interface{}) (targetRetriever, existingRetriever reconcileObjectRetriever, err error) {
	targetArray := reflect.ValueOf(target)
	existingArray := reflect.ValueOf(existing)

	targetRetriever, err = reconcileValidateInputArray(targetArray)
	if err != nil {
		return nil, nil, fmt.Errorf("target array invalid: %w", err)
	}

	existingRetriever, err = reconcileValidateInputArray(existingArray)
	if err != nil {
		return nil, nil, fmt.Errorf("existing array invalid: %w", err)
	}

	return
}

func reconcileValidateInputArray(v reflect.Value) (reconcileObjectRetriever, error) {
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: target and existing have to be arrays", ErrInvalidType)
	}

	objectType := reflect.TypeOf((*types.Object)(nil)).Elem()
	elemType := v.Type().Elem()

	if !elemType.Implements(objectType) {
		if !reflect.PointerTo(elemType).Implements(objectType) {
			return nil, fmt.Errorf("%w: input array elements do not implement types.Object", ErrInvalidType)
		}

		return reconcileIndirectObject{v}, nil
	}

	return reconcileDirectObject{v}, nil
}

type reconcileDirectObject struct {
	reflect.Value
}

func (rdo reconcileDirectObject) Object(i int) types.Object {
	return rdo.Index(i).Interface().(types.Object)
}

type reconcileIndirectObject struct {
	reflect.Value
}

func (rio reconcileIndirectObject) Object(i int) types.Object {
	return rio.Index(i).Addr().Interface().(types.Object)
}
