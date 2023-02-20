package v1

import (
	"testing"

	testutil "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKubernetes(t *testing.T) {
	testutil.Seed(GinkgoRandomSeed())

	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for Kubernetes API definition")
}
