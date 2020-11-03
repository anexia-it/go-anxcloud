package ips_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/ips"
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
	if location, set = os.LookupEnv(client.LocationEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.LocationEnvName))
	}
	if vlan, set = os.LookupEnv(client.VLANEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VLANEnvName))
	}
}

func TestGetFree(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = ips.NewAPI(c).GetFree(ctx, location, vlan)
	if err != nil {
		t.Fatalf("could not get free ips: %v", err)
	}
}
