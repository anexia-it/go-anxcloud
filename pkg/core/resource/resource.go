// Package resource implements API functions residing under /core/resource.
// This path contains methods for querying resources and attaching tags to them.
package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const pathPrefix = "/api/core/v1/resource.json"

// Summary describes a resource in short.
type Summary struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// Type is part of info.
type Type struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// Info contains all information about a resource.
type Info struct {
	Identifier string   `json:"identifier" anxcloud:"identifier"`
	Name       string   `json:"name"`
	Type       Type     `json:"resource_type"`
	Tags       []string `json:"tags"`
}

type listResponse struct {
	Data []Summary `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathPrefix, page, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create resource list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute resource list request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute resource list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	if err != nil {
		return nil, fmt.Errorf("could not decode resource list response: %w", err)
	}

	return responsePayload.Data, nil
}

func (a api) Get(ctx context.Context, id string) (Info, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("could not create resource get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute resource get request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Info{}, fmt.Errorf("could not execute resource get request, got response %s", httpResponse.Status)
	}

	var info Info
	err = json.NewDecoder(httpResponse.Body).Decode(&info)
	if err != nil {
		return Info{}, fmt.Errorf("could not decode resource get response: %w", err)
	}

	return info, nil
}

func (a api) AttachTag(ctx context.Context, resourceID, tagName string) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s/%v/tags/%v",
		a.client.BaseURL(),
		pathPrefix, resourceID, tagName,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not attach tag post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute attach tag request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute attach tag request, got response %s", httpResponse.Status)
	}

	var summary []Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return nil, fmt.Errorf("could not decode attach tag response: %w", err)
	}

	return summary, nil
}

func (a api) DetachTag(ctx context.Context, resourceID, tagName string) error {
	url := fmt.Sprintf(
		"%s%s/%v/tags/%v",
		a.client.BaseURL(),
		pathPrefix, resourceID, tagName,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create tag delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute tag delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute tag delete request, got response %s", httpResponse.Status)
	}

	return httpResponse.Body.Close()
}
