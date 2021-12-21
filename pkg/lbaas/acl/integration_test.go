//go:build integration
// +build integration

package acl

import (
	"context"
	"math/rand"

	lbaasBackend "go.anx.io/go-anxcloud/pkg/lbaas/backend"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const loadbalancerIdentifier = "fc5d7390e9e4400a9efc73b4d8e0613a"

var _ = Describe("lbaas/acl client", func() {
	var cli client.Client
	var api API
	var backend lbaasBackend.Backend

	BeforeEach(OncePerOrdered, func() {
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

	createACL := func() string {
		definition := &Definition{
			Name:       test.TestResourceName(),
			State:      common.NewlyCreated,
			ParentType: "backend",
			Criterion:  "src",
			Index:      rand.Intn(100),
			Value:      "10.0.0.0/8",
			Backend:    &backend.Identifier,
		}

		acl, err := api.Create(context.TODO(), *definition)
		Expect(err).To(BeNil())

		Expect(acl.Name).To(BeEquivalentTo(definition.Name))
		Expect(acl.ParentType).To(BeEquivalentTo(definition.ParentType))
		Expect(acl.Index).To(BeEquivalentTo(definition.Index))
		Expect(acl.Value).To(BeEquivalentTo(definition.Value))
		Expect(acl.Backend.Identifier).To(BeEquivalentTo(*definition.Backend))
		Expect(acl.Criterion).To(BeEquivalentTo(definition.Criterion))

		return acl.Identifier
	}

	deleteACL := func(identifier string) {
		err := api.DeleteByID(context.TODO(), identifier)
		Expect(err).To(BeNil())
	}

	Context("working on a fresh ACL", Ordered, func() {
		var aclIdentifier string

		It("creates an ACL", func() {
			aclIdentifier = createACL()
		})

		It("deletes the created ACL", func() {
			deleteACL(aclIdentifier)
		})
	})

	Context("with an ACL created for testing", func() {
		var aclIdentifier string

		BeforeEach(func() {
			aclIdentifier = createACL()

			DeferCleanup(func() {
				deleteACL(aclIdentifier)
			})
		})

		It("lists ACLs including our test ACL", func() {
			found := false
			page := 1

			for !found {
				acls, err := api.Get(context.TODO(), page, 20)
				Expect(err).To(BeNil())
				Expect(acls).NotTo(BeEmpty())

				for _, acl := range acls {
					if acl.Identifier == aclIdentifier {
						found = true
						break
					}
				}

				page++
			}
		})

		It("retrieves test ACL with expected data", func() {
			acl, err := api.GetByID(context.TODO(), aclIdentifier)
			Expect(err).To(BeNil())

			Expect(acl.ParentType).To(Equal("backend"))
			Expect(acl.Criterion).To(Equal("src"))
			Expect(acl.Value).To(Equal("10.0.0.0/8"))
			Expect(acl.Backend).NotTo(BeNil())
			Expect(acl.Backend.Identifier).To(Equal(backend.Identifier))
		})

		It("updates test ACL with changed values", func() {
			acl, err := api.Update(context.TODO(), aclIdentifier, Definition{
				Name:       test.TestResourceName(),
				State:      common.Updating,
				ParentType: "backend",
				Criterion:  "src",
				Index:      rand.Intn(100),
				Value:      "172.20.0.0/12",
				Frontend:   nil,
				Backend:    &backend.Identifier,
			})
			Expect(err).To(BeNil())

			Expect(acl.Identifier).To(BeEquivalentTo(aclIdentifier))

			Expect(acl.ParentType).To(Equal("backend"))
			Expect(acl.Criterion).To(Equal("src"))
			Expect(acl.Value).To(Equal("172.20.0.0/12"))
			Expect(acl.Backend).NotTo(BeNil())
			Expect(acl.Backend.Identifier).To(Equal(backend.Identifier))
		})
	})
})
