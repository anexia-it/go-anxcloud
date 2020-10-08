package vm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
)

const (
	progressPathPrefix                = "/api/vsphere/v1/provisioning/progress.json"
	provisionPollInterval             = 5 * time.Second
	provisioningProgressCompleteValue = 100
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

// AwaitProvisioning polls the status of a started provisioning request and blocks until it is done.
//
// ctx will be checked for cancellation and the method returns immidiatly if so.
// progressID identifies the running provisioning task and is contained within ProvisioningResponse.
// c is the HTTP to be used for polling requests.
//
// Returned will be the VM ID and an error if polling or ProvisioningError if provisioning failed.
func AwaitProvisioning(ctx context.Context, progressID string, c client.Client) (string, error) {
	ticker := time.NewTicker(provisionPollInterval)
	defer ticker.Stop()
	var responseError *client.ResponseError
	for {
		select {
		case <-ticker.C:
			progressResponse, err := GetProvisioningProgress(ctx, progressID, c)
			isProvisioningError := errors.As(err, &responseError)
			switch {
			case isProvisioningError && responseError.Response.StatusCode == 404:
			case err == nil:
				if progressResponse.Progress == provisioningProgressCompleteValue {
					return progressResponse.VMIdentifier, nil
				}
			default:
				return "", fmt.Errorf("could not query provision progress: %w", err)
			}
		case <-ctx.Done():
			return "", fmt.Errorf("vm did not get ready in time: %w", ctx.Err())
		}
	}
}
