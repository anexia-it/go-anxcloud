package vm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ProvisioningResponse contains information returned by the API regarding a newly created VM.
type ProvisioningResponse struct {
	Progress   int      `json:"progress"`
	Errors     []string `json:"errors"`
	Identifier string   `json:"identifier"`
	Queued     bool     `json:"queued"`
}

// ErrProvisioning is raised if the API returns an error.
var ErrProvisioning = errors.New("ProvisioningResponse contains errors")

// Provision issues a request to provision a new VM using the given VM definition.
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// definition contains the definition of the VM to be created.
//
// If the API call returns errors, they are raised as ErrProvisioning.
// The returned ProvisioningResponse is still valid in this case.
func (a api) Provision(ctx context.Context, definition Definition, scriptBase64Encoded bool) (ProvisioningResponse, error) {
	buf := bytes.Buffer{}

	if definition.Script != "" && scriptBase64Encoded {
		definition.Script = base64.StdEncoding.EncodeToString([]byte(definition.Script))
	}

	if err := json.NewEncoder(&buf).Encode(&definition); err != nil {
		panic(fmt.Sprintf("could not encode definition: %v", err))
	}

	url := fmt.Sprintf(
		"%s%s/%s/%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		definition.Location,
		definition.TemplateType,
		definition.TemplateID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not create VM provisioning request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not execute VM provisioning request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return ProvisioningResponse{}, fmt.Errorf("could not execute VM provisioning request, got response %s", httpResponse.Status)
	}

	var responsePayload ProvisioningResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return ProvisioningResponse{}, fmt.Errorf("could not decode VM provisioning response: %w", err)
	}

	if len(responsePayload.Errors) != 0 {
		err = fmt.Errorf("%w: %v", ErrProvisioning, responsePayload.Errors)
	}

	return responsePayload, err
}
