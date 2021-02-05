package tests

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)
var _ = Describe("CloudDNS API endpoint tests", func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Zone Endpoint", func() {
		It("Should list all avaliable zones", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).List(ctx)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
