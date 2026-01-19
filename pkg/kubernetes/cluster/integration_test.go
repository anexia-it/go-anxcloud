//go:build integration
// +build integration

package cluster

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("cluster client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})
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
		It("serializes object", func() {
			cli, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var def Definition

				err := json.NewDecoder(r.Body).Decode(&def)
				Expect(err).NotTo(HaveOccurred())

				resp := struct {
					State struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type int    `json:"type"`
					} `json:"state"`
					Name string `json:"name"`
				}{
					State: struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type int    `json:"type"`
					}{
						ID:   "2",
						Text: "Pending",
						Type: 2,
					},
					Name: "clustername",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})

			cCreated, err := api.Create(context.TODO(), Definition{
				State: StateNewlyCreated,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("clustername"))
			Expect(cCreated.State.Type).To(Equal(2))
			Expect(cCreated.State.ID).To(Equal("2"))
			Expect(cCreated.State.Text).To(Equal("Pending"))
		})

		It("serializes string", func() {
			cli, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var def Definition

				err := json.NewDecoder(r.Body).Decode(&def)
				Expect(err).NotTo(HaveOccurred())

				resp := struct {
					State string `json:"state"`
					Name  string `json:"name"`
				}{
					State: "2",
					Name:  "clustername",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})

			cCreated, err := api.Create(context.TODO(), Definition{
				State: StateNewlyCreated,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("clustername"))
			Expect(cCreated.State.Type).To(Equal(2))
			Expect(cCreated.State.ID).To(Equal("2"))
			Expect(cCreated.State.Text).To(Equal("Unknown"))
		})
	})
})
