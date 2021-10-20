// Package client provides a client for interacting with the anxcloud API.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func handleRequest(c *http.Client, req *http.Request, logger logr.Logger) (*http.Response, error) {
	logRequest(req, logger)

	response, err := c.Do(req)
	if err == nil && (response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices) {
		errResponse := ResponseError{Request: req, Response: response}
		if decodeErr := json.NewDecoder(response.Body).Decode(&errResponse); decodeErr != nil {
			return response, fmt.Errorf("could not decode error response: %w", decodeErr)
		}

		err = &errResponse
	}

	logResponse(response, logger)

	return response, err
}

type optionSet struct {
	httpClient *http.Client
	token      string
	logger     *logr.Logger
	userAgent  string
	baseURL    string
}

// Option is a optional parameter for the New method.
type Option func(o *optionSet) error

// AuthFromEnv uses any known environment variables to create a client.
func AuthFromEnv(unset bool) Option {
	return TokenFromEnv(unset)
}

// TokenFromString uses the given API auth token.
func TokenFromString(token string) Option {
	return func(o *optionSet) error {
		o.token = token

		return nil
	}
}

// UserAgent configures the user agent string to send with every HTTP request.
func UserAgent(userAgent string) Option {
	return func(o *optionSet) error {
		o.userAgent = userAgent
		return nil
	}
}

// TokenFromEnv fetches the API auth token from environment variables.
func TokenFromEnv(unset bool) Option {
	return func(o *optionSet) error {
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
	return func(o *optionSet) error {
		l := ioLogger(w)
		o.logger = &l

		l.Info("The LogWriter option of github.com/anexia-it/go-anxcloud/pkg/client is deprecated.")

		return nil
	}
}

// Logger configures where the client logs. Requests and responses are logged with verbosity LogVerbosityRequests
// on a logger derived from the one passed here with name LogNameTrace, replacing the previous LogWriter option.
func Logger(l logr.Logger) Option {
	return func(o *optionSet) error {
		o.logger = &l
		return nil
	}
}

// HTTPClient lets the client use the given http.Client.
func HTTPClient(c *http.Client) Option {
	return func(o *optionSet) error {
		o.httpClient = c

		return nil
	}
}

func BaseURL(baseURL string) Option {
	return func(o *optionSet) error {
		o.baseURL = baseURL
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
	optionSet := optionSet{}
	for _, option := range options {
		if err := option(&optionSet); err != nil {
			return nil, err
		}
	}
	if optionSet.httpClient == nil {
		optionSet.httpClient = http.DefaultClient
	}

	if optionSet.userAgent == "" {
		optionSet.userAgent = fmt.Sprintf("go-anxcloud/%s (%s)", version, runtime.GOOS)
	}

	if optionSet.logger == nil {
		logger := logr.Discard()
		optionSet.logger = &logger
	}

	if optionSet.baseURL == "" {
		optionSet.baseURL = defaultBaseURL
	}

	if optionSet.token != "" {
		return &tokenClient{
			token:      optionSet.token,
			httpClient: optionSet.httpClient,
			logger:     *optionSet.logger,
			userAgent:  optionSet.userAgent,
			baseURL:    optionSet.baseURL,
		}, nil
	}

	return nil, fmt.Errorf("%w: token not set", ErrConfiguration)
}
