package cpuperformancetype

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCPUPerformanceTypeSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vsphere/provisioning/cpuperformancetypes API client suite")
}
