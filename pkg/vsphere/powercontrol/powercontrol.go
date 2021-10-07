// Package powercontrol implements API functions residing under /powercontrol.
// This path contains methods for querying and setting the power state of VMs.
package powercontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// Request is the requested power state operation of a VM.
type Request string

// State is the current power state of a VM.
type State string

// Task inside the api to change the power state of a VM.
type Task struct {
	Progress       int    `json:"progress"`
	VMIdentifier   string `json:"identifier"`
	TaskIdentifier string `json:"task_id"`
	Error          string `json:"error"`
}

const (
	progressCompleteValue = 100
	pollInterval          = 5 * time.Second
	pathPrefix            = "/api/vsphere/v1/powercontrol.json"
)

var (
	// ErrSet is raised if the power state of a VM could not be set.
	ErrSet = errors.New("could not set powerstate of VM")
	// ErrInvalidState is raised if the API retuned an unknown power state.
	ErrInvalidState = errors.New("invalid power state received")

	// OnRequest indicates that the VM shall be on.
	OnRequest Request = "on"
	// HardRebootRequest indicates that the VM shall be rebooted without involving the OS. This is currently broken in API.
	HardRebootRequest Request = "hard_reboot"
	// HardShutdownRequest indicates that the VM shall shut down without involving the OS. This is currently broken in API.
	HardShutdownRequest Request = "hard_shutdown"
	// RebootRequest indicates that the VM shall be rebooted.
	RebootRequest Request = "reboot"
	// ShutdownRequest indicates that the VM shall be shut down.
	ShutdownRequest Request = "shutdown"

	// OnState indicates that the VM is on.
	OnState State = "VM_POWER_STATE_POWERED_ON"
	// OffState indicates that the VM is off.
	OffState State = "VM_POWER_STATE_POWERED_OFF"
)

// Do not use this, its broken.
func (a api) Set(ctx context.Context, identifier string, request Request) (Task, error) {
	url := fmt.Sprintf(
		"%s%s/%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
		request,
	)

	buf := &bytes.Buffer{}
	buf.WriteString("{}")

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, buf)
	if err != nil {
		return Task{}, fmt.Errorf("could not create powercontrol set request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Task{}, fmt.Errorf("could not execute powercontrol set request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Task{}, fmt.Errorf("could not execute powercontrol set request, got response %s", httpResponse.Status)
	}

	var task Task
	err = json.NewDecoder(httpResponse.Body).Decode(&task)
	if err != nil {
		return Task{}, fmt.Errorf("could not decode powercontrol set response: %w", err)
	}
	if task.Error != "" {
		return task, fmt.Errorf("%w: %s", ErrSet, task.Error)
	}

	return task, nil
}

// Get returns the power state of a given VM.
//
// ctx is attached to the request and will cancel it on cancelation.
// identifier is the ID of the VM to query.
func (a api) Get(ctx context.Context, identifier string) (State, error) {
	url := fmt.Sprintf(
		"%s%s/%s/info",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("could not create powercontrol get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute powercontrol get request: %w", err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return "", fmt.Errorf("could not execute powercontrol get request, got response %s", httpResponse.Status)
	}

	var responsePayload State
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	if err != nil {
		return "", fmt.Errorf("could not decode powercontrol get response: %w", err)
	}

	if responsePayload != OnState && responsePayload != OffState {
		return "", fmt.Errorf("%w: %s", ErrInvalidState, responsePayload)
	}

	return responsePayload, err
}

func (a api) AwaitCompletion(ctx context.Context, vmID, taskID string) error {
	url := fmt.Sprintf(
		"%s%s/%s/tasks/%s/info",
		a.client.BaseURL(),
		pathPrefix,
		vmID, taskID,
	)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("could not create powercontrol task get request: %w", err)
			}

			httpResponse, err := a.client.Do(req)
			if err != nil {
				return fmt.Errorf("could not execute powercontrol task get request: %w", err)
			}
			if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
				return fmt.Errorf("could not execute powercontrol task get request, got response %s", httpResponse.Status)
			}

			var responseError *client.ResponseError
			if errors.As(err, &responseError) && responseError.Response.StatusCode == 404 {
				continue
			}

			var task Task
			err = json.NewDecoder(httpResponse.Body).Decode(&task)
			_ = httpResponse.Body.Close()
			if err != nil {
				return fmt.Errorf("could not decode powercontrol task get response: %w", err)
			}

			if task.Progress == progressCompleteValue {
				return nil
			}

		case <-ctx.Done():
			return fmt.Errorf("powercontrol task did not complete in time: %w", ctx.Err())
		}
	}
}
