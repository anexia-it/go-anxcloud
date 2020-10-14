package ips_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/provisioning/ips"
)

var (
	location = ""
	vlan     = ""
)

func init() {
	var set bool
	if location, set = os.LookupEnv(client.LocationEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.LocationEnvName))
	}
	if vlan, set = os.LookupEnv(client.VLANEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VLANEnvName))
	}
}

func TestGetFree(t *testing.T) {
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = ips.GetFree(ctx, location, vlan, c)
	if err != nil {
		t.Fatalf("could not get free ips: %v", err)
	}
}
