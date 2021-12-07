package progress

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProgressSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vsphere/provisioning/progress API client suite")
}
