package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"

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

var _ = Describe("Create Cluster", func() {
	DescribeTable("prefix configurations", func(cluster Cluster, expected error) {
		_, err := cluster.FilterAPIRequestBody(types.ContextWithOperation(context.TODO(), types.OperationCreate))
		if expected == nil {
			Expect(err).ToNot(HaveOccurred())
		} else {
			Expect(err).To(MatchError(expected))
		}
	},
		Entry("all prefixes implicitly managed", Cluster{}, nil),
		Entry("all prefixes explicitly managed", Cluster{ManageInternalIPv4Prefix: pointer.Bool(true), ManageExternalIPv4Prefix: pointer.Bool(true), ManageExternalIPv6Prefix: pointer.Bool(true)}, nil),
		Entry("internal v4 prefix implicitly managed and explicitly provided", Cluster{InternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("internal v4 prefix explicitly managed and explicitly provided", Cluster{ManageInternalIPv4Prefix: pointer.Bool(true), InternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("internal v4 prefix explicitly unmanaged and provided", Cluster{ManageInternalIPv4Prefix: pointer.Bool(false), InternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, nil),
		Entry("external v4 prefix implicitly managed and explicitly provided", Cluster{ExternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("external v4 prefix explicitly managed and explicitly provided", Cluster{ManageExternalIPv4Prefix: pointer.Bool(true), ExternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("external v4 prefix explicitly unmanaged and provided", Cluster{ManageExternalIPv4Prefix: pointer.Bool(false), ExternalIPv4Prefix: &common.PartialResource{Identifier: "foo"}}, nil),
		Entry("external v6 prefix implicitly managed and explicitly provided", Cluster{ExternalIPv6Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("external v6 prefix explicitly managed and explicitly provided", Cluster{ManageExternalIPv6Prefix: pointer.Bool(true), ExternalIPv6Prefix: &common.PartialResource{Identifier: "foo"}}, ErrManagedPrefixSet),
		Entry("external v6 prefix explicitly unmanaged and provided", Cluster{ManageExternalIPv6Prefix: pointer.Bool(false), ExternalIPv6Prefix: &common.PartialResource{Identifier: "foo"}}, nil),
	)
})

var _ = DescribeTable("explicitlyFalse",
	func(b *bool, expected bool) {
		Expect(explicitlyFalse(b)).To(Equal(expected))
	},
	Entry("explicitly false", pointer.Bool(false), true),
	Entry("explicitly true", pointer.Bool(true), false),
	Entry("implicitly false (nil)", nil, false),
)
