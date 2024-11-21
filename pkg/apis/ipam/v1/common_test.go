package v1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
)

type TypeMarshallingTest struct {
	Type               ipamv1.AddressType
	JsonRepresentation string
}

var _ = DescribeTableSubtree("Type JSON encoding",
	func(t TypeMarshallingTest) {
		It("encodes correctly", func() {
			res, err := t.Type.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal([]byte(t.JsonRepresentation)))
		})
		It("decodes correctly", func() {
			var res ipamv1.AddressType
			Expect(res.UnmarshalJSON([]byte(t.JsonRepresentation))).To(Succeed())
			Expect(res).To(Equal(t.Type))
		})
	},
	Entry("private", TypeMarshallingTest{ipamv1.TypePrivate, "1"}),
	Entry("private", TypeMarshallingTest{ipamv1.TypePublic, "0"}),
	Entry("zero value", TypeMarshallingTest{"", ""}),
)
