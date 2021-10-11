package vm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a api) Update(ctx context.Context, identifier string, change Change) (ProvisioningResponse, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&change); err != nil {
		panic(fmt.Sprintf("could not encode update: %v", err))
	}

	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &buf)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not create VM update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not execute VM update request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return ProvisioningResponse{}, fmt.Errorf("could not execute VM update request, got response %s", httpResponse.Status)
	}

	var responsePayload ProvisioningResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not decode VM update response: %w", err)
	}

	if len(responsePayload.Errors) != 0 {
		err = fmt.Errorf("%w: %v", ErrProvisioning, responsePayload.Errors)
	}

	return responsePayload, err
}
