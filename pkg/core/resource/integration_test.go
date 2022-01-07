//go:build integration
// +build integration

package resource

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("core/resource API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should list resources", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		_, err := NewAPI(cli).List(ctx, 1, 1000)
		Expect(err).NotTo(HaveOccurred())
	})
})
