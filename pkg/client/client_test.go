package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_handleRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			formValue := r.FormValue("foo")
			assert.EqualValues(t, "bar", formValue)
			assert.EqualValues(t, "sensible-value", r.Header.Get("Authorization"))

			_, err := io.WriteString(w, formValue)
			assert.NoError(t, err)
		})

		srv := httptest.NewServer(handler)
		defer srv.Close()

		data := url.Values{}
		data.Set("foo", "bar")
		encData := data.Encode()
		req, err := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(encData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(encData)))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "sensible-value")

		buffer := make([]byte, 0)
		writeBuffer := bytes.NewBuffer(buffer)

		response, err := handleRequest(http.DefaultClient, req, writeBuffer)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			body, err := ioutil.ReadAll(response.Body)
			assert.NoError(t, err)
			assert.EqualValues(t, "bar", body)
		}
		assert.True(t, strings.Contains(writeBuffer.String(), "Authorization: REDACTED"))
	})

	t.Run("Error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			formValue := r.FormValue("foo")
			assert.EqualValues(t, "bar", formValue)
			assert.EqualValues(t, "sensible-value", r.Header.Get("Authorization"))

			errorMsg := map[string]string{
				"msg": "error message",
			}
			encoder := json.NewEncoder(w)
			w.WriteHeader(http.StatusBadRequest)
			assert.NoError(t, encoder.Encode(errorMsg))
		})

		srv := httptest.NewServer(handler)
		defer srv.Close()

		data := url.Values{}
		data.Set("foo", "bar")
		encData := data.Encode()
		req, err := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(encData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(encData)))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "sensible-value")

		buffer := make([]byte, 0)
		writeBuffer := bytes.NewBuffer(buffer)

		response, err := handleRequest(http.DefaultClient, req, writeBuffer)
		assert.Error(t, err)
		if assert.NotNil(t, response) {
			assert.EqualValues(t, response.StatusCode, http.StatusBadRequest)
		}
		assert.True(t, strings.Contains(writeBuffer.String(), "Authorization: REDACTED"))
	})
}
