package tests_test

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/test/echo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client tests", func() {

	Context("Echo endpoint", func() {

		It("Should be able to communicate with Anexia echo endpoint", func() {
			c, err := client.New(client.TokenFromEnv(false))
			Expect(err).NotTo(HaveOccurred())

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()

			err = echo.NewAPI(c).Echo(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
