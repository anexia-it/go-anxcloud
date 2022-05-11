package mock_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMockAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock API")
}
