package vm_test

import (
	"context"
	"testing"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/vm"
)

func TestGetIPs(t *testing.T) {
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = vm.GetFreeIPs(ctx, location, vlan, c)
	if err != nil {
		t.Fatalf("could not get free ips: %v", err)
	}
}
