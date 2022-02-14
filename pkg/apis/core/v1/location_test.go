package v1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Location Object", func() {
	var o *Location

	BeforeEach(func() {
		o = &Location{}
	})

	DescribeTable("gives ErrOperationNotSupported",
		func(op types.Operation) {
			_, err := o.EndpointURL(types.ContextWithOperation(context.TODO(), op))
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
		},
		Entry("for Create operation", types.OperationCreate),
		Entry("for Update operation", types.OperationUpdate),
		Entry("for Destroy operation", types.OperationDestroy),
	)
})
