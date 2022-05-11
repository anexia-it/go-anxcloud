package mock

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Mock iterator", func() {
	It("handles error for zero length page size", func() {
		_, err := newMockPageIter([]types.Object{}, 0, 1)
		Expect(err).To(MatchError(ErrPageSizeCannotBeZero))
	})

	It("supports resetting error", func() {
		iter, err := newMockPageIter([]types.Object{}, 1, 1)
		Expect(err).ToNot(HaveOccurred())
		iter.(*mockPageIter).err = errors.New("test")
		Expect(iter.Error()).To(HaveOccurred())
		iter.ResetError()
		Expect(iter.Error()).ToNot(HaveOccurred())
	})
})
