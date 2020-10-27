// Package client provides a client for interacting with the anxcloud API.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	// KeySecretEnvName is the name of the environment variable the signature key secret.
	KeySecretEnvName = "ANXCLOUD_KEY_SECRET" //nolint:gosec // This is a name, not a secret.
	// KeyIDEnvName is the name of the environment variable the signature key ID.
	KeyIDEnvName = "ANXCLOUD_KEY_ID"
	// TokenEnvName is the name of the environment variable that should contain the API token.
	TokenEnvName = "ANXCLOUD_TOKEN" //nolint:gosec // This is a name, not a secret.
	// LocationEnvName is the name of the environment variable that should contain the location of VMs to manage.
	LocationEnvName = "ANXCLOUD_LOCATION_ID"
	// VLANEnvName is the name of the environment variable that should contain the VLAN of VMs to manage.
	VLANEnvName = "ANXCLOUD_VLAN_ID"
	// IntegrationTestEnvName is the name of the environment variable that enables integration tests if present.
	IntegrationTestEnvName = "ANXCLOUD_INTEGRATION_TESTS_ON"
	// DefaultBaseURL is the default base URL used for requests.
	DefaultBaseURL = "https://engine.anexia-it.com"
	// EchoPath can be used to test connectivity with the API.
	EchoPath = "/api/v1/test/echo.json"
	// DefaultRequestTimeout is a suggested timeout for API calls.
	DefaultRequestTimeout = 10 * time.Second
)

// ErrEnvMissing indicates an environment variable is missing.
var ErrEnvMissing = errors.New("environment variable missing")

// ErrInvalidEchoResponse indicates that an error request returned an invalid value.
var ErrInvalidEchoResponse = errors.New("invalid echo value received")

// Client interacts with the anxcloud API.
type Client interface {
	// Do fires a given http.Request against the API.
	// This method behaves as http.Client.Do, but signs the request prior to sending it out
	// and returns an error is the response status is not OK.
	Do(req *http.Request) (*http.Response, error)
	BaseURL() string
}

// ResponseError is a response from the API that indicates an error.
type ResponseError struct {
	Request   *http.Request  `json:"-"`
	Response  *http.Response `json:"-"`
	ErrorData struct {
		Code       int               `json:"code"`
		Message    string            `json:"message"`
		Validation map[string]string `json:"validation"`
	} `json:"error"`
	Debug struct {
		Source string `json:"source"`
	} `json:"debug"`
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("received error from api: %+v", r.ErrorData)
}

// NewAnyClientFromEnvs tries to create a client from the present environment variables.
//
// unset can be set to true to let this method unset used environment variables after the client is
// successfully created.
// httpClient is the http.Client used for HTTP requests. Set the nil to use http.DefaultClient.
func NewAnyClientFromEnvs(unset bool, httpClient *http.Client) (Client, error) {
	_, tokenSet := os.LookupEnv(TokenEnvName)
	_, keyIDSet := os.LookupEnv(KeyIDEnvName)
	_, keySecretSet := os.LookupEnv(KeySecretEnvName)

	switch {
	case keyIDSet && keySecretSet:
		return NewSigningClientFromEnvs(unset, httpClient)
	case tokenSet:
		return NewTokenClientFromEnvs(unset, httpClient)
	default:
		return nil, fmt.Errorf("%w: either %s and %s must be set or %s", ErrEnvMissing, KeyIDEnvName, KeySecretEnvName, TokenEnvName)
	}
}

func handleRequest(c *http.Client, req *http.Request) (*http.Response, error) {
	response, err := c.Do(req)

	if err == nil && response.StatusCode != http.StatusOK {
		errResponse := ResponseError{Request: req, Response: response}
		if decodeErr := json.NewDecoder(response.Body).Decode(&errResponse); decodeErr != nil {
			return response, fmt.Errorf("could not decode error response: %w", decodeErr)
		}

		return response, &errResponse
	}

	return response, err
}
