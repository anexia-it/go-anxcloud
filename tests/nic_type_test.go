package tests_test

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/nictype"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NIC type API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("NIC type endpoint", func() {

		It("Should list all available NIC types", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			nicTypes, err := nictype.NewAPI(cli).List(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nicTypes)).To(BeNumerically(">", 1))
		})

	})

})
