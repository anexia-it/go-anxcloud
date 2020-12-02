[![Documentation](https://godoc.org/github.com/anexia-it/go-anxcloud?status.svg)](http://godoc.org/github.com/anexia-it/go-anxcloud)

# Go Client for the Anexia API

Go SDK for interacting with the Anexia multi purpose [API](https://engine.anexia-it.com/).

## Installing

To use the SDK, just add `github.com/anexia-it/go-anxcloud <version>` to your Go module.

## Getting started

Before using the SDK you should familiarize yourself with the API. See [here](https://engine.anexia-it.com/api/core/doc/) for more info.
I you crave for an example, take a look at [the terraform provider for this project](https://github.com/anexia-it/terraform-provider-anxcloud)

## Example 

The following code shows how to create a VM. To be able to do that you need to set the environment variable `ANEXIA_TOKEN` to your access token.
Afterwards you can run the following.

```go
package main

import (
	"context"
	"fmt"
	"time"

	anexia "github.com/anexia-it/go-anxcloud/pkg"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/vm"
)

func main() {
	vlan := "<ID of the VLAN the VM should have access to>"
	location := "<ID of the location the VM should be in>"

	// Create client from environment variables, do not unset env afterwards.
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		panic(fmt.Sprintf("could not create client: %v", err))
	}

	// Get some API.
	provisioning := anexia.NewAPI(c).VSphere().Provisioning()

	// Time out after 30 minutes. Yes it really takes that long sometimes.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	// Look for a free ip in the given VLAN. This IP is not reserved for you so better be quick.
	ips, err := provisioning.IPs().GetFree(ctx, location, vlan)
	defer cancel()
	if err != nil {
		panic(fmt.Sprintf("provisioning vm failed: %v", err))
	}
	if len(ips) < 1 {
		panic(fmt.Sprintf("no IPs left for testing in vlan"))
	}

	// Create a NIC for the VM and connect it to the VLAN.
	networkInterfaces := []vm.Network{{NICType: "vmxnet3", IPs: []string{ips[0].Identifier}, VLAN: vlan}}
	// Create the definition of the new VM. The ID you see here is Flatcar.
	definition := vm.NewAPI(c).NewDefinition(location, "template", "44b38284-6adb-430e-b4a4-1553e29f352f", "developersfirstvm", 2, 2048, 10, networkInterfaces)
	definition.SSH = "<your SSH pub key>"

	// Provision the VM.
	provisionResponse, err := provisioning.VM().Provision(ctx, definition)
	if err != nil {
		panic(fmt.Sprintf("provisioning vm failed: %v", err))
	}

	// Wait for the VM to be ready.
	_, err = provisioning.Progress().AwaitCompletion(ctx, provisionResponse.Identifier)
	if err != nil {
		panic(fmt.Sprintf("waiting for VM provisioning failed: %v", err))
	}
}
```
