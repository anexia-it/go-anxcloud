package api

import (
	"context"
	"fmt"
	"log"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	lbaasCommon "github.com/anexia-it/go-anxcloud/pkg/lbaas/common"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/loadbalancer"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

func ExampleNewAPI() {
	api, err := NewAPI(
		// you might find client.TokenFromEnv(false) useful
		WithClientOptions(client.TokenFromString("bogus auth token")),
	)

	if err != nil {
		log.Fatalf("Error creating api instance: %v\n", err)
	} else {
		// do something with api
		lb := loadbalancer.Loadbalancer{Identifier: "bogus identifier"}
		if err := api.Get(context.TODO(), &lb); IgnoreNotFound(err) != nil {
			fmt.Printf("Error retrieving loadbalancer with identifier '%v'\n", lb.Identifier)
		}
	}

	// fails because we didn't pass a valid auth token nor a valid identifier

	// Output: Error retrieving loadbalancer with identifier 'bogus identifier'
}

func Example_usage() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	// retrieve and create backend, handling errors along the way.
	backend := backend.Backend{Identifier: "bogus identifier 1"}
	if err := apiClient.Get(context.TODO(), &backend); IgnoreNotFound(err) != nil {
		fmt.Printf("Fatal error while retrieving backend: %v\n", err)
	} else if err != nil {
		fmt.Printf("Backend not yet existing, creating ...\n")

		backend.Name = "backend-01"
		backend.Mode = lbaasCommon.HTTP
		// [...]

		if err := apiClient.Create(context.TODO(), &backend); err != nil {
			fmt.Printf("Fatal error while creating backend: %v\n", err)
		}
	} else {
		fmt.Printf("Found backend with name %v and mode %v\n", backend.Name, backend.Mode)
		fmt.Printf("Deleting it for fun and profit :)\n")

		if err := apiClient.Destroy(context.TODO(), &backend); err != nil {
			fmt.Printf("Error destroying the backend: %v", err)
		}
	}

	// Output:
	// Found backend with name Example-Backend and mode tcp
	// Deleting it for fun and profit :)
}

func ExampleAPI_create() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	backend := backend.Backend{
		Name: "backend-01",
		Mode: lbaasCommon.HTTP,
		// [...]
	}

	if err := apiClient.Create(context.TODO(), &backend); err != nil {
		fmt.Printf("Error creating backend: %v\n", err)
	} else {
		fmt.Printf("Created backend '%v', engine assigned identifier '%v'\n", backend.Name, backend.Identifier)
	}

	// Output: Created backend 'backend-01', engine assigned identifier 'generated identifier 1'
}

func ExampleAPI_destroy() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	backend := backend.Backend{Identifier: "bogus identifier 1"}
	if err := apiClient.Destroy(context.TODO(), &backend); err != nil {
		fmt.Printf("Error destroying backend: %v\n", err)
	} else {
		fmt.Printf("Successfully destroyed backend\n")
	}

	// Output: Successfully destroyed backend
}

func ExampleAPI_get() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	backend := backend.Backend{Identifier: "bogus identifier 1"}
	if err := apiClient.Get(context.TODO(), &backend); err != nil {
		fmt.Printf("Error retrieving backend: %v\n", err)
	} else {
		fmt.Printf("Got backend named \"%v\"\n", backend.Name)
	}

	// Output: Got backend named "Example-Backend"
}

func ExampleAPI_listPaged() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	// List all backends, with 10 entries per page and starting on first page.

	// Beware: listing endpoints usually do not return all data for an object, sometimes
	// only the identifier is filled. This varies by specific API.
	b := backend.Backend{}
	var pageIter types.PageInfo
	if err := apiClient.List(context.TODO(), &b, Paged(1, 2, &pageIter)); err != nil {
		fmt.Printf("Error listing backends: %v\n", err)
	} else {
		var backends []backend.Backend
		for pageIter.Next(&backends) {
			fmt.Printf("Listing entries on page %v\n", pageIter.CurrentPage())

			for _, backend := range backends {
				fmt.Printf("  Got backend named \"%v\"\n", backend.Name)
			}
		}

		if err := pageIter.Error(); err != nil {
			// Handle error catched while iterating pages.
			// Errors will prevent pageIter.Next() to continue, you can call pageIter.ResetError() to resume.
			fmt.Printf("Error while iterating pages of backends: %v\n", err)
		}

		fmt.Printf("Last page listed was page %v, which returned %v entries\n", pageIter.CurrentPage(), len(backends))
	}

	// Output:
	// Listing entries on page 1
	//   Got backend named "Example-Backend"
	//   Got backend named "backend-01"
	// Listing entries on page 2
	//   Got backend named "test-backend-01"
	//   Got backend named "test-backend-02"
	// Listing entries on page 3
	//   Got backend named "test-backend-03"
	//   Got backend named "test-backend-04"
	// Last page listed was page 4, which returned 0 entries
}

func ExampleAPI_listChannel() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	channel := make(types.ObjectChannel)

	// list all backends using a channel and have the library handle the paging.
	// Oh and we filter by LoadBalancer, because we can and the example has to be somewhere.

	// Beware: listing endpoints usually do not return all data for an object, sometimes
	// only the identifier is filled. This varies by specific API.
	b := backend.Backend{LoadBalancer: loadbalancer.LoadBalancerInfo{Identifier: "bogus identifier 2"}}
	if err := apiClient.List(context.TODO(), &b, ObjectChannel(&channel)); err != nil {
		fmt.Printf("Error listing backends: %v\n", err)
	} else {
		for res := range channel {
			if err = res(&b); err != nil {
				fmt.Printf("Error retrieving backend from channel: %v\n", err)
				break
			}

			fmt.Printf("Got backend named \"%v\"\n", b.Name)
		}
	}

	// Output:
	// Got backend named "Example-Backend"
	// Got backend named "test-backend-02"
	// Got backend named "test-backend-04"
}

func ExampleAPI_update() {
	// see example on NewAPI how to implement this function
	apiClient := newExampleAPI()

	b := backend.Backend{
		Identifier: "bogus identifier 1",
		Name:       "Updated backend",
		Mode:       lbaasCommon.HTTP,
		// [...]
	}

	if err := apiClient.Update(context.TODO(), &b); err != nil {
		fmt.Printf("Error updating backend: %v\n", err)
	} else {
		fmt.Printf("Successfully updated backend\n")

		retrieved := backend.Backend{Identifier: "bogus identifier 1"}
		if err := apiClient.Get(context.TODO(), &retrieved); err != nil {
			fmt.Printf("Error verifying updated backend: %v\n", err)
		} else {
			fmt.Printf("Backend is now renamed to '%v' and has mode %v\n", retrieved.Name, retrieved.Mode)
		}
	}

	// Output:
	// Successfully updated backend
	// Backend is now renamed to 'Updated backend' and has mode http
}

// creates a new API instance for using the examples as tests. Includes a mock server.
func newExampleAPI() API {
	server := newMockServer()

	apiClient, err := NewAPI(
		WithClientOptions(
			client.BaseURL(server.URL()),
			client.TokenFromString("bogus testing token"),
		),
	)
	if err != nil {
		log.Fatalf("Error creating API client: %v\n", err)
		return nil
	}

	return apiClient
}
