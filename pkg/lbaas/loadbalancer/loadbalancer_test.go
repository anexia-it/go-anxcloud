package loadbalancer

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLoadbalancer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for Loadbalancer")
}
