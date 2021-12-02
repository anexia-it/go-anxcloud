package nictype

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNICTypeSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NIC type test suite")
}
