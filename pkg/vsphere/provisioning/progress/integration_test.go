// +build integration
// go:build integration

package progress

import (
	"context"
	"errors"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("vsphere/provisioning/disktype API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when using an identifier which does not exist", func() {
		It("receives a 404 response and handles it", func() {
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()
			progress, err := NewAPI(cli).AwaitCompletion(ctx, "this-id-does-not-exist")
			Expect(progress).To(BeEmpty())

			re := &client.ResponseError{}
			Expect(errors.As(err, &re)).To(BeTrue())
			Expect(re.Response.StatusCode).To(Equal(404))
		})
	})
})
