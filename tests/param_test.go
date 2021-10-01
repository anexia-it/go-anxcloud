package tests_test

import (
	"github.com/anexia-it/go-anxcloud/pkg/utils/param"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Parameter Builder Tests", func() {
	It("Should Create Parameter", func() {
		const testKey = "testKey"
		const testValue = "testValue"
		builder := param.ParameterBuilder(testKey)
		parameter := builder(testValue)

		values := url.Values{}
		parameter(values)
		Expect(values.Get(testKey)).To(BeEquivalentTo(testValue))
	})
})
