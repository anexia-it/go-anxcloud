package tests_test

import (
	"context"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/core/location"
	"github.com/anexia-it/go-anxcloud/pkg/core/resource"
	"github.com/anexia-it/go-anxcloud/pkg/core/service"
	"github.com/anexia-it/go-anxcloud/pkg/core/tags"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("core API endpoint tests", func() {

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

	Describe("location endpoint", func() {

		It("should list all available locations, and get the first entry by ID and code", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get the first location entry by ID", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			list, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
			l, err := location.NewAPI(cli).Get(ctx, list[0].ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(l).To(Equal(list[0]))
		})

		It("should get the first location entry by code", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			list, err := location.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())

			l, err := location.NewAPI(cli).GetByCode(ctx, list[0].Code)
			Expect(err).NotTo(HaveOccurred())
			Expect(l).To(Equal(list[0]))
		})
	})

	Describe("resource endpoint", func() {
		It("should list all created resources", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			_, err := resource.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("with at least one resource existing", func() {
			var ctx context.Context
			JustBeforeEach(func() {
				ctx = context.Background()
				createBackend(ctx, cli, nil)
			})

			It("should list resource using generic API client", func() {
				genericAPI, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
				Expect(err).ToNot(HaveOccurred())

				var pageIter types.PageInfo
				err = genericAPI.List(ctx, &resource.Info{}, api.Paged(1, 100, &pageIter))
				Expect(err).ToNot(HaveOccurred())

				var resInfo []resource.Info
				Expect(pageIter.Next(&resInfo)).To(BeTrue())
				Expect(resInfo).ToNot(BeEmpty())
				Expect(resInfo[0].Identifier).ToNot(BeEmpty())
			})

			It("should throw an error for unsupported operations for the genric API client", func() {
				genericAPI, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
				Expect(err).ToNot(HaveOccurred())
				err = genericAPI.Create(ctx, &resource.Info{})
				Expect(err).To(BeEquivalentTo(api.ErrOperationNotSupported))
			})
		})
	})

	Describe("service endpoint", func() {
		It("should list all created services", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := service.NewAPI(cli).List(ctx, 1, 1000)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("tags endpoint", func() {
		It("should list all created tags", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := tags.NewAPI(cli).List(ctx, 1, 1000, "", "", "", "", true)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
