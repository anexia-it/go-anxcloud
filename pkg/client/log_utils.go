package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	logfuncr "github.com/go-logr/logr/funcr"
)

const (
	// LogVerbosityRequests is the verbosity level for requests and responses.
	LogVerbosityRequests = 3

	// LogNameTrace is the name of the logger used for logging requests and responses.
	LogNameTrace = "trace"
)

func ioLogger(w io.Writer) logr.Logger {
	l := logfuncr.New(
		func(prefix, args string) {
			fmt.Fprintf(w, "%s\t%s\n", prefix, args)
		},
		logfuncr.Options{
			LogTimestamp: true,

			// we are logging request and response with verbosity 3 and since this is the only
			// thing we log and this option is for compatibility with previous versions, we set
			// this verbosity here.
			Verbosity: LogVerbosityRequests,
		},
	)

	return l
}

func traceLogger(log logr.Logger) logr.Logger {
	return log.WithName(LogNameTrace).V(LogVerbosityRequests)
}

func logRequest(req *http.Request, logger logr.Logger) {
	log := traceLogger(logger)

	if req == nil || !log.Enabled() {
		return
	}

	headers := req.Header.Clone()
	headers.Set("Authorization", "REDACTED")

	if body := stringifyBody(&req.Body, logger); body != nil {
		log = log.WithValues("body", body)
	}

	log.Info("Sending request to Engine",
		"url", req.URL.String(),
		"headers", stringifyHeaders(headers),
		"method", req.Method,
	)
}

func logResponse(res *http.Response, logger logr.Logger) {
	log := traceLogger(logger)

	if res == nil || !log.Enabled() {
		return
	}

	headers := res.Header.Clone()
	headers.Set("Set-Cookie", "REDACTED")

	if body := stringifyBody(&res.Body, logger); body != nil {
		log = log.WithValues("body", body)
	}

	if res.Request != nil {
		log = log.WithValues(
			"url", res.Request.URL.String(),
			"method", res.Request.Method,
		)
	}

	log.Info("Received response from Engine",
		"headers", stringifyHeaders(headers),
		"statusCode", res.StatusCode,
	)
}

func stringifyBody(body *io.ReadCloser, log logr.Logger) *string {
	if body != nil && *body != nil {
		b, err := io.ReadAll(*body)
		if err != nil {
			log.Error(err, "Error while preparing to log body, not logging it")
		} else {
			*body = io.NopCloser(bytes.NewBuffer(b))

			ret := string(b)
			return &ret
		}
	}

	return nil
}

func stringifyHeaders(headers http.Header) string {
	entries := make([]string, 0, len(headers))

	for header, values := range headers {
		// we sort the values to have an easier time reading the produced logs
		// ... also makes testing easier
		sort.Strings(values)

		// sometimes it's interesting to see what is _not_ in go's standard library ...
		// can I have some functional go pls? perhaps a slices.Map(array, func(entry))?
		quoted := make([]string, len(values))
		for i, val := range values {
			quoted[i] = fmt.Sprintf("'%s'", val)
		}

		var arrayEnclOpen, arrayEnclClose string

		if len(values) > 1 {
			arrayEnclOpen = "["
			arrayEnclClose = "]"
		}

		entries = append(entries, fmt.Sprintf("%s: %v%s%v", header, arrayEnclOpen, strings.Join(quoted, ", "), arrayEnclClose))
	}

	// we sort our headers to have an easier time reading the produced logs
	// ... also helps with testing
	sort.Strings(entries)

	return strings.Join(entries, ", ")
}
