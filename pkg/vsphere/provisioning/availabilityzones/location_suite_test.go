package availabilityzones

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLocationSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vsphere/provisioning/availabilityzones API client suite")
}
