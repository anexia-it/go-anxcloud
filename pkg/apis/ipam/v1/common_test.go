package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressSpace JSON encoding", func() {
	var t AddressSpace

	Context("initialized to public", func() {
		BeforeEach(func() {
			t = AddressSpacePublic
		})

		It("encodes correctly", func() {
			d, err := t.MarshalJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(d).To(Equal([]byte(`0`)))
		})

		It("decodes correctly", func() {
			err := t.UnmarshalJSON([]byte(`1`))
			Expect(err).NotTo(HaveOccurred())
			Expect(t).To(Equal(AddressSpacePrivate))
		})
	})

	Context("initialized to private", func() {
		BeforeEach(func() {
			t = AddressSpacePrivate
		})

		It("encodes correctly", func() {
			d, err := t.MarshalJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(d).To(Equal([]byte(`1`)))
		})

		It("decodes correctly", func() {
			err := t.UnmarshalJSON([]byte(`0`))
			Expect(err).NotTo(HaveOccurred())
			Expect(t).To(Equal(AddressSpacePublic))
		})
	})
})
