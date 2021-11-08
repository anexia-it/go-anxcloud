package tests_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/utils/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "02f39d20ca0f4adfb5032f88dbc26c39"
)

func init() {
	test.RunAsIntegrationTest = true
	test.InitFlags()

	vsphereTestInit()
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	rand.Seed(time.Now().Unix())
	RunSpecs(t, "Tests suite")
}
