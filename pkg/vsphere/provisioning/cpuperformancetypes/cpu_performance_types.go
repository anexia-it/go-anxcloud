// Package templates implements API functions residing under /provisioning/cpuperformancetype.
// This path contains methods for querying available VM templates.
package cpuperformancetype

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CPUPerformanceType
type CPUPerformanceType struct {
	ID             string  `json:"id"`
	Prioritization string  `json:"prioritization"`
	Limit          float64 `json:"limit"`
	Unit           string  `json:"unit"`
}

const (
	// TemplateTypeTemplates are templates that already contain a distribution.
	TemplateTypeTemplates string = "templates"
	// TemplateTypeFromScratch are templates that need to have distribution added to work.
	TemplateTypeFromScratch string = "from_scratch"
	pathPrefix              string = "/api/vsphere/v1/provisioning/cpu_performance_type.json"
)

func (a api) List(ctx context.Context) ([]CPUPerformanceType, error) {
	url := fmt.Sprintf("%s%s", a.client.BaseURL(), pathPrefix)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create cpu performance type list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute cpu performance type list request: %w", err)
	}
	var responsePayload []CPUPerformanceType
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode cpu performance type list response: %w", err)
	}

	return responsePayload, err
}
