package v2

import (
	"context"
	"net/url"
)

// EndpointURL returns the URL for the empty and delete trigger endpoint
func (t *bucketEmptyAndDelete) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.Parse("/api/object_storage/v2/bucket")
	if err != nil {
		return nil, err
	}

	// Add the bucket identifier and trigger path for empty_and_delete
	u.Path = u.Path + "/" + t.BucketIdentifier + "/trigger/empty_and_delete"

	return u, nil
}

// GetIdentifier returns the bucket identifier for the trigger
func (t *bucketEmptyAndDelete) GetIdentifier(ctx context.Context) (string, error) {
	return t.BucketIdentifier, nil
}

// FilterAPIRequestBody generates the request body for the empty and delete trigger
func (t *bucketEmptyAndDelete) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return map[string]bool{
		"empty_and_delete": t.EmptyAndDelete,
	}, nil
}
