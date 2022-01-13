package disktype

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDiskTypeSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vsphere/provisioning/disktype API client suite")
}
