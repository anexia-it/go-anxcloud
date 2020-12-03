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
	if location, set = os.LookupEnv(client.VsphereLocationEnvName); !set || location == "" {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VsphereLocationEnvName))
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
		t.Errorf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = ips.NewAPI(c).GetFree(ctx, location, vlan)
	if err != nil {
		t.Errorf("could not get free ips: %v", err)
	}
}
