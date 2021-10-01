// Package templates implements API functions residing under /provisioning/templates.
// This path contains methods for querying available VM templates.
package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// StringParameter is a string parameter for a template.
type StringParameter struct {
	Required bool   `json:"required"`
	Label    string `json:"label"`
	Default  string `json:"defaultValue"`
}

// BoolParameter is a bool parameter for a template.
type BoolParameter struct {
	Required bool   `json:"required"`
	Label    string `json:"label"`
	Default  bool   `json:"defaultValue"`
}

// IntParameter is an int parameter for a template.
type IntParameter struct {
	Minimum  int    `json:"minValue"`
	Maximum  int    `json:"maxValue"`
	Required bool   `json:"required"`
	Label    string `json:"label"`
	Default  int    `json:"defaultValue"`
}

// NIC is a single NIC in NICParameter.
type NIC struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

// NICParameter is a network interface card parameter for a template.
type NICParameter struct {
	Required bool   `json:"required"`
	Label    string `json:"label"`
	Default  int    `json:"defaultValue"`
	NICs     []NIC  `json:"data"`
}

// Parameters set of a VM.
type Parameters struct {
	Hostname         StringParameter `json:"hostname"`
	CPUs             IntParameter    `json:"cpus"`
	MemoryMB         IntParameter    `json:"memoryMB"`
	DiskGB           IntParameter    `json:"diskGB"`
	DNS0             StringParameter `json:"dns0"`
	DNS1             StringParameter `json:"dns1"`
	DNS2             StringParameter `json:"dns2"`
	DNS3             StringParameter `json:"dns3"`
	NICs             NICParameter    `json:"nics"`
	VLAN             StringParameter `json:"vlan"`
	IPs              StringParameter `json:"ips"`
	BootDelaySeconds IntParameter    `json:"bootDelaySeconds"`
	EnterBIOSSetup   BoolParameter   `json:"enterBIOSSetup"`
	Password         StringParameter `json:"password"`
	User             StringParameter `json:"user"`
	DiskType         StringParameter `json:"disk_type"`
}

// TemplateType defines which type of template is selected.
type TemplateType string

// Template contains a summary about the state of a template.
type Template struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	WordSize   string     `json:"bit"`
	Build      string     `json:"build"`
	Parameters Parameters `json:"param"`
}

const (
	// TemplateTypeTemplates are templates that already contain a distribution.
	TemplateTypeTemplates string = "templates"
	// TemplateTypeFromScratch are templates that need to have distribution added to work.
	TemplateTypeFromScratch string = "from_scratch"
	pathPrefix              string = "/api/vsphere/v1/provisioning/templates.json"
)

func (a api) List(ctx context.Context, locationID string, templateType string, page, limit int) ([]Template, error) {
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
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute template list request, got response %s", httpResponse.Status)
	}

	var responsePayload []Template
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode template list response: %w", err)
	}

	return responsePayload, err
}
