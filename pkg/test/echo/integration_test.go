//go:build integration
// +build integration

package echo

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("echo api client", func() {
	It("should be able to communicate with Anexia echo endpoint", func() {
		c, err := client.New(client.TokenFromEnv(false))
		Expect(err).NotTo(HaveOccurred())

		ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
		defer cancel()

		err = NewAPI(c).Echo(ctx)
		Expect(err).NotTo(HaveOccurred())
	})
})
