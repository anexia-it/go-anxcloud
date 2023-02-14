//go:build integration
// +build integration

package ips

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "16854ecb42af4fad89f9fcef26789d50"
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
