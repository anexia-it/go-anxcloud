package echo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEcho(t *testing.T) {
	c, server := client.NewTestClient(nil, TestMock(t))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := NewAPI(c).Echo(ctx); err != nil {
		t.Errorf("echo request failed: %v", err)
	}
}

func TestEchoInvalidStatusCode(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := client.ResponseError{}
		resp.ErrorData.Code = http.StatusInternalServerError
		resp.ErrorData.Message = "testerror"
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			t.Errorf("echo response could not be encoded: %v", err)
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	var responseError *client.ResponseError
	if err := NewAPI(c).Echo(ctx); !errors.As(err, &responseError) {
		t.Errorf("expected client.ResponseError but got %v", err)
	}
}

func TestEchoInvalidResponseEncoding(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("echo payload could not be decoded: %v", err)
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
	if err := NewAPI(c).Echo(ctx); !errors.As(err, &syntaxError) {
		t.Errorf("expected json.SyntaxError but got %v", err)
	}
}

func TestEchoOtherValue(t *testing.T) {
	c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "\"%s\"\n", "not the right value")
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := NewAPI(c).Echo(ctx); !errors.Is(err, ErrInvalidEchoResponse) {
		t.Errorf("expected ErrInvalidEchoResponse but got %v", err)
	}
}

func TestEchoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Echo test suite")
}
