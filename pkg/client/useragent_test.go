package client_test

import (
	"context"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/test/echo"
	"github.com/stretchr/testify/require"
	"net/http"
	"runtime"
	"testing"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()
	testUserAgent := "userAgentTest"
	c, err := client.New(client.TokenFromString("Test"), client.UserAgent(testUserAgent))
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAgent := r.Header.Get("User-Agent")
		require.Equal(t, testUserAgent, receivedAgent)
		echo.TestMock(t).ServeHTTP(w, r)
	})

	cw, server := client.NewTestClient(c, handler)

	ctx := context.Background()
	err = echo.NewAPI(cw).Echo(ctx)
	require.NoError(t, err)
	defer server.Close()
}

func TestWithNoAgent(t *testing.T) {
	t.Parallel()
	c, err := client.New(client.TokenFromString("Test"))
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAgent := r.Header.Get("User-Agent")
		require.Equal(t, fmt.Sprintf("go-anxcloud / %s (%s)", "snapshot", runtime.GOOS), receivedAgent)
		echo.TestMock(t).ServeHTTP(w, r)
	})

	cw, server := client.NewTestClient(c, handler)

	ctx := context.Background()
	err = echo.NewAPI(cw).Echo(ctx)
	require.NoError(t, err)
	defer server.Close()
}
