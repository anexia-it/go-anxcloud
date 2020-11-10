package echo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestMock creates a new http.Handler that emulated the echo API endpoint.
func TestMock(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != EchoPath {
			t.Fatalf("not using the correct echo path but: %s", r.URL.Path)
		}
		payload := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("echo payload could not be decoded: %v", err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "\"%s\"\n", payload["value"])
	})
}
