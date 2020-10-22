package vm

import (
	"context"
	"fmt"
	"net/http"
)

// Deprovision issues a request to deprovision an existing VM using.
//
// ctx is attached to the request and will cancel it on cancelation.
// It does not affect the provisioning request after it was issued.
// identifier is the VM identifier string returned when querying the
// provisioning task which ID was returned on VM provisioning.
// delayed indicated that the VM shall be removed with a delay of 24h.
//
// If the API returns errors, they are raised as ResponseError error.
func (a api) Deprovision(ctx context.Context, identifier string, delayed bool) error {
	url := fmt.Sprintf(
		"%s%s/%s?delayed=%t",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
		delayed,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create VM deprovisioning request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute VM deprovisioning request: %w", err)
	}
	_ = httpResponse.Body.Close()

	if err != nil {
		return fmt.Errorf("could not decode VM deprovisioning response: %w", err)
	}

	return err
}
