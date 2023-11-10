package v1

import (
	"context"
	"log"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

func Example() {
	a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
	if err != nil {
		log.Fatalf("failed to initialize API client: %s", err)
	}

	cluster := Cluster{Name: "example", NeedsServiceVMs: pointer.Bool(true)}

	if err := a.Create(context.TODO(), &cluster); err != nil {
		log.Fatalf("failed to create cluster: %s", err)
	}

	if err := gs.AwaitCompletion(context.TODO(), a, &cluster); err != nil {
		log.Fatalf("failed to await cluster creation: %s", err)
	}

	// define node pool with a single replica, 2 GiB of system memory and 20 GiB of disk space
	nodePool := NodePool{Name: "example-np-00", Cluster: cluster, Replicas: pointer.Int(1), Memory: 2 * 1073741824, DiskSize: 20 * 1073741824}

	if err := a.Create(context.TODO(), &nodePool); err != nil {
		log.Fatalf("failed to create nodepool: %s", err)
	}

	if err := gs.AwaitCompletion(context.TODO(), a, &nodePool); err != nil {
		log.Fatalf("failed to await nodepool creation: %s", err)
	}
}
