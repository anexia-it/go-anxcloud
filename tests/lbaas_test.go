package tests

import (
	"context"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/acl"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/bind"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/frontend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/loadbalancer"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
)

type CleanUpHandler = func() error

var cleanupHandlers []CleanUpHandler

var _ = FDescribe("LBaaS Service Tests", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		for _, handler := range cleanupHandlers {
			err := handler()
			if err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "error when cleaning up tests: %s", err.Error())
			}
		}

		cleanupHandlers = []CleanUpHandler{}
	})

	Context("LBaaS - Loadbalancers", func() {
		It("Get load balancers", func() {
			ctx := context.Background()

			loadBalancers, err := lbaas.NewAPI(cli).LoadBalancer().Get(ctx, 1, 50)

			Expect(err).To(BeNil())
			Expect(loadBalancers).NotTo(BeNil())
			Expect(len(loadBalancers)).Should(BeNumerically(">=", 2))
		})

		It("Get a specific load balancer", func() {
			ctx := context.Background()
			api := lbaas.NewAPI(cli).LoadBalancer()

			loadBalancers, err := api.Get(ctx, 1, 50)

			Expect(err).To(BeNil())
			Expect(loadBalancers).NotTo(BeNil())
			Expect(len(loadBalancers)).Should(BeNumerically(">=", 2))

			loadBalancer, err := api.GetByID(ctx, loadBalancers[0].Identifier)

			Expect(err).To(BeNil())
			Expect(loadBalancer.Identifier).To(BeEquivalentTo(loadBalancers[0].Identifier))
			Expect(loadBalancer.Name).To(BeEquivalentTo(loadBalancers[0].Name))
		})
	})

	Context("LBaaS - Backend", func() {
		It("Create a Backend", func() {
			ctx := context.Background()
			definition := &backend.Definition{
				Name:         randomName(),
				State:        common.NewlyCreated,
				LoadBalancer: getFirstLB(ctx, cli).Identifier,
				Mode:         common.TCP,
			}

			backend := createBackend(ctx, cli, definition)

			Expect(backend.Name).To(BeEquivalentTo(definition.Name))
			Expect(backend.Mode).To(BeEquivalentTo(definition.Mode))
			Expect(backend.LoadBalancer.Identifier).To(BeEquivalentTo(definition.LoadBalancer))
			Expect(backend.Identifier).ToNot(BeEmpty())
		})

		It("Get Backends", func() {
			ctx := context.Background()
			createBackend(ctx, cli, nil)

			backends, err := backend.NewAPI(cli).Get(ctx, 1, 5)

			Expect(err).To(BeNil())
			Expect(backends).ToNot(BeEmpty())
		})

		It("Get a specific backend", func() {
			ctx := context.Background()
			testBackend := createBackend(ctx, cli, nil)

			fetchedBackend, err := backend.NewAPI(cli).GetByID(ctx, testBackend.Identifier)

			Expect(err).To(BeNil())
			Expect(fetchedBackend).To(BeEquivalentTo(testBackend))
		})
	})

	Context("LBaaS - Servers", func() {
		It("Create server", func() {
			ctx := context.Background()

			definition := &server.Definition{
				Name:    randomName(),
				State:   common.NewlyCreated,
				IP:      "8.8.8.8",
				Port:    8080,
				Backend: createBackend(ctx, cli, nil).Identifier,
			}

			createdServer := createServer(ctx, cli, definition)

			Expect(createdServer.Name).To(BeEquivalentTo(definition.Name))
			Expect(createdServer.Port).To(BeEquivalentTo(definition.Port))
			Expect(createdServer.IP).To(BeEquivalentTo(definition.IP))
			Expect(createdServer.Backend.Identifier).To(BeEquivalentTo(definition.Backend))
		})

		It("Get servers", func() {
			ctx := context.Background()
			createServer(ctx, cli, nil)

			servers, err := server.NewAPI(cli).Get(ctx, 1, 5)
			Expect(err).To(BeNil())
			Expect(servers).ToNot(BeEmpty())
		})

		It("Get a specific server", func() {
			ctx := context.Background()
			createdServer := createServer(ctx, cli, nil)

			fetchedServer, err := server.NewAPI(cli).GetByID(ctx, createdServer.Identifier)
			Expect(err).To(BeNil())
			Expect(fetchedServer).To(BeEquivalentTo(createdServer))
		})
	})

	Context("LBaaS - Binds", func() {
		It("Create Bind", func() {
			ctx := context.Background()
			definition := &bind.Definition{
				Name:     randomName(),
				Frontend: createFrontend(ctx, cli, nil).Identifier,
				State:    common.NewlyCreated,
			}
			createdBind := createBind(ctx, cli, definition)
			Expect(createdBind.Name).To(BeEquivalentTo(definition.Name))
			Expect(createdBind.Frontend.Identifier).To(BeEquivalentTo(definition.Frontend))
		})
		It("Get Binds", func() {
			ctx := context.Background()
			createBind(ctx, cli, nil)

			binds, err := bind.NewAPI(cli).Get(ctx, 1, 5)
			Expect(err).To(BeNil())
			Expect(binds).ToNot(HaveLen(0))
		})
		It("Get a specific Bind", func() {
			ctx := context.Background()
			createdBind := createBind(ctx, cli, nil)

			fetchedBind, err := bind.NewAPI(cli).GetByID(ctx, createdBind.Identifier)
			Expect(err).To(BeNil())
			Expect(fetchedBind).To(BeEquivalentTo(createdBind))
		})
	})

	Context("LBaaS - Frontends", func() {
		It("Create frontend", func() {
			ctx := context.Background()
			backend := createBackend(ctx, cli, nil)
			definition := frontend.Definition{
				Name:           randomName(),
				LoadBalancer:   getFirstLB(ctx, cli).Identifier,
				DefaultBackend: backend.Identifier,
				Mode:           common.TCP,
				State:          common.NewlyCreated,
			}

			frontend := createFrontend(ctx, cli, &definition)

			Expect(frontend.Name).To(BeEquivalentTo(definition.Name))
			Expect(frontend.Mode).To(BeEquivalentTo(definition.Mode))
			Expect(frontend.LoadBalancer.Identifier).To(BeEquivalentTo(definition.LoadBalancer))
		})

		It("Get load balancer frontends", func() {
			ctx := context.Background()
			createFrontend(ctx, cli, nil)

			frontends, err := lbaas.NewAPI(cli).Frontend().Get(ctx, 1, 50)

			Expect(err).To(BeNil())
			Expect(frontends).ToNot(BeEmpty())
		})

		It("Get a specific frontend", func() {
			ctx := context.Background()
			api := lbaas.NewAPI(cli).Frontend()
			createdFrontend := createFrontend(ctx, cli, nil)

			fetchedFrontend, err := api.GetByID(ctx, createdFrontend.Identifier)

			Expect(err).To(BeNil())
			Expect(fetchedFrontend).To(BeEquivalentTo(createdFrontend))
		})
	})

	Context("LBaaS - ACLs", func() {
		It("Create ACL", func() {
			ctx := context.Background()
			backend := createBackend(ctx, cli, nil)
			definition := &acl.Definition{
				Name:       randomName(),
				State:      common.NewlyCreated,
				ParentType: "backend",
				Criterion:  "src",
				Index:      rand.Intn(100),
				Value:      "10.0.0.0/4",
				Backend:    &backend.Identifier,
			}

			createdACL := createACL(ctx, cli, definition)

			Expect(createdACL.Name).To(BeEquivalentTo(definition.Name))
			Expect(createdACL.ParentType).To(BeEquivalentTo(definition.ParentType))
			Expect(createdACL.Index).To(BeEquivalentTo(definition.Index))
			Expect(createdACL.Value).To(BeEquivalentTo(definition.Value))
			Expect(createdACL.Backend.Identifier).To(BeEquivalentTo(*definition.Backend))
			Expect(createdACL.Criterion).To(BeEquivalentTo(definition.Criterion))
		})

		It("Get ACLs", func() {
			ctx := context.Background()
			createACL(ctx, cli, nil)
			api := acl.NewAPI(cli)

			acls, err := api.Get(ctx, 1, 5)
			Expect(err).To(BeNil())
			Expect(acls).NotTo(BeEmpty())
		})

		It("Get specific ACL", func() {
			ctx := context.Background()
			createdACL := createACL(ctx, cli, nil)
			api := acl.NewAPI(cli)

			fetchedACL, err := api.GetByID(ctx, createdACL.Identifier)
			Expect(err).To(BeNil())
			Expect(fetchedACL).To(BeEquivalentTo(createdACL))
		})
	})
})

func createBind(ctx context.Context, cli client.Client, definition *bind.Definition) bind.Bind {
	api := bind.NewAPI(cli)
	if definition == nil {
		definition = &bind.Definition{
			Name:     randomName(),
			State:    common.NewlyCreated,
			Frontend: createFrontend(ctx, cli, nil).Identifier,
		}
	}
	createdBind, err := api.Create(ctx, *definition)
	Expect(err).To(BeNil())
	cleanUpAfterTest(bindWithID(createdBind.Identifier))
	return createdBind
}

func createBackend(ctx context.Context, cli client.Client, definition *backend.Definition) backend.Backend {
	api := backend.NewAPI(cli)
	if definition == nil {
		definition = &backend.Definition{
			Name:         randomName(),
			State:        common.NewlyCreated,
			LoadBalancer: getFirstLB(ctx, cli).Identifier,
			Mode:         common.TCP,
		}
	}

	backend, err := api.Create(ctx, *definition)
	Expect(err).To(BeNil())
	cleanUpAfterTest(backendWithID(backend.Identifier))
	return backend
}

func createServer(ctx context.Context, cli client.Client, definition *server.Definition) server.Server {
	api := server.NewAPI(cli)
	if definition == nil {
		definition = &server.Definition{
			Name:    randomName(),
			State:   common.NewlyCreated,
			IP:      "8.8.8.8",
			Port:    8080,
			Backend: createBackend(ctx, cli, nil).Identifier,
		}
	}
	createdServer, err := api.Create(ctx, *definition)
	Expect(err).To(BeNil())
	cleanUpAfterTest(serverWithID(createdServer.Identifier))
	return createdServer
}

func createACL(ctx context.Context, cli client.Client, definition *acl.Definition) acl.ACL {
	api := acl.NewAPI(cli)
	if definition == nil {
		backend := createBackend(ctx, cli, nil)
		definition = &acl.Definition{
			Name:       randomName(),
			State:      common.NewlyCreated,
			ParentType: "backend",
			Criterion:  "src",
			Index:      rand.Intn(100),
			Value:      "10.0.0.0/4",
			Backend:    &backend.Identifier,
		}
	}
	acl, err := api.Create(ctx, *definition)
	Expect(err).To(BeNil())
	cleanUpAfterTest(aclWithID(acl.Identifier))
	return acl
}

func createFrontend(ctx context.Context, cli client.Client, definition *frontend.Definition) frontend.Frontend {
	api := frontend.NewAPI(cli)
	if definition == nil {
		backend := createBackend(ctx, cli, nil)
		definition = &frontend.Definition{
			Name:           randomName(),
			State:          common.NewlyCreated,
			LoadBalancer:   getFirstLB(ctx, cli).Identifier,
			Mode:           common.TCP,
			DefaultBackend: backend.Identifier,
		}
	}

	frontend, err := api.Create(ctx, *definition)
	Expect(err).To(BeNil())
	cleanUpAfterTest(frontendWithID(frontend.Identifier))
	return frontend
}

func getFirstLB(ctx context.Context, cli client.Client) loadbalancer.Loadbalancer {
	api := lbaas.NewAPI(cli).LoadBalancer()
	loadBalancers, err := api.Get(ctx, 1, 50)
	Expect(err).To(BeNil())
	Expect(loadBalancers).ToNot(BeEmpty())

	loadbalancer, err := api.GetByID(ctx, loadBalancers[0].Identifier)
	Expect(err).To(BeNil())
	return loadbalancer
}

func randomName() string {
	return fmt.Sprintf("go-anxcloud-integration-test-random-resource-name-%d", rand.Intn(100000))
}

// cleanUpAfterTest registers cleanupHandlers which are executed after every test
func cleanUpAfterTest(handler ...CleanUpHandler) {
	// cleanup needs to happen in reverse order but keeping the order of the handlers that were passed in
	cleanupHandlers = append(handler, cleanupHandlers...)
}

func frontendWithID(identifier string) CleanUpHandler {
	return func() error {
		cli, err := client.New(client.AuthFromEnv(false))
		if err != nil {
			return err
		}
		return lbaas.NewAPI(cli).Frontend().DeleteByID(context.Background(), identifier)
	}
}

func backendWithID(identifier string) CleanUpHandler {
	return func() error {
		cli, err := client.New(client.AuthFromEnv(false))
		if err != nil {
			return err
		}
		return lbaas.NewAPI(cli).Backend().DeleteByID(context.Background(), identifier)
	}
}

func bindWithID(identifier string) CleanUpHandler {
	return func() error {
		cli, err := client.New(client.AuthFromEnv(false))
		if err != nil {
			return err
		}
		return lbaas.NewAPI(cli).Bind().DeleteByID(context.Background(), identifier)
	}
}

func aclWithID(identifier string) CleanUpHandler {
	return func() error {
		cli, err := client.New(client.AuthFromEnv(false))
		if err != nil {
			return err
		}
		return lbaas.NewAPI(cli).ACL().DeleteByID(context.Background(), identifier)
	}
}

func serverWithID(identifier string) CleanUpHandler {
	return func() error {
		cli, err := client.New(client.AuthFromEnv(false))
		if err != nil {
			return err
		}
		return lbaas.NewAPI(cli).Server().DeleteByID(context.Background(), identifier)
	}
}
