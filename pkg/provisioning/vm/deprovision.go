package vm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// Deprovision issues a request to deprovision an existing VM using.
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// identifier is the VM identifier string returned when querying the
// provisioning task which ID was returned on VM provisioning.
// delayed indicated that the VM shall be removed with a delay of 24h.
// client is the HTTP to be used for the request.
//
// If the API returns errors, they are raised as ResponseError error.
func Deprovision(ctx context.Context, identifier string, delayed bool, c client.Client) error {
	url := fmt.Sprintf(
		"https://%s%s/%s?delayed=%t",
		client.DefaultHost,
		pathPrefix,
		identifier,
		delayed,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create VM deprovisioning request: %w", err)
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute VM deprovisioning request: %w", err)
	}
	_ = httpResponse.Body.Close()

	if err != nil {
		return fmt.Errorf("could not decode VM deprovisioning response: %w", err)
	}

	return err
}
