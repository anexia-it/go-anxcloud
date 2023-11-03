package v2

import (
	"context"
	"encoding/json"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"

	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLBaaSv2APIBindings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LBaaS v2 API Bindings Suite")
}

var _ = Describe("mock", func() {
	var (
		engine api.API
		srv    *ghttp.Server
	)

	BeforeEach(func() {
		var err error
		srv = ghttp.NewServer()
		engine, err = api.NewAPI(
			api.WithClientOptions(
				client.BaseURL(srv.URL()),
				client.IgnoreMissingToken(),
			),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		srv.Close()
	})

	Context("Cluster", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/clusters.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			c := Cluster{}
			err := engine.Create(context.TODO(), &c)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/clusters.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"implementation": "haproxy",
					"frontend_prefixes": "foo,bar",
					"backend_prefixes": "bar,foo",
					"replicas": 3
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			c := Cluster{
				Name:             "foo",
				Implementation:   LoadBalancerImplementationHAProxy,
				FrontendPrefixes: &gs.PartialResourceList{{Identifier: "foo"}, {Identifier: "bar"}},
				BackendPrefixes:  &gs.PartialResourceList{{Identifier: "bar"}, {Identifier: "foo"}},
				Replicas:         pointer.Int(3),
			}
			err := engine.Create(context.TODO(), &c)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Context("Node", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/nodes.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			n := Node{}
			err := engine.Create(context.TODO(), &n)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/nodes.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"cluster": "bar"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			n := Node{
				Name:    "foo",
				Cluster: &common.PartialResource{Identifier: "bar"},
			}
			err := engine.Create(context.TODO(), &n)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly applies filter parameters on List operations", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/LBaaSv2/v1/nodes.json", "limit=10&page=1&filters=cluster=bar"),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			n := Node{Cluster: &common.PartialResource{Identifier: "bar"}}
			var channel types.ObjectChannel
			err := engine.List(context.TODO(), &n, api.ObjectChannel(&channel))
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Context("LoadBalancer", func() {
		It("correctly encodes zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/load_balancers.json"),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			lb := LoadBalancer{}
			err := engine.Create(context.TODO(), &lb)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly encodes non-zero values", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/LBaaSv2/v1/load_balancers.json"),
				ghttp.VerifyJSON(`{
					"name": "foo",
					"generation": 5,
					"cluster": "bar",
					"frontend_ips": "foo,bar",
					"ssl_certificates": "bar,foo",
					"definition": "{\"frontends\":[{\"name\":\"foo\",\"protocol\":\"TCP\",\"backend\":{\"protocol\":\"TCP\"}}]}"
				}`),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			lb := LoadBalancer{
				Name:            "foo",
				Generation:      5,
				Cluster:         &common.PartialResource{Identifier: "bar"},
				FrontendIPs:     &gs.PartialResourceList{{Identifier: "foo"}, {Identifier: "bar"}},
				SSLCertificates: &gs.PartialResourceList{{Identifier: "bar"}, {Identifier: "foo"}},
				Definition:      &Definition{Frontends: []Frontend{{Name: "foo", Protocol: "TCP", Backend: FrontendBackend{Protocol: "TCP"}}}},
			}
			err := engine.Create(context.TODO(), &lb)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly applies filter parameters on List operations", func() {
			srv.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/LBaaSv2/v1/load_balancers.json", "limit=10&page=1&filters=cluster=bar"),
				ghttp.RespondWithJSONEncoded(200, map[string]any{}),
			))
			lb := LoadBalancer{Cluster: &common.PartialResource{Identifier: "bar"}}
			var channel types.ObjectChannel
			err := engine.List(context.TODO(), &lb, api.ObjectChannel(&channel))
			Expect(err).ToNot(HaveOccurred())
		})

		It("supports decoding nested json (Definition)", func() {
			lbJson := `{"definition": "{\"frontends\":[{\"name\":\"foo\",\"protocol\":\"TCP\",\"backend\":{\"protocol\":\"TCP\"}}]}"}`
			var lb LoadBalancer
			err := json.Unmarshal([]byte(lbJson), &lb)
			Expect(err).ToNot(HaveOccurred())
			Expect(*lb.Definition).To(Equal(Definition{Frontends: []Frontend{{Name: "foo", Protocol: "TCP", Backend: FrontendBackend{Protocol: "TCP"}}}}))
		})
	})
})
