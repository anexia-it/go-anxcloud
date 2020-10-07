package vm_test

import (
	"context"
	"errors"
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

	progressID := provisionResponse.Identifier
	progressResponse := vm.ProgressResponse{}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	isDone := false
	var responseError *client.ResponseError
	for !isDone {
		select {
		case <-ticker.C:
			progressResponse, err := vm.GetProvisioningProgress(ctx, progressID, c)
			isProvisioningError := errors.As(err, &responseError)
			switch {
			case isProvisioningError && responseError.Response.StatusCode == 404:
			case err == nil:
				if progressResponse.Progress == 100 {
					isDone = true
				}
			default:
				t.Fatalf(fmt.Sprintf("could not query provision progress: %v", err))
			}
		case <-ctx.Done():
			t.Fatalf("vm did not get ready in time: %+v, %+v, %+v", progressResponse, *responseError, progressID)
		}
	}

	if err = vm.DeprovisionVM(ctx, progressResponse.VMIdentifier, false, c); err == nil {
		t.Fatalf(fmt.Sprintf("could not deprovision progress: %v", err))
	}
}
