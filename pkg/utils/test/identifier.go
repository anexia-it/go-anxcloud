package test

import (
	"math/rand"
	"time"
)

const (
	identifierCharset      = "abcdef0123456789"
	randomIdentifierLength = 32
)

// generates random identifier similar to engine identifiers
func RandomIdentifier() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // No crypto needed here.
	identifier := make([]byte, randomIdentifierLength)
	for i := range identifier {
		identifier[i] = identifierCharset[r.Intn(len(identifierCharset))]
	}

	return string(identifier)
}
