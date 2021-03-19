package tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "883003c636644f08a5356078b5fc189d"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests suite")
}
