package v1

import (
	"context"
	"fmt"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/core/v1/helper"
)

// Tag adds tags to an object resource
func Tag(ctx context.Context, a types.API, obj types.IdentifiedObject, tags ...string) error {
	objects, err := resourceWithTagObjects(obj, tags...)
	if err != nil {
		return fmt.Errorf("generating ResourceWithTag objects failed: %w", err)
	}

	for _, obj := range objects {
		if err := a.Create(ctx, obj); err != nil {
			if err, ok := err.(api.HTTPError); ok && err.StatusCode() == 422 {
				// already tagged -> skip
				continue
			}
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
		objects = append(objects, &ResourceWithTag{Identifier: identifier, Tag: tag})
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
