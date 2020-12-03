package tags_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/core/tags"
)

var (
	skipIntegration = true
	location        = ""
)

func init() {
	var set bool
	if _, set = os.LookupEnv(client.IntegrationTestEnvName); !set {
		return
	}
	if location, set = os.LookupEnv(client.CoreLocationEnvName); !set || location == "" {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.CoreLocationEnvName))
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
	_, err = tags.NewAPI(c).List(ctx, 1, 1000)
	if err != nil {
		t.Errorf("could not list VLAN: %v", err)
	}
	cancel()
}

func TestCreateDelete(t *testing.T) {
	t.Parallel()
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}

	api := tags.NewAPI(c)

	serviceID := "ff543fc08b3149ee9a8c50ee018b15a6"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	summary, err := api.Create(ctx, tags.Create{
		Name:      "go sdk integration test",
		ServiceID: "ff543fc08b3149ee9a8c50ee018b15a6",
	})
	if err != nil {
		t.Errorf("could not create vlan: %v", err)
	}

	_, err = api.Get(ctx, summary.Identifier)
	if err != nil {
		t.Errorf("could not get vlan: %v", err)
	}

	err = api.Delete(ctx, summary.Identifier, serviceID)
	if err != nil {
		t.Errorf("could not delete vlan: %v", err)
	}
}
