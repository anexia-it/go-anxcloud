package types

import "errors"

var (
	// ErrKeyNotSet is returned when trying to options.Get(key) a key that is not set.
	ErrKeyNotSet = errors.New("requested key not set on the given Options")

	// ErrKeyAlreadySet is returned when trying to options.Set(key, v, false) and the key that is already set.
	ErrKeyAlreadySet = errors.New("given key is already set on the given Options")
)

type commonOptions struct {
	additional   map[string]interface{}
	environments map[string]string
}

// Options is the interface all operation-specific options implement, making it possible to pass all the specific options to the same functions.
// Specific APIs can use this interface to set additional options, keys should be prefixed with the name of the API. This might be enforced in the future.
type Options interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, overwrite bool) error
	GetEnvironment(key string) (string, error)
	SetEnvironment(key, value string, overwrite bool) error
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

// AnyOption is the interface options have to implement to be usable with any operation.
type AnyOption func(Options) error

// ApplyToGet applies the AnyOption to the GetOptions
func (ao AnyOption) ApplyToGet(opts *GetOptions) error {
	return ao(opts)
}

// ApplyToList applies the AnyOption to the ListOptions
func (ao AnyOption) ApplyToList(opts *ListOptions) error {
	return ao(opts)
}

// ApplyToCreate applies the AnyOption to the CreateOptions
func (ao AnyOption) ApplyToCreate(opts *CreateOptions) error {
	return ao(opts)
}

// ApplyToUpdate applies the AnyOption to the UpdateOptions
func (ao AnyOption) ApplyToUpdate(opts *UpdateOptions) error {
	return ao(opts)
}

// ApplyToDestroy applies the AnyOption to the DestroyOptions
func (ao AnyOption) ApplyToDestroy(opts *DestroyOptions) error {
	return ao(opts)
}

// Option is a dummy interface used for any type of request Options
type Option interface{}

// Get retrieves a custom value from request options
func (o commonOptions) Get(key string) (interface{}, error) {
	if o.additional != nil {
		if v, ok := o.additional[key]; ok {
			return v, nil
		}
	}

	return nil, ErrKeyNotSet
}

// Set stores a custom value in request options
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

// GetEnvironment retrieves an environment value from request options
func (o commonOptions) GetEnvironment(key string) (string, error) {
	if o.environments != nil {
		if v, ok := o.environments[key]; ok {
			return v, nil
		}
	}

	return "", ErrKeyNotSet
}

// SetEnvironment stores an environment value in request options
func (o *commonOptions) SetEnvironment(key, val string, overwrite bool) error {
	if o.environments == nil {
		o.environments = make(map[string]string, 1)
	}

	if _, alreadySet := o.environments[key]; alreadySet && !overwrite {
		return ErrKeyAlreadySet
	}

	o.environments[key] = val
	return nil
}
