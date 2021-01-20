package tests_test

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/core/location"
	"github.com/anexia-it/go-anxcloud/pkg/core/resource"
	"github.com/anexia-it/go-anxcloud/pkg/core/service"
	"github.com/anexia-it/go-anxcloud/pkg/core/tags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Location endpoint", func() {

		It("Should list all available locations", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := location.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Resource endpoint", func() {

		It("Should list all created resources", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := resource.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Service endpoint", func() {

		It("Should list all created services", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := service.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Tags endpoint", func() {

		It("Should list all created tags", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := tags.NewAPI(cli).List(ctx, 1, 1000, "", "", "", "", true)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
