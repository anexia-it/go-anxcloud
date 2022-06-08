package v1

import (
	"context"
	"log"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

func ExampleACL() {
	a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
	if err != nil {
		log.Fatalf("failed to initialize api client: %s", err)
	}

	acl := &ACL{
		Name:       "destination port 8080",
		ParentType: "frontend",
		Index:      pointer.Int(5),
		Criterion:  "dst_port",
		Value:      "8080",
		Frontend: Frontend{
			Identifier: "<frontend-identifier>",
		},
	}

	if err := a.Create(context.TODO(), acl); err != nil {
		log.Fatalf("failed to create ACL: %s", err)
	}
}
