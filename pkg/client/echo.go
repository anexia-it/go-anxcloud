package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type echoRequest struct {
	Value string `json:"value"`
}

// Echo to test connectivity with the API.
func Echo(ctx context.Context, c Client) error {
	value := fmt.Sprintf("%v", rand.Int()) //nolint: gosec // No secure generator required.
	requestPayload := echoRequest{value}

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&requestPayload); err != nil {
		panic(fmt.Sprintf("could not encode definition: %v", err))
	}

	url := fmt.Sprintf("%s%s", c.BaseURL(), EchoPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &buf)
	if err != nil {
		panic(fmt.Sprintf("could not create echo request: %v", err))
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return err
	}

	var responsePayload string
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return fmt.Errorf("could not decode echo response: %w", err)
	}

	if responsePayload != value {
		return fmt.Errorf("%w: expected %v , was %v", ErrInvalidEchoResponse, value, responsePayload)
	}

	return err
}
