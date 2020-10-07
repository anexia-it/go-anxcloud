package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
)

const (
	progressPathPrefix = "/api/vsphere/v1/provisioning/progress.json"
)

// ProgressResponse contains information regarding the provisioning of a VM returned by the API .
type ProgressResponse struct {
	// TaskIdentifier is the identifier of the provisioning task.
	TaskIdentifier string `json:"identifier"`
	// Queued indicated that the task is waiting to be started-
	Queued bool `json:"queued"`
	// Progress is the provisioning progress in percent (queuing not included).
	Progress int `json:"progress"`
	// VMIdentifier of the new VM.
	VMIdentifier string `json:"vm_identifier"`
	// Errors encountered while provisioning.
	Errors []string `json:"errors"`
}

// GetProvisioningProgress queries the current progress of the provisioning of a VM.
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// identifier is the ID of the provisioning task to query. This is returned when
// provisioning the VM.
// client is the HTTP to be used for the request.
//
// If the API returns errors, they are raised as ResponseError error.
func GetProvisioningProgress(ctx context.Context, identifier string, c client.Client) (ProgressResponse, error) {
	url := fmt.Sprintf(
		"https://%s%s/%s",
		client.DefaultHost,
		progressPathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("could not create progress get request: %w", err)
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("could not execute progress get request: %w", err)
	}
	var responsePayload ProgressResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return ProgressResponse{}, fmt.Errorf("could not decode progress get response: %w", err)
	}

	if len(responsePayload.Errors) != 0 {
		err = &ProvisioningError{responsePayload.Errors}
	}

	return responsePayload, err
}
