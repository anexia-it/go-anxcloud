package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KubeConfig", Ordered, func() {
	var (
		a   api.API
		srv *ghttp.Server

		clusterIdentifier = "mock-cluster-identifier"
	)

	BeforeEach(func() {
		srv = ghttp.NewServer()

		var err error
		a, err = api.NewAPI(api.WithClientOptions(
			client.BaseURL(srv.URL()),
			client.IgnoreMissingToken(),
		))
		Expect(err).ToNot(HaveOccurred())
	})

	It("can Request and Get kubeconfig", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": nil})
		appendRequestKubeConfigHandler(srv, clusterIdentifier)
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": nil})
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": "<kubeconfig>"})

		config, err := GetKubeConfig(context.TODO(), a, clusterIdentifier)
		Expect(err).ToNot(HaveOccurred())
		Expect(config).To(Equal("<kubeconfig>"))
	})

	It("can Get existing kubeconfig", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": "<kubeconfig>"})

		config, err := GetKubeConfig(context.TODO(), a, clusterIdentifier)
		Expect(err).ToNot(HaveOccurred())
		Expect(config).To(Equal("<kubeconfig>"))
	})

	It("supports Context cancelation", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": nil})
		appendRequestKubeConfigHandler(srv, clusterIdentifier)

		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()

		_, err := GetKubeConfig(ctx, a, clusterIdentifier)
		Expect(err).To(MatchError(context.DeadlineExceeded))
	})

	It("returns api.ErrNotFound when specified cluster does not exist", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusNotFound, nil)

		_, err := GetKubeConfig(context.TODO(), a, clusterIdentifier)
		Expect(err).To(MatchError(api.ErrNotFound))
	})

	It("returns an error when RequestKubeConfig fails", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"kubeconfig": nil})

		srv.AllowUnhandledRequests = true
		srv.UnhandledRequestStatusCode = 500

		_, err := GetKubeConfig(context.TODO(), a, clusterIdentifier)
		Expect(err).To(HaveOccurred())
		var he api.HTTPError
		Expect(errors.As(err, &he)).To(BeTrue())
	})
})
