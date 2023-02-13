package gs

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommonGenericServiceAPIResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for common gs API resources")
}

var _ = DescribeTable("PartialResourceList",
	func(list PartialResourceList, expected string) {
		data, err := json.Marshal(&list)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal(expected))
	},
	Entry("empty list", PartialResourceList{}, `""`),
	Entry("single resource", PartialResourceList{{Identifier: "foo"}}, `"foo"`),
	Entry("multiple resources", PartialResourceList{{Identifier: "foo"}, {Identifier: "bar"}}, `"foo,bar"`),
)
