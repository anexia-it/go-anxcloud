package types

import "errors"

var (
	// ErrKeyNotSet is returned when trying to options.Get(key) a key that is not set.
	ErrKeyNotSet = errors.New("requested key not set on the given Options")

	// ErrKeyAlreadySet is returned when trying to options.Set(key, v, false) and the key that is already set.
	ErrKeyAlreadySet = errors.New("given key is already set on the given Options")
)

type commonOptions struct {
	additional map[string]interface{}
}

// Options is the interface all operation-specific options implement, making it possible to pass all the specific options to the same functions.
// Specific APIs can use this interface to set additional options, keys should be prefixed with the name of the API. This might be enforced in the future.
type Options interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, overwrite bool) error
}

// GetOption is the interface options have to implement to be usable with Get operation.
type GetOption interface {
	// Apply this option to the set of all options
	ApplyToGet(*GetOptions) error
}

// ListOption is the interface options have to implement to be usable with List operation.
type ListOption interface {
	// Apply this option to the set of all options
	ApplyToList(*ListOptions) error
}

// CreateOption is the interface options have to implement to be usable with Create operation.
type CreateOption interface {
	// Apply this option to the set of all options
	ApplyToCreate(*CreateOptions) error
}

// UpdateOption is the interface options have to implement to be usable with Update operation.
type UpdateOption interface {
	// Apply this option to the set of all options
	ApplyToUpdate(*UpdateOptions) error
}

// DestroyOption is the interface options have to implement to be usable with Destroy operation.
type DestroyOption interface {
	// Apply this option to the set of all options
	ApplyToDestroy(*DestroyOptions) error
}

type AnyOption interface {
	ApplyToAny(Options) error
}

type Option interface{}

func (o commonOptions) Get(key string) (interface{}, error) {
	if o.additional != nil {
		if v, ok := o.additional[key]; ok {
			return v, nil
		}
	}

	return nil, ErrKeyNotSet
}

func (o *commonOptions) Set(key string, val interface{}, overwrite bool) error {
	if o.additional == nil {
		o.additional = make(map[string]interface{}, 1)
	}

	if _, alreadySet := o.additional[key]; alreadySet && !overwrite {
		return ErrKeyAlreadySet
	}

	o.additional[key] = val
	return nil
}
