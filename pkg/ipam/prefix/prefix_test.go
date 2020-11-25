package prefix_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/prefix"
)

var (
	location        = ""
	vlan            = ""
	skipIntegration = true
)

func init() {
	var set bool
	if _, set = os.LookupEnv(client.IntegrationTestEnvName); !set {
		return
	}
	skipIntegration = false
	if location, set = os.LookupEnv(client.CoreLocationEnvName); !set || location == "" {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.CoreLocationEnvName))
	}
	if vlan, set = os.LookupEnv(client.VLANEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VLANEnvName))
	}
}

func TestList(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	_, err = prefix.NewAPI(c).List(ctx, 1, 1000)
	if err != nil {
		t.Fatalf("could not list prefix: %v", err)
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
		t.Fatalf("could not create client: %v", err)
	}

	api := prefix.NewAPI(c)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	summary, err := api.Create(ctx, prefix.NewCreate(location, vlan, 4, prefix.TypePrivate, 24))
	if err != nil {
		t.Fatalf("could not create prefix: %v", err)
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		info, err := api.Get(ctx, summary.ID)
		if err != nil {
			t.Fatalf("could not get vlan: %v", err)
		}
		if info.Status == "Active" {
			break
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			t.Fatalf(ctx.Err().Error())
		}
	}

	_, err = api.Update(ctx, summary.ID, prefix.Update{CustomerDescription: "something else"})
	if err != nil {
		t.Fatalf("could not update vlan: %v", err)
	}

	err = api.Delete(ctx, summary.ID)
	if err != nil {
		t.Fatalf("could not delete vlan: %v", err)
	}
}
