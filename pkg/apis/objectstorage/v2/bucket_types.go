package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object:hooks=RequestBodyHook

// Bucket represents a bucket resource in the Object Storage API.
type Bucket struct {
	gs.GenericService
	gs.HasState

	CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
	ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
	Identifier         string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	Tags               gs.PartialResourceList `json:"tags,omitempty"`
	Reseller           string                 `json:"reseller,omitempty"`
	Customer           string                 `json:"customer,omitempty"`
	Share              bool                   `json:"share,omitempty"`

	Name               string                 `json:"name"`
	State              *GenericAttributeState `json:"state,omitempty"`
	Region             common.PartialResource `json:"region"`
	ObjectCount        interface{}            `json:"object_count,omitempty"`
	ObjectSize         interface{}            `json:"object_size,omitempty"`
	Backend            common.PartialResource `json:"backend"`
	Tenant             common.PartialResource `json:"tenant"`
	ObjectLockLifetime *int                   `json:"object_lock_lifetime,omitempty"`
	VersioningActive   bool                   `json:"versioning_active,omitempty"`
	Embed              []string               `json:"-"`
}

// GetObjectCount returns the object count as a float64, handling both string and numeric values.
func (b *Bucket) GetObjectCount() (float64, error) {
	if b.ObjectCount == nil {
		return 0, nil
	}

	switch v := b.ObjectCount.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, nil
	}
}

// GetObjectSize returns the object size as a float64, handling both string and numeric values.
func (b *Bucket) GetObjectSize() (float64, error) {
	if b.ObjectSize == nil {
		return 0, nil
	}

	switch v := b.ObjectSize.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, nil
	}
}

// DecodeAPIResponse handles custom JSON unmarshaling for Bucket to fix type mismatches
func (b *Bucket) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	// Create a temporary struct with the problematic fields as interface{}
	var temp struct {
		gs.GenericService
		gs.HasState

		CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
		ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
		Identifier         string                 `json:"identifier,omitempty"`
		Tags               gs.PartialResourceList `json:"tags,omitempty"`
		Reseller           string                 `json:"reseller,omitempty"`
		Customer           string                 `json:"customer,omitempty"`
		Share              bool                   `json:"share,omitempty"`

		Name        string                 `json:"name"`
		State       *GenericAttributeState `json:"state,omitempty"`
		Region      common.PartialResource `json:"region"`
		ObjectCount interface{}            `json:"object_count,omitempty"`
		ObjectSize  interface{}            `json:"object_size,omitempty"`
		Backend     common.PartialResource `json:"backend"`
		Tenant      common.PartialResource `json:"tenant"`
	}

	// Unmarshal into the temp struct
	if err := json.NewDecoder(data).Decode(&temp); err != nil {
		return fmt.Errorf("failed to decode bucket JSON: %w", err)
	}

	// Copy all fields to the bucket
	b.GenericService = temp.GenericService
	b.HasState = temp.HasState
	b.CustomerIdentifier = temp.CustomerIdentifier
	b.ResellerIdentifier = temp.ResellerIdentifier
	b.Identifier = temp.Identifier
	b.Tags = temp.Tags
	b.Reseller = temp.Reseller
	b.Customer = temp.Customer
	b.Share = temp.Share
	b.Name = temp.Name
	b.State = temp.State
	b.Region = temp.Region
	b.Backend = temp.Backend
	b.Tenant = temp.Tenant

	// Handle the interface{} fields that can be strings or numbers
	b.ObjectCount = temp.ObjectCount
	b.ObjectSize = temp.ObjectSize

	return nil
}

// bucketEmptyAndDelete represents the trigger object for empty and delete operations
type bucketEmptyAndDelete struct {
	BucketIdentifier string `json:"-" anxcloud:"identifier"`
	EmptyAndDelete   bool   `json:"empty_and_delete"`
}

// EmptyAndDelete empties the bucket and then deletes it using the trigger/empty_and_delete endpoint.
// This is the proper way to delete a bucket that contains objects.
func EmptyAndDelete(ctx context.Context, a api.API, bucketID string) error {
	trigger := &bucketEmptyAndDelete{
		BucketIdentifier: bucketID,
		EmptyAndDelete:   true,
	}
	return a.Create(ctx, trigger)
}

// EmptyAndDelete is a convenience method that calls the package-level EmptyAndDelete function
// using the bucket's identifier.
func (b *Bucket) EmptyAndDelete(ctx context.Context, a api.API) error {
	if b.Identifier == "" {
		return fmt.Errorf("bucket identifier is required for empty and delete operation")
	}
	return EmptyAndDelete(ctx, a, b.Identifier)
}
