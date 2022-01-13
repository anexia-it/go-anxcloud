package vsphere

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// This versioning scheme that currently seems to be in place for template build numbers.
var buildNumberRegex = regexp.MustCompile(`[bB]?(\d+)`)

func extractBuildNumber(build string) int {
	matches := buildNumberRegex.FindStringSubmatch(build)
	if len(matches) != 2 {
		// panic here since someone needs to check on the regex
		panic("build does not match the buildNumberRegex")
	}

	number, err := strconv.ParseInt(matches[1], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("could not extract build for %s: %s", build, err.Error()))
	}
	return int(number)
}

var _ = Describe("extractBuildNumber()", func() {
	It("extracts build number from string", func() {
		Expect(extractBuildNumber("b5555")).To(BeEquivalentTo(5555))
		Expect(extractBuildNumber("B111")).To(BeEquivalentTo(111))
		Expect(extractBuildNumber("123")).To(BeEquivalentTo(123))
	})
})
