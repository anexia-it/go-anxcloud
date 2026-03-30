//go:build integration
// +build integration

package storageserverinterface

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("storageserverinterface client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		api = NewAPI(cli)
	})

	Context("with a storageserverinterface", func() {
		var ssi StorageServerInterface

		It("lists storageserverinterfaces", func() {
			ssiInfos, err := api.Get(context.TODO(), 1, 20)
			Expect(err).NotTo(HaveOccurred())
			Expect(ssiInfos).NotTo(BeEmpty())

			ssi.Identifier = ssiInfos[0].Identifier
		})

		It("retrieves first ssi with expected values", func() {
			var err error
			ssi, err = api.GetByID(context.TODO(), ssi.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(ssi).NotTo(BeNil())
		})
	})

	Context("storageserverinterface full cycle", func() {
		var id string

		It("creates a storageserverinterface", func() {
			ssi, err := api.Create(context.TODO(), Definition{
				Name: "integration-test-ssi",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(ssi).NotTo(BeNil())
			id = ssi.Identifier
		})

		It("updates the storageserverinterface", func() {
			ssi, err := api.Update(context.TODO(), id, Definition{
				Name: "integration-test-ssi-updated",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(ssi).NotTo(BeNil())
		})

		It("deletes the storageserverinterface", func() {
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
					Name: "storageserverinterfacename",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli)

			cCreated, err := api.Create(context.TODO(), Definition{})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("storageserverinterfacename"))
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
					Name:  "storageserverinterfacename",
				}

				err = json.NewEncoder(w).Encode(resp)
				Expect(err).NotTo(HaveOccurred())
			}))
			defer server.Close()

			api := NewAPI(cli)

			cCreated, err := api.Create(context.TODO(), Definition{})

			Expect(err).NotTo(HaveOccurred())
			Expect(cCreated.Name).To(Equal("storageserverinterfacename"))
			Expect(cCreated.State.Type).To(Equal(2))
			Expect(cCreated.State.ID).To(Equal("2"))
			Expect(cCreated.State.Text).To(Equal("Unknown"))
		})
	})
})
