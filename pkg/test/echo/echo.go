package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

const (
	// EchoPath can be used to test connectivity with the API.
	EchoPath = "/api/v1/test/echo.json"
)

// ErrInvalidEchoResponse indicates that an error request returned an invalid value.
var ErrInvalidEchoResponse = errors.New("invalid echo value received")

type echoRequest struct {
	Value string `json:"value"`
}

// Echo to test connectivity with the API.
func (a api) Echo(ctx context.Context) error {
	value := fmt.Sprintf("%v", rand.Int()) //nolint: gosec // No secure generator required.
	requestPayload := echoRequest{value}

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&requestPayload); err != nil {
		panic(fmt.Sprintf("could not encode definition: %v", err))
	}

	url := fmt.Sprintf("%s%s", a.client.BaseURL(), EchoPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &buf)
	if err != nil {
		panic(fmt.Sprintf("could not create echo request: %v", err))
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute echo request, got response %s", httpResponse.Status)
	}

	var responsePayload string
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return fmt.Errorf("could not decode echo response: %w", err)
	}

	if responsePayload != value {
		return fmt.Errorf("%w: expected %v , was %v", ErrInvalidEchoResponse, value, responsePayload)
	}

	return err
}
