// Package progress implements API functions residing under /provisioning/progress.
// This path contains methods for querying the status of VM provisioning tasks.
package progress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.anx.io/go-anxcloud/pkg/client"
)

const (
	pathPrefix            = "/api/vsphere/v1/provisioning/progress.json"
	pollInterval          = 5 * time.Second
	progressCompleteValue = 100
	statusFailed          = -1
	statusSuccess         = 1
	statusInProgress      = 2
	statusCancelled       = 3
)

// Progress contains information regarding the provisioning of a VM returned by the API .
type Progress struct {
	// TaskIdentifier is the identifier of the provisioning task.
	TaskIdentifier string `json:"identifier"`
	// Queued indicates that the task is waiting to be started.
	Queued bool `json:"queued"`
	// Progress is the provisioning progress in percent (queuing not included).
	Progress int `json:"progress"`
	// VMIdentifier of the new VM.
	VMIdentifier string `json:"vm_identifier"`
	// Errors encountered while provisioning.
	Errors []string `json:"errors"`
	// Status of the task.
	Status int `json:"status"`
}

// ErrProgress is raised if a poll request completes but the result contains errors.
var ErrProgress = errors.New("progress response contains errors")

// Get queries the current progress of the provisioning of a VM.
//
// ctx is attached to the request and will cancel it on cancelation.
// identifier is the ID of the provisioning task to query. This is returned when
// provisioning the VM.
//
// If the API call returns errors, they are raised as ErrProgress.
// The returned progress response is still valid in this case.
func (a api) Get(ctx context.Context, identifier string) (Progress, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Progress{}, fmt.Errorf("could not create progress get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Progress{}, fmt.Errorf("could not execute progress get request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Progress{}, fmt.Errorf("could not execute progress get request, got response %s", httpResponse.Status)
	}

	var responsePayload Progress
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return Progress{}, fmt.Errorf("could not decode progress get response: %w", err)
	}

	switch {
	case len(responsePayload.Errors) == 0:
	case len(responsePayload.Errors) == 1 && responsePayload.Errors[0] == "The attempted operation cannot be performed in the current state (Powered on).":
	default:
		err = fmt.Errorf("%w: %v", ErrProgress, responsePayload.Errors)
	}

	return responsePayload, err
}

// AwaitCompletion polls the status of a started provisioning request and blocks until it is done.
//
// ctx will be checked for cancellation and the method returns immidiatly if so.
// progressID identifies the running provisioning task and is contained within ProvisioningResponse.
//
// Returned will be the VM ID and an error if polling or ProvisioningError if provisioning failed.
func (a api) AwaitCompletion(ctx context.Context, progressID string) (string, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	var responseError *client.ResponseError
	for {
		select {
		case <-ticker.C:
			progressResponse, err := a.Get(ctx, progressID)
			isProvisioningError := errors.As(err, &responseError)
			switch {
			case isProvisioningError && responseError.Response.StatusCode == 404:
				return "", fmt.Errorf("could not get progress. Endpoint returned 404: %w", err)
			case err == nil:
				if progressResponse.Progress == progressCompleteValue {
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
