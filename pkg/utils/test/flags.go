package test

import "flag"

var RunAsIntegrationTest = false

func InitFlags() {
	flag.BoolVar(&RunAsIntegrationTest, "integration-test", RunAsIntegrationTest, "Run the test suite as integration tests against the real Engine API")
}
