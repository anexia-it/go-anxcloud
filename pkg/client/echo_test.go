package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

func echoTestHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != client.EchoPath {
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

func TestEcho(t *testing.T) {
	c, server := client.NewTestClient(nil, echoTestHandler(t))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := client.Echo(ctx, c); err != nil {
		t.Fatalf("echo request failed: %v", err)
	}
}

func TestEchoInvalidStatusCode(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := client.ResponseError{}
		resp.ErrorData.Code = http.StatusInternalServerError
		resp.ErrorData.Message = "testerror"
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			t.Fatalf("echo response could not be encoded: %v", err)
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	var responseError *client.ResponseError
	if err := client.Echo(ctx, c); !errors.As(err, &responseError) {
		t.Fatalf("expected client.ResponseError but got %v", err)
	}
}

func TestEchoInvalidResponseEncoding(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("echo payload could not be decoded: %v", err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "\"%s\n", payload["value"])
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	var syntaxError *json.SyntaxError
	if err := client.Echo(ctx, c); !errors.As(err, &syntaxError) {
		t.Fatalf("expected json.SyntaxError but got %v", err)
	}
}

func TestEchoOtherValue(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "\"%s\"\n", "not the right value")
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := client.Echo(ctx, c); !errors.Is(err, client.ErrInvalidEchoResponse) {
		t.Fatalf("expected ErrInvalidEchoResponse but got %v", err)
	}
}
