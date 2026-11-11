package prefix

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
			ghttp.VerifyRequest(http.MethodPut, "/api/ipam/v1/prefix.json/prefix-id"),
			ghttp.VerifyJSONRepresenting(map[string]interface{}{"description_customer": ""}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, Summary{ID: "prefix-id"}),
		))

		cli, err := client.New(client.BaseURL(server.URL()), client.IgnoreMissingToken())
		Expect(err).NotTo(HaveOccurred())

		_, err = NewAPI(cli).Update(context.Background(), "prefix-id", Update{CustomerDescription: pointer.String("")})
		Expect(err).NotTo(HaveOccurred())
	})

	It("leaves an unset customer description out of the request", func() {
		server := ghttp.NewServer()
		DeferCleanup(server.Close)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/api/ipam/v1/prefix.json/prefix-id"),
			ghttp.VerifyJSONRepresenting(map[string]interface{}{}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, Summary{ID: "prefix-id"}),
		))

		cli, err := client.New(client.BaseURL(server.URL()), client.IgnoreMissingToken())
		Expect(err).NotTo(HaveOccurred())

		_, err = NewAPI(cli).Update(context.Background(), "prefix-id", Update{})
		Expect(err).NotTo(HaveOccurred())
	})
})
