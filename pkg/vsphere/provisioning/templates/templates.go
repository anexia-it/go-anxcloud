// Package templates implements API functions residing under /provisioning/templates.
// This path contains methods for querying available VM templates.
package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TemplateType defines which type of template is selected.
type TemplateType string

// Template contains a summary about the state of a template.
type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	WordSize string `json:"bits"`
	Build    string `json:"build"`
}

const (
	// TemplateTypeTemplates are templates that already contain a distribution.
	TemplateTypeTemplates TemplateType = "templates"
	// TemplateTypeFromScratch are templates that need to have distribution added to work.
	TemplateTypeFromScratch TemplateType = "from_scratch"
	pathPrefix              string       = "/api/vsphere/v1/provisioning/templates.json"
)

func (a api) List(ctx context.Context, locationID string, templateType TemplateType, page, limit int) ([]Template, error) {
	url := fmt.Sprintf(
		"%s%s/%s/%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathPrefix, locationID, templateType, page, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create template list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute template list request: %w", err)
	}
	var responsePayload []Template
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode template list response: %w", err)
	}

	return responsePayload, err
}
