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

		It("Should list all available locations, and get the first entry by ID and code", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should paginate correctly", func() {
			ctx := context.Background()

			api := location.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">", 1))

			for i := 2; i < 5; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})

		It("Should get the first location entry by ID", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			list, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
			l, err := location.NewAPI(cli).Get(ctx, list[0].ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(l).To(Equal(list[0]))
		})

		It("Should get the first location entry by code", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			list, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())

			l, err := location.NewAPI(cli).GetByCode(ctx, list[0].Code)
			Expect(err).NotTo(HaveOccurred())
			Expect(l).To(Equal(list[0]))
		})
	})

	Context("Resource endpoint", func() {

		It("Should list all created resources", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := resource.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should paginate correctly", func() {
			ctx := context.Background()

			api := resource.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">", 1))

			for i := 2; i < 5; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})

	})

	Context("Service endpoint", func() {

		It("Should list all created services", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := service.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should paginate correctly", func() {
			ctx := context.Background()

			api := service.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">", 1))

			for i := 2; i < 5; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
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
