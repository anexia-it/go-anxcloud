package test

import "math/rand" // #nosec G404 - math/rand is acceptable for test utilities

var random *rand.Rand = nil

func Seed(seed int64) {
	random = rand.New(rand.NewSource(seed)) // #nosec G404 - test utility doesn't need cryptographic randomness
}

func getRandom() *rand.Rand {
	if random == nil {
		panic("using pkg/utils/test.getRandom without seeding it first!")
	}

	return random
}
