package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/core/service"
)

var skipIntegration = true

func init() {
	var set bool
	if _, set = os.LookupEnv(client.IntegrationTestEnvName); !set {
		return
	}
	skipIntegration = false
}

func TestList(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	_, err = service.NewAPI(c).List(ctx, 1, 1000)
	cancel()
	if err != nil {
		t.Errorf("could not list services: %v", err)
	}
}
