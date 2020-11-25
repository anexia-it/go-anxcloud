// Package service implements API functions residing under /core/service.
// This path contains methods for listing services.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/core/v1/service.json"
)

// Service describes a service provided by Anexia.
type Service struct {
	Name     string `json:"name"`
	ID       string `json:"identifier"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type listResponse struct {
	Data []Service `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int) ([]Service, error) {
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathPrefix, page, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create service list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute service list request: %w", err)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode service list response: %w", err)
	}

	return responsePayload.Data, nil
}
