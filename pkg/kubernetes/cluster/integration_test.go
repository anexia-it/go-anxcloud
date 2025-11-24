//go:build integration
// +build integration

package cluster

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("cluster client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli, ClientOpts{Environment: EnvironmentDev})
	})

	Context("with a cluster", func() {
		var cluster Cluster

		It("lists clusters", func() {
			clusterInfos, err := api.Get(context.TODO(), 1, 20)
			Expect(err).NotTo(HaveOccurred())
			Expect(clusterInfos).NotTo(BeEmpty())

			cluster.Identifier = clusterInfos[0].Identifier
		})

		It("retrieves first cluster with expected values", func() {
			var err error
			cluster, err = api.GetByID(context.TODO(), cluster.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(cluster).NotTo(BeNil())
		})

		It("can trigger GA request kubeconfig", func() {
			err := api.RequestKubeConfig(context.TODO(), &cluster)

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
