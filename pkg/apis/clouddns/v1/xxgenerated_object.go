// Code generated by go.anx.io/go-anxcloud/tools object-generator - DO NOT EDIT!

package v1

import (
	"context"
)

// GetIdentifier returns the primary identifier of a Record object
func (o *Record) GetIdentifier(ctx context.Context) (string, error) {
	return o.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Revision object
func (o *Revision) GetIdentifier(ctx context.Context) (string, error) {
	return o.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Zone object
func (o *Zone) GetIdentifier(ctx context.Context) (string, error) {
	return o.Name, nil
}
