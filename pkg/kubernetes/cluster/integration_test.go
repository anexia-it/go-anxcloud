package cluster

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
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

	Context("json serialization", func() {
		It("serializes correct", func() {
			cli, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var def Definition

				err := json.NewDecoder(r.Body).Decode(&def)
				Expect(err).NotTo(HaveOccurred())

				type c struct {
					State State  `json:"state"`
					Name  string `json:"name"`
				}

				err = json.NewEncoder(w).Encode(c{State{ID: "4", Text: "NewlyCreated", Type: gs.StateTypePending}, "clustername"})
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli, ClientOpts{Environment: EnvironmentDev})

			cCreated, err := api.Create(context.TODO(), Definition{
				State: StateNewlyCreated,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("clustername"))
		})
	})
})
