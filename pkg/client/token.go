package client

import (
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

type tokenClient struct {
	token      string
	httpClient *http.Client
	logger     logr.Logger
	userAgent  string
}

func (t tokenClient) BaseURL() string {
	return DefaultBaseURL
}

func (t tokenClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %v", t.token))
	req.Header.Set("User-Agent", t.userAgent)
	return handleRequest(t.httpClient, req, t.logger)
}
