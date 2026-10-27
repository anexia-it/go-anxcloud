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
	decodeErr error
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

// StatusCode returns the HTTP status code of the response.
func (r ResponseError) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode
	}

	return r.ErrorData.Code
}

// Unwrap returns the error while decoding the error response.
func (r ResponseError) Unwrap() error {
	return r.decodeErr
}

func parseEngineError(req *http.Request, res *http.Response) error {
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		errResponse := ResponseError{Request: req, Response: res}
		if decodeErr := json.NewDecoder(res.Body).Decode(&errResponse); decodeErr != nil {
			errResponse.decodeErr = decodeErr
			errResponse.ErrorData.Code = res.StatusCode
			errResponse.ErrorData.Message = res.Status
		} else if errResponse.ErrorData.Code == 0 {
			errResponse.ErrorData.Code = res.StatusCode
		}

		return &errResponse
	}

	return nil
}
