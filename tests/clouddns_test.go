package tests

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("CloudDNS API endpoint tests", func() {
	var cli client.Client

	const TestZone string = "xocp.de"

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Definition List Endpoint", func() {
		It("Should list all available zones", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).List(ctx)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Definition Get Endpoint", func() {
		It("Should return the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).Get(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("Definition Create Endpoint", func() {
		It("Should create the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			createDefinition := zone.Definition{
				ZoneName:   "sdk-test.xocp.de",
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "amdin@xocp.de",
				Refresh:    300,
				Retry:      300,
				Expire:     3600,
				TTL:        300,
			}
			response, err := zone.NewAPI(cli).Create(ctx, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.Name).To(Equal("sdk-test.xocp.de"))
		})
	})
	Context("Definition Update Endpoint", func() {
		It("Should update the zone", func() {
			// TODO
		})
	})

	Context("Definition Delete Endpoint", func() {
		It("Should delete the zone", func() {
			// TODO
		})
	})

	Context("Definition ChangeSet Endpoint", func() {
		It("Should apply the changeset", func() {
			// TODO
		})
	})

	Context("Definition Import Endpoint", func() {
		It("Should import the zone", func() {
			// TODO
		})
	})

})
