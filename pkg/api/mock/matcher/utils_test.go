package matcher

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("custom gomega matcher utils", func() {
	Context("isMap", func() {
		It("tests if interface is of type Map", func() {
			Expect(isMap(map[string]string{"f": "oo", "b": "ar"})).To(BeTrue())
			Expect(isMap(nil)).To(BeFalse())
			Expect(isMap("not a map")).To(BeFalse())
		})
	})

	Context("isArrayOrSlice", func() {
		It("tests if interface is of type Array or Slice", func() {
			array := []int{1, 2, 3, 4, 5}
			Expect(isArrayOrSlice(array)).To(BeTrue())
			Expect(isArrayOrSlice(array[2:4])).To(BeTrue())
			Expect(isArrayOrSlice(nil)).To(BeFalse())
			Expect(isArrayOrSlice("not array or slice")).To(BeFalse())
		})
	})
})
