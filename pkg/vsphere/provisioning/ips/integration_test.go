// +build integration
// go:build integration

package ips

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "02f39d20ca0f4adfb5032f88dbc26c39"
)

var _ = Describe("vsphere/provisioning/ips API client", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	It("retrieves a free IP address", func() {
		ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
		defer cancel()

		_, err := NewAPI(cli).GetFree(ctx, locationID, vlanID)
		Expect(err).NotTo(HaveOccurred())
	})
})
