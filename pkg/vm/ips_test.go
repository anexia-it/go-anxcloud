package vm_test

import (
	"context"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vm"
)

func TestGetIPs(t *testing.T) {
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("[%s] could not create client: %v", time.Now(), err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = vm.GetFreeIPs(ctx, location, vlan, c)
	if err != nil {
		t.Fatalf("[%s] could not get free ips: %v", time.Now(), err)
	}
}
