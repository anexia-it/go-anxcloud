package v2_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestObjectStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ObjectStorage v2 Suite")
}
