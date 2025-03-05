[![Documentation](https://godoc.org/go.anx.io/go-anxcloud?status.svg)](http://godoc.org/go.anx.io/go-anxcloud)
[![codecov](https://codecov.io/gh/anexia-it/go-anxcloud/branch/main/graph/badge.svg?token=G4XZW5U5WT)](https://codecov.io/gh/anexia-it/go-anxcloud)

# Go Client for the Anexia API

Go SDK for interacting with the Anexia Engine [API](https://engine.anexia-it.com/).

## Installing

To use the SDK, just add `go.anx.io/go-anxcloud <version>` to your Go module.

## Getting started

Before using the SDK you should familiarize yourself with the [Anexia Engine API](https://engine.anexia-it.com/docs/).

The library is used in [our terraform provider](https://github.com/anexia-it/terraform-provider-anxcloud), check it out if you want some examples how to use it.

### Example

Below is a short example using the new generic client in this package. Not all APIs can already be used with it, but we are working on that.
Find [more examples in the docs](https://pkg.go.dev/go.anx.io/go-anxcloud@main/pkg/api#example-package-Usage) (linked to docs for
main branch, not the latest (or any) release).

```go
package main

import (
	"context"
	"log"

	"go.anx.io/go-anxcloud/pkg/api"
	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"

	// apis usable with the generic client have their own package in a location analog to this
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
)

func main() {
	apiClient, err := api.NewAPI(
		api.WithClientOptions(
			// Get auth token from ANEXIA_TOKEN environment variable.
			// The boolean parameter specifies if the environment variable should be unset.
			client.TokenFromEnv(false),
		),
	)
	if err != nil {
		log.Fatalf("Error creating ANX API client: %v", err)
	}

	// let's list LBaaS backends of a known LoadBalancer
	frontend := lbaasv1.Frontend{
		LoadBalancer: &lbaasv1.LoadBalancer{Identifier: "285b954fdf2a449c8fdae01cc6074025"},
	}

	var frontends apiTypes.ObjectChannel
	err = apiClient.List(context.TODO(), &frontend,
		// Listing can be done with either a page iterator or a channel, we use a channel here.
		api.ObjectChannel(&frontends),

		// Most APIs only give a very small subset when listing resources, add this flag to
		// get all attributes, at the cost of doing lots of API requests.
		api.FullObjects(true),
	)
	if err != nil {
		log.Fatalf("Error listing backends for LoadBalancer '%v': %v", frontend.LoadBalancer.Identifier, err)
	}

	for retriever := range frontends {
		// reinitialise frontend every loop to reset pointers and avoid potential overwriting of data in the next loop
		var frontend lbaasv1.Frontend
		if err := retriever(&frontend); err != nil {
			log.Fatalf("Error retrieving Frontend: %v", err)
		}

		log.Printf("Got Frontend named '%v' with mode '%v'", frontend.Name, frontend.Mode)
	}
}
```

This new generic client will one day be the only client in go-anxcloud. The legacy API-specific clients are deprecated and will be removed in the
go-anxcloud release following the one with all APIs go-anxcloud supports usable with the generic client (so if the generic client in 0.5.0 supports
at least everything there is another client for in go-anxcloud, 0.6.0 will drop the API-specific clients).

<details>
<summary>Example how to create a VM with the API-specific, deprecated client.</summary>

```go
package main

import (
	"context"
	"fmt"
	"time"

	anexia "go.anx.io/go-anxcloud/pkg"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/vm"
)

func main() {
	vlan := "<ID of the VLAN the VM should have access to>"
	location := "<ID of the location the VM should be in>"

	// Create client using the auth token in environment variable ANEXIA_TOKEN and do not unset the environment variable.
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
	networkInterfaces := []vm.Network{{NICType: "virtio", IPs: []string{ips[0].Identifier}, VLAN: vlan}}
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
</details>
