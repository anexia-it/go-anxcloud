package v1_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLBaaS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for LBaaS API definition")
}
