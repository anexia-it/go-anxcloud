package param

import (
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParameterBuilder", func() {
	It("should Create Parameter", func() {
		const testKey = "testKey"
		const testValue = "testValue"
		builder := ParameterBuilder(testKey)
		parameter := builder(testValue)

		values := url.Values{}
		parameter(values)
		Expect(values.Get(testKey)).To(BeEquivalentTo(testValue))
	})
})

func TestParam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Param test suite")
}
