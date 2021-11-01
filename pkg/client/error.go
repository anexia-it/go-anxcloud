package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func parseEngineError(req *http.Request, res *http.Response) error {
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		errResponse := ResponseError{Request: req, Response: res}
		if decodeErr := json.NewDecoder(res.Body).Decode(&errResponse); decodeErr != nil {
			return fmt.Errorf("could not decode error response: %w", decodeErr)
		}

		return &errResponse
	}

	return nil
}
