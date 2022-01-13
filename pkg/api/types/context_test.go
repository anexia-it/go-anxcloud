package types

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("context key retriever functions", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.TODO()
	})

	Context("with no attributes set", func() {
		It("returns ErrContextKeyNotSet error for URL", func() {
			u, err := URLFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(u).To(BeZero())
		})

		It("returns ErrContextKeyNotSet error for Operation", func() {
			o, err := OperationFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(o).To(BeZero())
		})

		It("returns ErrContextKeyNotSet error for Options", func() {
			o, err := OptionsFromContext(ctx)
			Expect(err).To(MatchError(ErrContextKeyNotSet))
			Expect(o).To(BeNil())
		})
	})
})
