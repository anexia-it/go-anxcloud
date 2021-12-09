//go:build integration
// +build integration

package nictype

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NIC type API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("lists available NIC types", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		nicTypes, err := NewAPI(cli).List(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(nicTypes)).To(BeNumerically(">", 1))
	})
})
