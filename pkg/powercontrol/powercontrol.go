// Package powercontrol implements API functions residing under /powercontrol.
// This path contains methods for querying and setting the power state of VMs.
package powercontrol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

type (
	// Request is the requested power state operation of a VM.
	Request string
	// State is the current power state of a VM.
	State string
)

const (
	pathPrefix = "/api/vsphere/v1/powercontrol.json"
)

var (
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

// Set issues a request change the power state of a given VM..
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// identifier is the ID of the VM to change.
// request is the desired operation to perform.
// c is the HTTP to be used for the request.
func Set(ctx context.Context, identifier string, request Request, c client.Client) error {
	url := fmt.Sprintf(
		"%s%s/%s/%s",
		c.BaseURL(),
		pathPrefix,
		identifier,
		request,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("could not create powercontrol set request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute powercontrol set request: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	return nil
}

// Get returns the power state of a given VM.
//
// ctx is attached to the request and will cancel it on cancelation.
// identifier is the ID of the VM to query.
// client is the HTTP to be used for the request.
func Get(ctx context.Context, identifier string, c client.Client) (State, error) {
	url := fmt.Sprintf(
		"%s%s/%s/info",
		c.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("could not create powercontrol get request: %w", err)
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute powercontrol get request: %w", err)
	}
	var responsePayload State
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return "", fmt.Errorf("could not decode powercontrol get response: %w", err)
	}

	if responsePayload != OnState && responsePayload != OffState {
		return "", fmt.Errorf("%w: %s", ErrInvalidState, responsePayload)
	}

	return responsePayload, err
}
