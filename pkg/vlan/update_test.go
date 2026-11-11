package vlan

import (
	"context"
	"net/http"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", func() {
	It("sends an empty customer description", func() {
		server := ghttp.NewServer()
		DeferCleanup(server.Close)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/api/vlan/v1/vlan.json/vlan-id"),
			ghttp.VerifyJSONRepresenting(map[string]interface{}{"description_customer": ""}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, Summary{}),
		))

		cli, err := client.New(client.BaseURL(server.URL()), client.IgnoreMissingToken())
		Expect(err).NotTo(HaveOccurred())

		err = NewAPI(cli).Update(context.Background(), "vlan-id", UpdateDefinition{CustomerDescription: pointer.String("")})
		Expect(err).NotTo(HaveOccurred())
	})

	It("leaves unset fields out of the request", func() {
		server := ghttp.NewServer()
		DeferCleanup(server.Close)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/api/vlan/v1/vlan.json/vlan-id"),
			ghttp.VerifyJSONRepresenting(map[string]interface{}{}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, Summary{}),
		))

		cli, err := client.New(client.BaseURL(server.URL()), client.IgnoreMissingToken())
		Expect(err).NotTo(HaveOccurred())

		err = NewAPI(cli).Update(context.Background(), "vlan-id", UpdateDefinition{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("sends explicitly disabled VM provisioning", func() {
		server := ghttp.NewServer()
		DeferCleanup(server.Close)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/api/vlan/v1/vlan.json/vlan-id"),
			ghttp.VerifyJSONRepresenting(map[string]interface{}{"vm_provisioning": false}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, Summary{}),
		))

		cli, err := client.New(client.BaseURL(server.URL()), client.IgnoreMissingToken())
		Expect(err).NotTo(HaveOccurred())

		err = NewAPI(cli).Update(context.Background(), "vlan-id", UpdateDefinition{VMProvisioning: pointer.Bool(false)})
		Expect(err).NotTo(HaveOccurred())
	})
})
