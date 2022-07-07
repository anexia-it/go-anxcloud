// DO NOT EDIT, auto generated

package v1

import (
	"context"
)

// GetIdentifier returns the primary identifier of a Resource object
func (x *Resource) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a ResourceWithTag object
func (x *ResourceWithTag) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}
