package common

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommonAPIResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for common API resources")
}

var _ = Describe("PartialResource", func() {
	It("json marshalls to the resources identifier", func() {
		pr := PartialResource{Identifier: "foo", Name: "bar"}
		data, err := json.Marshal(&pr)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal(`"foo"`))
	})
})
