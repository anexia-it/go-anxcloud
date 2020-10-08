package vm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/client"
	"anxkube-gitlab-dev.se.anx.io/anxkube/go-anxcloud/pkg/vm"
)

func TestVMProvisioningDeprovisioning(t *testing.T) {
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	ips, err := vm.GetFreeIPs(ctx, location, vlan, c)
	defer cancel()
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}
	if len(ips) < 1 {
		t.Fatalf("no IPs left for testing in vlan")
	}

	networkInterfaces := []vm.Network{{
		NICType: "vmxnet3",
		IPs:     []string{ips[0].Identifier},
		VLAN:    vlan,
	}}

	definition := vm.NewDefinition(location, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
	definition.SSH = randomPublicSSHKey()

	provisionResponse, err := vm.ProvisionVM(ctx, definition, c)
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}

	vmID, err := vm.AwaitProvisioning(ctx, provisionResponse.Identifier, c)
	if err != nil {
		t.Fatalf("waiting for VM provisioning failed: %v", err)
	}

	if err = vm.DeprovisionVM(ctx, vmID, false, c); err != nil {
		t.Fatalf(fmt.Sprintf("could not deprovision VM: %v", err))
	}
}
