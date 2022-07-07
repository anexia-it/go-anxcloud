// DO NOT EDIT, auto generated

package v1

import (
	"context"
)

// GetIdentifier returns the primary identifier of a Record object
func (x *Record) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Zone object
func (x *Zone) GetIdentifier(ctx context.Context) (string, error) {
	return x.Name, nil
}
