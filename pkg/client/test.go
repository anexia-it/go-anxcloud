package client

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-logr/logr"
)

type testClient struct {
	baseClient Client
	baseURL    string
	httpClient *http.Client
	logger     logr.Logger
}

func (t testClient) BaseURL() string {
	return t.baseURL
}

func (t testClient) Do(req *http.Request) (*http.Response, error) {
	if t.baseClient != nil {
		return t.baseClient.Do(req)
	}

	return handleRequest(t.httpClient, req, t.logger)
}

// NewTestClient creates a new client for testing.
//
// c may be used to specify an other client implementation that needs to be tested
// or may be nil.
// handler is a http.Handler that mocks parts of the API functionality that shall be tested.
//
// Returned will be a client.Client that can be passed to the method under test and the
// used httptest.Server that should be closed after test completion.
func NewTestClient(c Client, handler http.Handler) (Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	cw := testClient{c, server.URL, &http.Client{}, logr.Discard()}

	return cw, server
}
