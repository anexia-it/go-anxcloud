package ips

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIPsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vsphere/provisioning/ips API client suite")
}
