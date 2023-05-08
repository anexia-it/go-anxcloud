package vm

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM definition", func() {
	var api API

	BeforeEach(func() {
		api = NewAPI(nil)
	})

	It("should ignore no dns servers", func() {
		definition := api.NewDefinitionWithDNS("foo", "template", "44b38284-6adb-430e-b4a4-1553e29f352f", "developersfirstvm", 2, 2048, 10, []Network{}, []string{})

		Expect(definition.DNS1).To(BeEmpty())
		Expect(definition.DNS2).To(BeEmpty())
		Expect(definition.DNS3).To(BeEmpty())
		Expect(definition.DNS4).To(BeEmpty())
	})

	It("should ignore empty dns server", func() {
		definition := api.NewDefinitionWithDNS("foo", "template", "44b38284-6adb-430e-b4a4-1553e29f352f", "developersfirstvm", 2, 2048, 10, []Network{}, []string{"dns1", "", "dns3", "dns4"})

		Expect(definition.DNS1).To(Equal("dns1"))
		Expect(definition.DNS2).To(BeEmpty())
		Expect(definition.DNS3).To(Equal("dns3"))
		Expect(definition.DNS4).To(Equal("dns4"))
	})

	It("should only assign first four dns servers", func() {
		definition := api.NewDefinitionWithDNS("foo", "template", "44b38284-6adb-430e-b4a4-1553e29f352f", "developersfirstvm", 2, 2048, 10, []Network{}, []string{"dns1", "dns2", "dns3", "dns4", "dns5"})

		Expect(definition.DNS1).To(Equal("dns1"))
		Expect(definition.DNS2).To(Equal("dns2"))
		Expect(definition.DNS3).To(Equal("dns3"))
		Expect(definition.DNS4).To(Equal("dns4"))
	})
})

func TestProvisioningVMUnits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for pkg/vsphere/provisioning/vm")
}
