package test

import (
	"fmt"
)

// TestResourceName generates a name suitable for e2e test resources.
func TestResourceName() string {
	return fmt.Sprintf("go-test-%s", RandomHostname())
}
