package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// An IntegerRequest is a request to do something with an integer.
type IntegerRequest struct {
	Input int `json:"input"`
}

// DecodeIntegerRequest decodes an IntegerRequest assuming JSON encoding.
func DecodeIntegerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req IntegerRequest
	return req, json.NewDecoder(r.Body).Decode(&req)
}

// An IntegerResponse is an integer result of a successful calculation.
type IntegerResponse struct {
	// Response body
	Body struct {
		// The result of the calculation.
		Output int `json:"output"`
	} `json:"body"`
}

// DecodeIntegerResponse decodes integer response from the body, given the right
// status code.
func DecodeIntegerResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if code := r.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", code)
	}

	var ret IntegerResponse
	return ret, json.NewDecoder(r.Body).Decode(&ret)
}
