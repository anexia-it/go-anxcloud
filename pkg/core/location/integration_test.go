// +build integration
// go:build integration

package location

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("core/location API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should list available locations", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		_, err := NewAPI(cli).List(ctx, 1, 1000, "")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should get the first location entry by ID", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		list, err := NewAPI(cli).List(ctx, 1, 1000, "")
		Expect(err).NotTo(HaveOccurred())

		l, err := NewAPI(cli).Get(ctx, list[0].ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(l).To(Equal(list[0]))
	})

	It("should get the first location entry by code", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		list, err := NewAPI(cli).List(ctx, 1, 1000, "")
		Expect(err).NotTo(HaveOccurred())

		l, err := NewAPI(cli).GetByCode(ctx, list[0].Code)
		Expect(err).NotTo(HaveOccurred())
		Expect(l).To(Equal(list[0]))
	})
})
