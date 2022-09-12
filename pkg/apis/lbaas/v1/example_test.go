package v1_test

import (
	"context"
	"log"

	"go.anx.io/go-anxcloud/pkg/api"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

func ExampleACL() {
	a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
	if err != nil {
		log.Fatalf("failed to initialize api client: %s", err)
	}

	acl := &lbaasv1.ACL{
		Name:       "destination port 8080",
		ParentType: "frontend",
		Index:      pointer.Int(5),
		Criterion:  "dst_port",
		Value:      "8080",
		Frontend: lbaasv1.Frontend{
			Identifier: "<frontend-identifier>",
		},
	}

	if err := a.Create(context.TODO(), acl); err != nil {
		log.Fatalf("failed to create ACL: %s", err)
	}
}
