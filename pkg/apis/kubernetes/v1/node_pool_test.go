package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodePool resource", func() {
	It("can filter by cluster", func() {
		np := &NodePool{Cluster: Cluster{Identifier: "foo"}}
		url, err := np.EndpointURL(
			types.ContextWithOperation(types.ContextWithOptions(context.TODO(), &types.ListOptions{}), types.OperationList))
		Expect(err).ToNot(HaveOccurred())
		Expect(url.Query().Encode()).To(Equal("filters=cluster%3Dfoo"))
	})
})
