package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
)

func Test_traceLogging(t *testing.T) {
	logger := funcr.New(
		func(prefix, args string) {
			message := fmt.Sprintf("%s\t%s\n", prefix, args)

			if strings.Contains(message, "Authorization") {
				assert.Contains(t, message, "REDACTED")
				assert.NotContains(t, message, "auth_token")
			}

			if strings.Contains(message, "Set-Cookie") {
				assert.Contains(t, message, "REDACTED")
				assert.NotContains(t, message, "session_id")
			}
		},
		funcr.Options{
			Verbosity: 3,
		},
	)

	t.Run("Check Authorization redacted in request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo", nil)
		req.Header.Add("Authorization", "auth_token")

		logRequest(req, logger)
	})

	t.Run("Check Set-Cookie header redacted in response", func(t *testing.T) {
		w := httptest.NewRecorder()

		cookie := http.Cookie{
			Name:  "Session-Cookie",
			Value: "session_id",
		}

		http.SetCookie(w, &cookie)
		written, err := w.Write([]byte("OK"))
		assert.Equal(t, 2, written)
		assert.NoError(t, err)

		logResponse(w.Result(), logger)
	})

	t.Run("Check request body intact after logging", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo", bytes.NewBuffer([]byte("OK")))

		logRequest(req, logger)

		if body, err := ioutil.ReadAll(req.Body); err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, []byte("OK"), body)
		}
	})

	t.Run("Check response body intact after logging", func(t *testing.T) {
		w := httptest.NewRecorder()

		written, err := w.Write([]byte("OK"))
		assert.Equal(t, 2, written)
		assert.NoError(t, err)

		res := w.Result()

		logResponse(res, logger)

		if body, err := ioutil.ReadAll(res.Body); err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, []byte("OK"), body)
		}
	})

	t.Run("Check if logger is disabled on lower verbosity", func(t *testing.T) {
		logger := funcr.New(
			func(prefix, args string) {
				assert.Fail(t, "We shouldn't have logged anything")
			},
			funcr.Options{
				Verbosity: 0,
			},
		)

		assert.False(t, logger.V(3).Enabled())

		req := httptest.NewRequest("GET", "/foo", bytes.NewBuffer([]byte("OK")))
		logRequest(req, logger)
	})

	t.Run("Check not crashing on nil requests", func(t *testing.T) {
		logRequest(nil, logger)
	})

	t.Run("Check not crashing on nil responses", func(t *testing.T) {
		logResponse(nil, logger)
	})

	t.Run("Test Header stringification", func(t *testing.T) {
		headers := make(http.Header, 2)
		headers.Add("Foxes", "are cool")

		// this should be sorted to the start
		headers.Add("000", "1234")

		// add values for Foo in wrongly sorted order to test values are sorted
		headers.Add("Foo", "Baz")
		headers.Add("Foo", "Bar")

		headerString := stringifyHeaders(headers)

		assert.Contains(t, headerString, "Foo: ['Bar', 'Baz']")
		assert.Contains(t, headerString, "Foxes: 'are cool'")

		assert.Equal(t, "000: '1234', Foo: ['Bar', 'Baz'], Foxes: 'are cool'", headerString)
	})
}
