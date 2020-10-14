package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type tokenClient struct {
	token      string
	httpClient *http.Client
}

func (t tokenClient) BaseURL() string {
	return DefaultBaseURL
}

func (t tokenClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %v", t.token))

	response, err := t.httpClient.Do(req)

	if err == nil && response.StatusCode != http.StatusOK {
		errResponse := ResponseError{Request: req, Response: response}
		if decodeErr := json.NewDecoder(response.Body).Decode(&errResponse); decodeErr != nil {
			return response, fmt.Errorf("could not decode error response: %w. Original error was: %v", decodeErr, err)
		}

		return response, &errResponse
	}

	return response, err
}

// NewTokenClient creates a new token client for the anxcloud that uses tokens.
//
// token is the token you received from the webinterface.
// httpClient is the http.Client used for HTTP requests. Set the nil to use http.DefaultClient.
func NewTokenClient(token string, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &tokenClient{token, httpClient}
}

// NewTokenClientFromEnvs extracts token settings from environment variables and uses
// NewTokenClient to create a client.
//
// unset can be set to true to let this method unset used environment variables after the client is
// successfully created.
// httpClient is the http.Client used for HTTP requests. Set the nil to use http.DefaultClient.
func NewTokenClientFromEnvs(unset bool, httpClient *http.Client) (Client, error) {
	token, tokenPresent := os.LookupEnv(TokenEnvName)
	if !tokenPresent {
		return nil, fmt.Errorf("%w: %s", ErrEnvMissing, TokenEnvName)
	}

	client := NewTokenClient(token, httpClient)
	if unset {
		if err := os.Unsetenv(KeyIDEnvName); err != nil {
			return client, fmt.Errorf("could not unset %s: %w", KeyIDEnvName, err)
		}
		if err := os.Unsetenv(KeySecretEnvName); err != nil {
			return client, fmt.Errorf("could not unset %s: %w", KeySecretEnvName, err)
		}
	}

	return client, nil
}
