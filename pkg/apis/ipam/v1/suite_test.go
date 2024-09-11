package v1

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testutil "go.anx.io/go-anxcloud/pkg/utils/test"
)

func TestIPAM(t *testing.T) {
	testutil.Seed(GinkgoRandomSeed())

	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for IPAM API definition")
}
