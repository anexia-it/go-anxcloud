// Package client provides a client for interacting with the anxcloud API.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"time"

	"github.com/go-logr/logr"
)

const (
	// TokenEnvName is the name of the environment variable that should contain the API token.
	TokenEnvName = "ANEXIA_TOKEN" //nolint:gosec // This is a name, not a secret.
	// VsphereLocationEnvName is the name of the environment variable that should contain a test location for paths that need a provisioning location.
	VsphereLocationEnvName = "ANEXIA_VSPHERE_LOCATION_ID"
	// CoreLocationEnvName is the name of the environment variable that should contain a test location for paths that need a core location.
	CoreLocationEnvName = "ANEXIA_CORE_LOCATION_ID"
	// VLANEnvName is the name of the environment variable that should contain the VLAN of VMs to manage.
	VLANEnvName = "ANEXIA_VLAN_ID"
	// IntegrationTestEnvName is the name of the environment variable that enables integration tests if present.
	IntegrationTestEnvName = "ANEXIA_INTEGRATION_TESTS_ON"
	// DefaultRequestTimeout is a suggested timeout for API calls.
	DefaultRequestTimeout = 10 * time.Second

	// defaultBaseURL is the default base URL used for requests.
	defaultBaseURL = "https://engine.anexia-it.com"
)

// ErrEnvMissing indicates an environment variable is missing.
var ErrEnvMissing = errors.New("environment variable missing")

var (
	// Version gets set by linker at build time
	version = "snapshot"
)

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

type client struct {
	httpClient        *http.Client
	token             string
	logger            logr.Logger
	userAgent         string
	baseURL           string
	parseEngineErrors bool
}

type clientOptions struct {
	client
	ignoreMissingToken bool
}

// Option is a optional parameter for the New method.
type Option func(o *clientOptions) error

// AuthFromEnv uses any known environment variables to create a client.
func AuthFromEnv(unset bool) Option {
	return TokenFromEnv(unset)
}

// TokenFromString uses the given API auth token.
func TokenFromString(token string) Option {
	return func(o *clientOptions) error {
		o.token = token

		return nil
	}
}

// UserAgent configures the user agent string to send with every HTTP request.
func UserAgent(userAgent string) Option {
	return func(o *clientOptions) error {
		o.userAgent = userAgent
		return nil
	}
}

// TokenFromEnv fetches the API auth token from environment variables.
func TokenFromEnv(unset bool) Option {
	return func(o *clientOptions) error {
		token, tokenPresent := os.LookupEnv(TokenEnvName)
		if !tokenPresent {
			return fmt.Errorf("%w: %s", ErrEnvMissing, TokenEnvName)
		}
		o.token = token
		if unset {
			if err := os.Unsetenv(TokenEnvName); err != nil {
				return fmt.Errorf("could not unset %s: %w", TokenEnvName, err)
			}
		}

		return nil
	}
}

// LogWriter configures the debug writer for logging requests and responses. Deprecated, use Logger instead.
func LogWriter(w io.Writer) Option {
	return func(o *clientOptions) error {
		o.logger = ioLogger(w)

		o.logger.Info("The LogWriter option of github.com/anexia-it/go-anxcloud/pkg/client is deprecated.")

		return nil
	}
}

// Logger configures where the client logs. Requests and responses are logged with verbosity LogVerbosityRequests
// on a logger derived from the one passed here with name LogNameTrace, replacing the previous LogWriter option.
func Logger(l logr.Logger) Option {
	return func(o *clientOptions) error {
		o.logger = l
		return nil
	}
}

// HTTPClient lets the client use the given http.Client.
func HTTPClient(c *http.Client) Option {
	return func(o *clientOptions) error {
		o.httpClient = c

		return nil
	}
}

// BaseURL configures the base URL for the client to use. Defaults to the production engine, but changing this
// can be useful for testing.
func BaseURL(baseURL string) Option {
	return func(o *clientOptions) error {
		o.baseURL = baseURL
		return nil
	}
}

// ParseEngineErrors is an option to chose if the client is supposed to parse http error responses into go errors or not.
func ParseEngineErrors(parseEngineErrors bool) Option {
	return func(o *clientOptions) error {
		o.parseEngineErrors = parseEngineErrors
		return nil
	}
}

// WithClient can be used to use another underlying http.Client.
func WithClient(hc *http.Client) Option {
	return func(o *clientOptions) error {
		o.httpClient = hc
		return nil
	}
}

// IgnoreMissingToken makes New() not return an error when no token is supplied.
func IgnoreMissingToken() Option {
	return func(o *clientOptions) error {
		o.ignoreMissingToken = true
		return nil
	}
}

// ErrConfiguration is raised when the given configuration is insufficient or erroneous.
var ErrConfiguration = errors.New("could not configure client")

// New creates a new client with the given options.
//
// The options need to contain a method of authentication with the API. If you are
// unsure what to use pass AuthFromEnv.
func New(options ...Option) (Client, error) {
	co := clientOptions{
		client: defaultClient(),
	}

	for _, option := range options {
		if err := option(&co); err != nil {
			return nil, err
		}
	}

	if co.client.userAgent == "" {
		co.client.userAgent = fmt.Sprintf("go-anxcloud/%s (%s)", version, runtime.GOOS)
	}

	if co.client.token == "" && !co.ignoreMissingToken {
		return nil, fmt.Errorf("%w: token not set", ErrConfiguration)
	}

	return &co.client, nil
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

	if c != nil {
		// TODO(LittleFox94): is this used somewhere and what for?
		return c, server
	}

	ret, err := New(
		BaseURL(server.URL),
		IgnoreMissingToken(),
	)

	if err != nil {
		panic(fmt.Errorf("error creating test client: %w", err))
	}

	return ret, server
}

func (c client) BaseURL() string {
	return c.baseURL
}

func (c client) Do(req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token %v", c.token))
	}

	req.Header.Set("User-Agent", c.userAgent)
	return c.handleRequest(req)
}

func (c client) handleRequest(req *http.Request) (*http.Response, error) {
	logRequest(req, c.logger)

	response, err := c.httpClient.Do(req)

	// TODO: we should probably handle redirects here. The Engine might not use them im Responses right now, but
	// it's a common HTTP future and the Engine might use them in the future.

	if c.parseEngineErrors && err == nil && (response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices) {
		errResponse := ResponseError{Request: req, Response: response}
		if decodeErr := json.NewDecoder(response.Body).Decode(&errResponse); decodeErr != nil {
			return response, fmt.Errorf("could not decode error response: %w", decodeErr)
		}

		err = &errResponse
	}

	logResponse(response, c.logger)

	return response, err
}

func defaultClient() client {
	return client{
		parseEngineErrors: true,
		logger:            logr.Discard(),
		baseURL:           defaultBaseURL,
		httpClient:        http.DefaultClient,
	}
}
