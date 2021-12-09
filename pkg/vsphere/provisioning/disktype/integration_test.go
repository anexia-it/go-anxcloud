//go:build integration
// +build integration

package disktype

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
)

var _ = Describe("vsphere/provisioning/disktype API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("lists available disk types", func() {
		ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
		defer cancel()

		_, err := NewAPI(cli).List(ctx, locationID, 1, 1000)
		Expect(err).NotTo(HaveOccurred())
	})
})
