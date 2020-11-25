package address_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
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
	_, err = address.NewAPI(c).List(ctx, 1, 1000)
	if err != nil {
		t.Fatalf("could not list address: %v", err)
	}
	cancel()
}
