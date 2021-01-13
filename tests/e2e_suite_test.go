package tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "02f39d20ca0f4adfb5032f88dbc26c39"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests suite")
}
