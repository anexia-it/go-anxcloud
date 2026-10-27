package v1_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	clouddnsv1 "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	"go.anx.io/go-anxcloud/pkg/client"
)

var _ = Describe("Record ZoneName in the request path", func() {
	var (
		srv *ghttp.Server
		a   api.API
	)

	BeforeEach(func() {
		srv = ghttp.NewServer()
		DeferCleanup(srv.Close)

		var err error
		a, err = api.NewAPI(api.WithClientOptions(client.BaseURL(srv.URL()), client.IgnoreMissingToken()))
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("sends the zone name as one escaped path segment and keeps it on listed records",
		func(zoneName, rawPath string) {
			srv.AppendHandlers(ghttp.CombineHandlers(
				func(_ http.ResponseWriter, req *http.Request) {
					Expect(req.URL.EscapedPath()).To(Equal(rawPath))
					Expect(req.URL.RawQuery).To(BeEmpty())
				},
				ghttp.RespondWithJSONEncoded(200, []clouddnsv1.Record{{Identifier: "r1", Name: "www", Type: "A", RData: "10.0.0.1"}}),
			))

			channel := make(types.ObjectChannel)
			Expect(a.List(context.TODO(), &clouddnsv1.Record{ZoneName: zoneName}, api.ObjectChannel(&channel))).To(Succeed())

			var r clouddnsv1.Record
			for res := range channel {
				Expect(res(&r)).To(Succeed())
			}

			Expect(r.Identifier).To(Equal("r1"))
			Expect(r.ZoneName).To(Equal(zoneName))
		},
		Entry("plain name", "example.com", "/api/clouddns/v1/zone.json/example.com/records"),
		Entry("question mark", "a?b", "/api/clouddns/v1/zone.json/a%3Fb/records"),
		Entry("percent sign", "a%2Fb", "/api/clouddns/v1/zone.json/a%252Fb/records"),
	)

	It("rejects a zone name containing a slash before sending anything", func() {
		err := a.Destroy(context.TODO(), &clouddnsv1.Record{ZoneName: "a/b", Identifier: "r1"})
		Expect(err).To(MatchError(clouddnsv1.ErrInvalidZoneName))
		Expect(srv.ReceivedRequests()).To(BeEmpty())
	})
})
