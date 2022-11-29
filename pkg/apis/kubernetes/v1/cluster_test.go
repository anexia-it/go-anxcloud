package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
