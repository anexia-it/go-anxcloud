//go:build integration
// +build integration

package server

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/utils/test"

	lbaasBackend "github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/server client", func() {
	var cli client.Client
	var api API
	var backend lbaasBackend.Backend

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)

		backendAPI := lbaasBackend.NewAPI(cli)

		b, err := backendAPI.Create(context.TODO(), lbaasBackend.Definition{
			Name:         test.TestResourceName(),
			State:        common.NewlyCreated,
			LoadBalancer: loadbalancerIdentifier,
			Mode:         common.TCP,
		})
		Expect(err).To(BeNil())

		DeferCleanup(func() {
			err := backendAPI.DeleteByID(context.TODO(), b.Identifier)
			Expect(err).To(BeNil())
		})

		backend = b
	})

	createServer := func(definition Definition) Server {
		s, err := api.Create(context.TODO(), definition)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			err := api.DeleteByID(context.TODO(), s.Identifier)
			Expect(err).NotTo(HaveOccurred())
		})

		Expect(s.Name).To(Equal(definition.Name))
		Expect(s.Port).To(Equal(definition.Port))
		Expect(s.IP).To(Equal(definition.IP))
		Expect(s.Backend).NotTo(BeNil())
		Expect(s.Backend.Identifier).To(Equal(backend.Identifier))

		return s
	}

	Context("with a server created for testing", func() {
		var definition Definition
		var server Server

		BeforeEach(func() {
			definition = Definition{
				Name:    test.TestResourceName(),
				State:   common.NewlyCreated,
				IP:      "8.8.8.8",
				Port:    8080,
				Backend: backend.Identifier,
			}

			server = createServer(definition)
		})

		It("lists servers including our test server", func() {
			found := false
			page := 1

			for !found {
				ss, err := api.Get(context.TODO(), page, 20)
				Expect(err).To(BeNil())
				Expect(ss).NotTo(BeEmpty())

				for _, s := range ss {
					if s.Identifier == server.Identifier {
						found = true
						break
					}
				}

				page++
			}
		})

		It("retrieves test server with expected data", func() {
			s, err := api.GetByID(context.TODO(), server.Identifier)
			Expect(err).To(BeNil())

			Expect(s).To(Equal(server))
		})

		It("updates test server with changed values", func() {
			definition := Definition{
				Name:    test.TestResourceName(),
				State:   common.Updating,
				IP:      "8.8.4.4",
				Port:    5353,
				Backend: backend.Identifier,
			}

			s, err := api.Update(context.TODO(), server.Identifier, definition)
			Expect(err).To(BeNil())

			Expect(s.Identifier).To(Equal(server.Identifier))
			Expect(s.Name).To(Equal(definition.Name))
			Expect(s.IP).To(Equal(definition.IP))
			Expect(s.Port).To(Equal(definition.Port))
		})
	})

	Context("with some servers created for testing", func() {
		const numberOfTestServers = 5

		BeforeEach(func() {
			for i := 0; i < numberOfTestServers; i++ {
				createServer(Definition{
					Name:    test.TestResourceName(),
					State:   common.NewlyCreated,
					IP:      "8.8.8.8",
					Port:    8080,
					Backend: backend.Identifier,
				})
			}
		})

		It("iterates through pages as expected", func() {
			page, err := api.GetPage(context.TODO(), 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTestServers))

			// we already had the first page
			for i := 2; i < numberOfTestServers+1; i++ {
				page, err = api.NextPage(context.TODO(), page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})
	})
})
