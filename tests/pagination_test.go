package tests_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/bind"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/frontend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/server"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LBaaS Service Tests", func() {
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

	Context("pagination Core Functionality", func() {
		const numberOfTests = 5
		It("testing pagination for frontend", func() {
			ctx := context.Background()

			for i := 0; i < numberOfTests; i++ {
				createFrontend(ctx, cli, nil)
			}

			api := frontend.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTests))

			// we already had the first page
			for i := 2; i < numberOfTests+1; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})

		It("testing pagination for backends", func() {
			ctx := context.Background()

			for i := 0; i < numberOfTests; i++ {
				createBackend(ctx, cli, nil)
			}

			api := backend.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTests))

			// we already had the first page
			for i := 2; i < numberOfTests+1; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})

		It("testing pagination for server", func() {
			ctx := context.Background()

			for i := 0; i < numberOfTests; i++ {
				createServer(ctx, cli, nil)
			}

			api := server.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTests))

			// we already had the first page
			for i := 2; i < numberOfTests+1; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})

		It("testing pagination for binds", func() {
			ctx := context.Background()

			for i := 0; i < numberOfTests; i++ {
				createBind(ctx, cli, nil)
			}

			api := bind.NewAPI(cli)
			page, err := api.GetPage(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(page.Size()).To(BeEquivalentTo(1))
			Expect(page.Total()).To(BeNumerically(">=", numberOfTests))

			// we already had the first page
			for i := 2; i < numberOfTests+1; i++ {
				page, err = api.NextPage(ctx, page)
				Expect(err).NotTo(HaveOccurred())
				Expect(page.Num()).To(BeEquivalentTo(i))
			}
		})
	})

	Context("pagination Utility Function", func() {
		It("pagination As Go Channels", func() {
			ctx := context.Background()

			const numberOfBackends = 5
			for i := 0; i < numberOfBackends; i++ {
				createBackend(ctx, cli, nil)
			}

			asChan, cancelFunc := pagination.AsChan(ctx, backend.NewAPI(cli))
			defer cancelFunc()

			counter := 0
			for elem := range asChan {
				Expect(elem).To(BeAssignableToTypeOf(backend.BackendInfo{}))
				counter++
			}
			Expect(counter).To(BeNumerically(">=", numberOfBackends))
		})

		It("pagination Loop until (all)", func() {
			ctx := context.Background()
			const numberOfTests = 5
			for i := 0; i < numberOfTests; i++ {
				createBackend(ctx, cli, nil)
			}

			counter := 0
			err := pagination.LoopUntil(ctx, backend.NewAPI(cli), func(i interface{}) (bool, error) {
				Expect(i).To(BeAssignableToTypeOf(backend.BackendInfo{}))
				counter++
				return false, nil
			})
			Expect(counter).To(BeNumerically(">=", numberOfTests))
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pagination.ErrConditionNeverMet))
		})

		It("pagination Loop until (subset)", func() {
			ctx := context.Background()
			const numberOfTests = 5
			for i := 0; i < numberOfTests; i++ {
				createBackend(ctx, cli, nil)
			}

			counter := 0
			err := pagination.LoopUntil(ctx, backend.NewAPI(cli), func(i interface{}) (bool, error) {
				Expect(i).To(BeAssignableToTypeOf(backend.BackendInfo{}))
				if counter == numberOfTests-1 {
					return true, nil
				}
				counter++
				return false, nil
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(counter).To(BeNumerically(">=", numberOfTests-1))
		})

		It("pagination Loop until (subset)", func() {
			ctx := context.Background()
			const numberOfTests = 5
			for i := 0; i < numberOfTests; i++ {
				createBackend(ctx, cli, nil)
			}

			counter := 0
			expectedErr := errors.New("test error")
			err := pagination.LoopUntil(ctx, backend.NewAPI(cli), func(i interface{}) (bool, error) {
				Expect(i).To(BeAssignableToTypeOf(backend.BackendInfo{}))
				if counter == numberOfTests-1 {
					return false, expectedErr
				}
				counter++
				return false, nil
			})

			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(expectedErr))
			Expect(counter).To(BeNumerically(">=", numberOfTests-1))
		})
	})
})
