package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	bodylessHeaders = `(request-target) host date`
	bodyHeaders     = `(request-target) host date content-type content-length digest`
)

type signingClient struct {
	keySecret  string
	keyID      string
	httpClient *http.Client
}

func (s signingClient) BaseURL() string {
	return DefaultBaseURL
}

func (s signingClient) Do(req *http.Request) (*http.Response, error) {
	headers := []string{fmt.Sprintf("(request-target): %s %s", strings.ToLower(req.Method), req.URL.RequestURI())}
	req.Header.Set("host", req.Host)
	headers = append(headers, fmt.Sprintf("host: %s", req.Header.Get("host")))
	req.Header.Set("date", time.Now().Format(time.RFC1123))
	headers = append(headers, fmt.Sprintf("date: %s", req.Header.Get("date")))

	headerSet := bodylessHeaders
	if req.Body != nil {
		bodyReader, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("could not extract request body for signing: %w", err)
		}
		buf := bytes.Buffer{}
		if _, err = io.Copy(&buf, bodyReader); err != nil {
			return nil, fmt.Errorf("could not extract request body for signing: %w", err)
		}
		body := buf.Bytes()
		sha := sha512.New()
		if _, err := sha.Write(body); err != nil {
			panic(fmt.Sprintf("could not write hash: %v", err))
		}

		req.Header.Set("content-type", "application/json")
		headers = append(headers, fmt.Sprintf("content-type: %s", req.Header.Get("content-type")))
		req.Header.Set("content-length", fmt.Sprintf("%v", len(body)))
		headers = append(headers, fmt.Sprintf("content-length: %s", req.Header.Get("content-length")))
		req.Header.Set("digest", fmt.Sprintf("SHA-512=%s", base64.StdEncoding.EncodeToString(sha.Sum(nil))))
		headers = append(headers, fmt.Sprintf("digest: %s", req.Header.Get("digest")))
		req.Header.Set("headers", bodyHeaders)
		headerSet = bodyHeaders
	}

	hash := hmac.New(sha512.New, []byte(s.keySecret))
	if _, err := hash.Write([]byte(strings.Join(headers, "\n"))); err != nil {
		panic(fmt.Sprintf("could not write hash: %v", err))
	}
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf(`Signature keyId="%s",algorithm="hmac-sha512",signature="%s",headers="%s"`, s.keyID, sig, headerSet))

	response, err := s.httpClient.Do(req)
	if err == nil && response.StatusCode != http.StatusOK {
		errResponse := ResponseError{Request: req, Response: response}
		if decodeErr := json.NewDecoder(response.Body).Decode(&errResponse); decodeErr != nil {
			return response, fmt.Errorf("could not decode error response: %w. Original error was: %v", decodeErr, err)
		}

		return response, &errResponse
	}

	return response, err
}

// NewSigningClient creates a new signing client for the anxcloud that uses HTTP Signature Authentication.
//
// keySecret and keyID are signature key and ID you can fetch from the anxcloud webinterface.
// httpClient is the http.Client used for HTTP requests. Set the nil to use http.DefaultClient.
func NewSigningClient(keySecret, keyID string, httpClient *http.Client) (Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &signingClient{keySecret, keyID, httpClient}, nil
}

// NewSigningClientFromEnvs extracts HTTP Signature Authentication settings from environment variables and uses
// NewSigningClient to create a client.
//
// unset can be set to true to let this method unset used environment variables after the client is
// successfully created.
// httpClient is the http.Client used for HTTP requests. Set the nil to use http.DefaultClient.
func NewSigningClientFromEnvs(unset bool, httpClient *http.Client) (Client, error) {
	id, idPresent := os.LookupEnv(KeyIDEnvName)
	if !idPresent {
		return nil, fmt.Errorf("%w: %s", ErrEnvMissing, KeyIDEnvName)
	}
	secret, secretPresent := os.LookupEnv(KeySecretEnvName)
	if !secretPresent {
		return nil, fmt.Errorf("%w: %s", ErrEnvMissing, KeySecretEnvName)
	}

	client, err := NewSigningClient(secret, id, httpClient)
	if err != nil {
		return client, err
	}
	if unset {
		if err = os.Unsetenv(KeyIDEnvName); err != nil {
			return client, fmt.Errorf("could not unset %s: %w", KeyIDEnvName, err)
		}
		if err = os.Unsetenv(KeySecretEnvName); err != nil {
			return client, fmt.Errorf("could not unset %s: %w", KeySecretEnvName, err)
		}
	}

	return client, nil
}
