package client_test

import (
	"context"
	"testing"
	"time"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
)

func TestEchoWithToken(t *testing.T) {
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := client.ExecuteEcho(ctx, c); err != nil {
		t.Fatalf("[%s] echo test failed: %v", time.Now(), err)
	}
}
