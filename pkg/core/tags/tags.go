// Package tags implements API functions residing under /core/tags.
// This path contains methods for querying, creating and deleting of tags.
package tags

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pathPrefix = "/api/core/v1/tags.json"
)

// Create contains information of a tag to create.
type Create struct {
	Name       string `json:"name"`
	ServiceID  string `json:"service_identifier"`
	CustomerID string `json:"customer_identifier"`
}

// Summary contains short info about a tag.
type Summary struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

// Service is part of Organisation.
type Service struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

// Customer is part of Organisation.
type Customer struct {
	CustomerID string `json:"customer_id"`
	Demo       bool   `json:"demo"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Slug       string `json:"name_slug"`
	Reseller   string `json:"reseller"`
}

// Organisation is part of info.
type Organisation struct {
	Customer Customer `json:"customer"`
	Service  Service  `json:"service"`
}

// Info contains all info about a tag.
type Info struct {
	Name          string         `json:"name"`
	Identifier    string         `json:"identifier"`
	Organisations []Organisation `json:"organisation_assignments"`
}

type listResponse struct {
	Data []Summary `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int, query, serviceIdentifier, organizationIdentifier, order string, sortAscending bool) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v&query=%s&service_identifier=%s&organization_identifier=%s&order=%s&sort_descending=%s",
		a.client.BaseURL(), pathPrefix,
		page, limit, query, serviceIdentifier, organizationIdentifier, order, strconv.FormatBool(sortAscending),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create tags list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute tags list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute tags list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode tags list response: %w", err)
	}

	return responsePayload.Data, nil
}

func (a api) Get(ctx context.Context, identifier string) (Info, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("could not create tags get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute tags get request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Info{}, fmt.Errorf("could not execute tags get request, got response %s", httpResponse.Status)
	}

	var info Info
	err = json.NewDecoder(httpResponse.Body).Decode(&info)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Info{}, fmt.Errorf("could not decode tags get response: %w", err)
	}

	return info, nil
}

func (a api) Create(ctx context.Context, create Create) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(create); err != nil {
		panic(fmt.Sprintf("could not create request data for tag creation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create tag post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute tag post request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute tag post request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode tag post response: %w", err)
	}

	return summary, nil
}

func (a api) Delete(ctx context.Context, tagID, serviceID string) error {
	url := fmt.Sprintf(
		"%s%s/%s?service_identifier=%v",
		a.client.BaseURL(),
		pathPrefix, tagID, serviceID,
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
