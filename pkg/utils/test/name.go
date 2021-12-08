package test

import "fmt"

func TestResourceName() string {
	return fmt.Sprintf("go-test-%s", RandomHostname())
}
