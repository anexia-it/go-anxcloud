package v1

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

var (
	// ErrTemplateNotFound is returned when the named template was not found at a given location
	ErrTemplateNotFound = fmt.Errorf("%w: named template was not found at specified location", api.ErrNotFound)
)

const (
	// LatestTemplateBuild is used to find the template with the highest build number
	LatestTemplateBuild = "latest"
)

// FindNamedTemplate retrieves a template by name and build at a specified location.
// Empty and LatestTemplateBuild build identifier will yield the highest available build.
// It returns ErrTemplateNotFound if no matching template was found.
func FindNamedTemplate(ctx context.Context, a api.API, name, build string, location corev1.Location) (*Template, error) {
	var (
		match    *Template
		fallback *Template
	)
	buildNo := -1
	latest := build == "" || build == LatestTemplateBuild

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var channel types.ObjectChannel

	err := a.List(ctx, &Template{Type: TypeTemplate, Location: location}, api.ObjectChannel(&channel))
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	for res := range channel {
		var template Template
		err := res(&template)
		if err != nil {
			return nil, fmt.Errorf("error retrieving template: %w", err)
		}

		if template.Name != name {
			continue
		}

		if latest {
			currentTemplateBuildNo, err := template.BuildNumber()
			if err != nil {
				fallback = &template
				logr.FromContextOrDiscard(ctx).Info("couldn't parse build %q of template %q from location %q", template.Build, template.Identifier, location.Identifier)
				continue
			}

			if latest && (match == nil || currentTemplateBuildNo > buildNo) {
				match = &template
				buildNo = currentTemplateBuildNo
			}
		} else if template.Build == build {
			match = &template
			break
		}

	}

	if match == nil && fallback == nil {
		return nil, fmt.Errorf("%w (name: %q, build: %q, location: %q)", ErrTemplateNotFound, name, build, location.Identifier)
	} else if match == nil {
		match = fallback
	}

	return match, nil
}
