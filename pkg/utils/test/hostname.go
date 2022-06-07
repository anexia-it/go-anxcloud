package test

import (
	"math/rand"
	"time"
)

const (
	hostnameCharset      = "abcdefghijklmnopqrstuvwxyz"
	randomHostnameLength = 8
)

// RandomHostname generates a random hostname, useful for e2e tests and primarily used there.
func RandomHostname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // No crypto needed here.
	hostname := make([]byte, randomHostnameLength)
	for i := range hostname {
		hostname[i] = hostnameCharset[r.Intn(len(hostnameCharset))]
	}

	return string(hostname)
}
