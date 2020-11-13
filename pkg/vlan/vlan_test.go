package vlan_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vlan"
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
	if location, set = os.LookupEnv(client.LocationEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.LocationEnvName))
	}
	skipIntegration = false
}

func TestAll(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	_, err = vlan.NewAPI(c).All(ctx)
	if err != nil {
		t.Fatalf("could not list VLANSt: %v", err)
	}
	cancel()
}

func TestCreateetDelete(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	api := vlan.NewAPI(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	summary, err := api.Create(ctx, vlan.CreateDefinition{Location: location})
	cancel()
	if err != nil {
		t.Fatalf("could not create vlan: %v", err)
	}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		info, err := api.Get(ctx, summary.Identifier)
		cancel()
		if err != nil {
			t.Fatalf("could not get vlan: %v", err)
		}
		if info.Status == "Active" {
			break
		}
		time.Sleep(3 * time.Second)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = api.Delete(ctx, summary.Identifier)
	cancel()
	if err != nil {
		t.Fatalf("could not delete vlan: %v", err)
	}
}
