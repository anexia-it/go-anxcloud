package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLBaaS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for LBaaS API definition")
}
