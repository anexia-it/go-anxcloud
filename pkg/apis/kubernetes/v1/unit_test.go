//go:build !integration
// +build !integration

package v1

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api/types"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"go.anx.io/go-anxcloud/pkg/client"
)

const isIntegrationTest = false

var _ = Describe("NodePool resource", func() {
	It("can filter by cluster", func() {
		np := &NodePool{Cluster: Cluster{Identifier: "foo"}}
		url, err := np.EndpointURL(types.ContextWithOperation(context.TODO(), types.OperationList))
		Expect(err).ToNot(HaveOccurred())
		Expect(url.Query().Encode()).To(Equal("filters=cluster%3Dfoo"))
	})
})

var mockStateError = map[string]interface{}{"type": gs.StateTypeError}
var mockStatePending = map[string]interface{}{"type": gs.StateTypePending}

var _ = Describe("AwaitCompletion", Ordered, func() {
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

	It("can wait until state is ready", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateOK})

		cluster := Cluster{Identifier: clusterIdentifier}
		err := gs.AwaitCompletion(context.TODO(), a, &cluster)
		Expect(err).ToNot(HaveOccurred())
		Expect(cluster.StateOK()).To(BeTrue())
		Expect(srv.ReceivedRequests()).To(HaveLen(3))
	})

	It("returns an error when updating the state fails", func() {
		cluster := Cluster{}
		err := gs.AwaitCompletion(context.TODO(), a, &cluster)
		Expect(err).To(MatchError(types.ErrUnidentifiedObject))
	})

	It("returns ErrClusterProvisioning when state is Error", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateError})

		cluster := Cluster{Identifier: clusterIdentifier}
		err := gs.AwaitCompletion(context.TODO(), a, &cluster)
		Expect(err).To(MatchError(gs.ErrStateError))
		Expect(cluster.StateError()).To(BeTrue())
		Expect(srv.ReceivedRequests()).To(HaveLen(2))
	})

	It("supports Context cancelation", func() {
		appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		cluster := Cluster{Identifier: clusterIdentifier}
		err := gs.AwaitCompletion(ctx, a, &cluster)
		Expect(err).To(MatchError(context.DeadlineExceeded))
	})
})
