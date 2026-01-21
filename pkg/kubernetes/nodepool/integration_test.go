//go:build integration
// +build integration

package nodepool

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("nodepool client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})
	})

	Context("with a nodepool", func() {
		var nodepool Nodepool

		It("lists nodepools", func() {
			nodepoolInfos, err := api.Get(context.TODO(), 1, 20)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodepoolInfos).NotTo(BeEmpty())

			nodepool.Identifier = nodepoolInfos[0].Identifier
		})

		It("retrieves first nodepool with expected values", func() {
			var err error
			nodepool, err = api.GetByID(context.TODO(), nodepool.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(nodepool).NotTo(BeNil())
		})
	})

	Context("nodepool full cycle", func() {
		var id string

		It("creates a nodepool", func() {
			nodepool, err := api.Create(context.TODO(), Definition{
				Name:            "integration-test-nodepool",
				State:           StatePending,
				ClusterID:       "d0d9c6f26dd2489281e0bd86d79de572",
				Replicas:        4,
				CPUs:            3,
				MemoryBytes:     5 * Gibibyte,
				DiskSizeBytes:   22 * Gibibyte,
				OperatingSystem: OSFlatcar,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(nodepool).NotTo(BeNil())
			id = nodepool.Identifier
		})

		It("updates the nodepool", func() {
			nodepool, err := api.Update(context.TODO(), id, Definition{
				Name:            "integration-test-nodepool-updated",
				ClusterID:       "d0d9c6f26dd2489281e0bd86d79de572",
				Replicas:        5,
				CPUs:            4,
				MemoryBytes:     6 * Gibibyte,
				DiskSizeBytes:   23 * Gibibyte,
				OperatingSystem: OSFlatcar,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(nodepool).NotTo(BeNil())
		})

		It("deletes the nodepool", func() {
			err := api.DeleteByID(context.TODO(), id)
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
					Name: "nodepoolname",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})

			cCreated, err := api.Create(context.TODO(), Definition{
				State: StatePending,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("nodepoolname"))
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
					Name:  "nodepoolname",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli, common.ClientOpts{Environment: common.EnvironmentDev})

			cCreated, err := api.Create(context.TODO(), Definition{
				State: StatePending,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("nodepoolname"))
			Expect(cCreated.State.Type).To(Equal(2))
			Expect(cCreated.State.ID).To(Equal("2"))
			Expect(cCreated.State.Text).To(Equal("Unknown"))
		})
	})
})
