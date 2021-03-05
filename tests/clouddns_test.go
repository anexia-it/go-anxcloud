package tests_test

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
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
		var createTestZoneName = "sdk-create-test.xocp.de"
		AfterEach(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			var _ = zone.NewAPI(cli).Delete(ctx, createTestZoneName)
		})

		It("Should create the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			createDefinition := zone.Definition{
				ZoneName:   createTestZoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "admin@xocp.de",
				Refresh:    300,
				Retry:      300,
				Expire:     3600,
				TTL:        300,
			}
			response, err := zone.NewAPI(cli).Create(ctx, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.Name).To(Equal(createTestZoneName))
			Expect(response.AdminEmail).To(Equal("admin@xocp.de"))
		})
	})

	Context("Definition Update Endpoint", func() {
		updateTestZoneName := "sdk-update-test.xocp.de"

		It("Should update the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			randRefresh := rand.Intn(10) * 100
			randRetry := rand.Intn(10) * 100
			randExpire := rand.Intn(10) * 1000
			randTTL := rand.Intn(10) * 100
			createDefinition := zone.Definition{
				ZoneName:   updateTestZoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "test@xocp.de",
				Refresh:    randRefresh,
				Retry:      randRetry,
				Expire:     randExpire,
				TTL:        randTTL,
			}
			response, err := zone.NewAPI(cli).Update(ctx, updateTestZoneName, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.AdminEmail).To(Equal("test@xocp.de"))
			Expect(response.Refresh).To(Equal(randRefresh))
			Expect(response.Retry).To(Equal(randRetry))
			Expect(response.Expire).To(Equal(randExpire))
			Expect(response.TTL).To(Equal(randTTL))
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

	Context("Definition List Recoreds Endpoint", func() {
		It("Should list all available records for the test zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).ListRecords(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
