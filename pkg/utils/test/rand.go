package test

import "math/rand"

var random *rand.Rand = nil

func Seed(seed int64) {
	random = rand.New(rand.NewSource(seed))
}

func getRandom() *rand.Rand {
	if random == nil {
		panic("using pkg/utils/test.getRandom without seeding it first!")
	}

	return random
}
