package filter

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("filter.Helper", func() {
	var o types.Object

	BeforeEach(func() {
		o = &testObject{}
	})

	It("returns the correct error when trying to retrieve an unknown field", func() {
		helper, err := NewHelper(o)
		Expect(err).NotTo(HaveOccurred())

		_, ok, err := helper.Get("fooooNotAConfiguredFilter")
		Expect(ok).NotTo(BeTrue())
		Expect(err).To(MatchError(ErrUnknownField))
	})

	Context("with no filters set", func() {
		It("gives an empty query string", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			query := helper.BuildQuery().Encode()
			Expect(query).To(Equal(""))
		})

		It("returns not-present for every filterable field", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			for field := range helper.(filterHelper).fields {
				_, ok, err := helper.Get(field)
				Expect(ok).To(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Context("with identifier set in referenced Object", func() {
		JustBeforeEach(func() {
			o.(*testObject).Parent.Identifier = "foo"
		})

		It("gives the expected query string", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			query := helper.BuildQuery().Encode()
			Expect(query).To(Equal("parent=foo"))
		})

		It("returns the filter value and it being set", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			val, present, err := helper.Get("parent")
			Expect(err).NotTo(HaveOccurred())
			Expect(present).To(BeTrue())
			Expect(val).To(Equal("foo"))
		})
	})

	Context("with identifier set in pointer-referenced Object", func() {
		JustBeforeEach(func() {
			o.(*testObject).PointerParent = &parentObject{
				Identifier: "pointerfoo",
			}
		})

		It("gives the expected query string", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			query := helper.BuildQuery().Encode()
			Expect(query).To(Equal("ptrParent=pointerfoo"))
		})

		It("returns the filter value and it being set", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			val, present, err := helper.Get("ptrParent")
			Expect(err).NotTo(HaveOccurred())
			Expect(present).To(BeTrue())
			Expect(val).To(Equal("pointerfoo"))
		})
	})

	Context("with description set", func() {
		JustBeforeEach(func() {
			desc := "some random description"
			o.(*testObject).Description = &desc
		})

		It("gives the expected query string", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			query := helper.BuildQuery().Encode()
			Expect(query).To(Equal("desc=some+random+description"))
		})

		It("returns the filter value and it being set", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			val, present, err := helper.Get("desc")
			Expect(err).NotTo(HaveOccurred())
			Expect(present).To(BeTrue())
			Expect(val).To(Equal("some random description"))
		})
	})
})

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "filter test suite")
}
