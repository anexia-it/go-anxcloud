package v1

import (
	"errors"
	"fmt"
	"strconv"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

var (
	// ErrFailedToParseTemplateBuildNumber is returned when the template build couldn't be converted to an int
	ErrFailedToParseTemplateBuildNumber = errors.New("failed to parse template build number")
)

// anxcloud:object:hooks=PaginationSupportHook,FilterRequestURLHook,ResponseDecodeHook

// Template represents a vSphere template used for vm provisioning
type Template struct {
	Identifier string          `json:"id" anxcloud:"identifier"`
	Name       string          `json:"name"`
	Bit        string          `json:"bit"`
	Build      string          `json:"build"`
	Location   corev1.Location `json:"-"`
	Type       TemplateType    `json:"-"`
}

// BuildNumber returns the parsed build number
func (t *Template) BuildNumber() (int, error) {
	if t.Build == "" || t.Build[0] != 'b' {
		return 0, fmt.Errorf("%w: template build does not start with \"b\"", ErrFailedToParseTemplateBuildNumber)
	}

	buildNumber, err := strconv.Atoi(t.Build[1:])
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrFailedToParseTemplateBuildNumber, err)
	}

	return buildNumber, nil
}

// TemplateType specifies the type of template
type TemplateType string

const (
	// TypeTemplate is used for prebuilt templates
	TypeTemplate TemplateType = "templates"

	// TypeFromScratch is used for custom templates
	TypeFromScratch TemplateType = "from_scratch"
)
