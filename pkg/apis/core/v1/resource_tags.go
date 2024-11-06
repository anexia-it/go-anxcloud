package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/core/v1/helper"
	"go.anx.io/go-anxcloud/pkg/utils/retry"
)

// Tag adds tags to an object resource and internally retries 2 more times on conflict
func Tag(ctx context.Context, a types.API, obj types.IdentifiedObject, tags ...string) error {
	objects, err := resourceWithTagObjects(obj, tags...)
	if err != nil {
		return fmt.Errorf("generating ResourceWithTag objects failed: %w", err)
	}
	for _, obj := range objects {
		tag := func() (bool, error) {
			var (
				httpError api.HTTPError
				retryable bool
				err       error
			)
			if err = a.Create(ctx, obj); err != nil && errors.As(err, &httpError) {
				if httpError.StatusCode() == http.StatusUnprocessableEntity {
					// already tagged -> skip
					err = nil
				} else if httpError.StatusCode() == http.StatusConflict {
					retryable = true
				}
			}
			return retryable, err
		}

		if err := retry.Retry(ctx, 3, time.Second, tag); err != nil {
			return err
		}
	}

	return nil
}

// Untag removes tags from an object resource
func Untag(ctx context.Context, a types.API, obj types.IdentifiedObject, tags ...string) error {
	objects, err := resourceWithTagObjects(obj, tags...)
	if err != nil {
		return fmt.Errorf("generating ResourceWithTag objects failed: %w", err)
	}

	for _, obj := range objects {
		if err := a.Destroy(ctx, obj); api.IgnoreNotFound(err) != nil {
			return err
		}
	}
	return nil
}

func resourceWithTagObjects(obj types.IdentifiedObject, tags ...string) ([]*ResourceWithTag, error) {
	identifier, err := types.GetObjectIdentifier(obj, true)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving Object identifier: %w", err)
	}

	objects := make([]*ResourceWithTag, 0, len(tags))

	for _, tag := range tags {
		objects = append(objects, &ResourceWithTag{ResourceIdentifier: identifier, Tag: tag})
	}

	return objects, nil
}

// ListTags retrieves tags of given object
func ListTags(ctx context.Context, a types.API, obj types.IdentifiedObject) ([]string, error) {
	identifier, err := types.GetObjectIdentifier(obj, true)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving Object identifier: %w", err)
	}

	r := &Resource{Identifier: identifier}
	if err := a.Get(ctx, r); err != nil {
		return nil, err
	}

	return r.Tags, nil
}

type taggerImplementation int

func (ti taggerImplementation) Tag(ctx context.Context, a types.API, obj types.IdentifiedObject, tags ...string) error {
	return Tag(ctx, a, obj, tags...)
}

func (ti taggerImplementation) Untag(ctx context.Context, a types.API, obj types.IdentifiedObject, tags ...string) error {
	return Untag(ctx, a, obj, tags...)
}

func (ti taggerImplementation) ListTags(ctx context.Context, a types.API, obj types.IdentifiedObject) ([]string, error) {
	return ListTags(ctx, a, obj)
}

func init() {
	// This is a workaround to solve import cycles between `pkg/apis/core/v1` <--> `pkg/api`
	// initially caused by the AutoTag Create option.
	helper.TaggerImplementation = taggerImplementation(42)
}
