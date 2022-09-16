// Package helper exists to break a circular-dependency between pkg/api and pkg/apis/core/v1, that
// was introduced when adding the AutoTag option to the generic client.
// Be careful what you use from this package, it should not contain anything valuable for users outside
// go-anxcloud itself.
package helper

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// Tagger is a helper to Tag or Untag generic client objects.
type Tagger interface {
	Tag(context.Context, types.API, types.IdentifiedObject, ...string) error
	Untag(context.Context, types.API, types.IdentifiedObject, ...string) error
	ListTags(context.Context, types.API, types.IdentifiedObject) ([]string, error)
}

// TaggerImplementation is the Tagger to use, set on startup of a program. This is a really ugly workaround
// for a circular-dependency between pkg/api and pkg/apis/core/v1 for automatically tagging Objects after
// they are created.
// Please do not use this directly.
var TaggerImplementation Tagger
