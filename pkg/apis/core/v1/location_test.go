package v1_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

var _ = Describe("Location Object", func() {
	var o *corev1.Location

	BeforeEach(func() {
		o = &corev1.Location{}
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
