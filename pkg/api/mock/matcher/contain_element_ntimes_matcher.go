package matcher

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

var (
	// ErrContainElementNTimesInvalidArgument is returned when ContainElementNTimes ACTUAL is not array/slice/map
	ErrContainElementNTimesInvalidArgument = errors.New("ContainElementNTimes matcher expects an array/slice/map")
)

// ContainElementNTimes succeeds if actual contains the given element n-times.
// This matcher is a slightly modified version of gomegas own ContainElement matcher.
func ContainElementNTimes(element interface{}, n int) types.GomegaMatcher {
	return &containElementNTimesMatcher{
		element: element,
		n:       n,
	}
}

type containElementNTimesMatcher struct {
	element interface{}
	n       int
}

func (matcher *containElementNTimesMatcher) Match(actual interface{}) (success bool, err error) {
	if !isArrayOrSlice(actual) && !isMap(actual) {
		return false, fmt.Errorf("%w. Got:\n%s", ErrContainElementNTimesInvalidArgument, format.Object(actual, 1))
	}

	elemMatcher, elementIsMatcher := matcher.element.(types.GomegaMatcher)
	if !elementIsMatcher {
		elemMatcher = &matchers.EqualMatcher{Expected: matcher.element}
	}

	value := reflect.ValueOf(actual)
	var valueAt func(int) interface{}
	if isMap(actual) {
		keys := value.MapKeys()
		valueAt = func(i int) interface{} {
			return value.MapIndex(keys[i]).Interface()
		}
	} else {
		valueAt = func(i int) interface{} {
			return value.Index(i).Interface()
		}
	}

	successCount := 0

	var lastError error
	for i := 0; i < value.Len(); i++ {
		success, err := elemMatcher.Match(valueAt(i))
		if err != nil {
			lastError = err
			continue
		}
		if success {
			successCount++
		}
	}

	return successCount == matcher.n, lastError
}

func (matcher *containElementNTimesMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain n elements matching", matcher.element)
}

func (matcher *containElementNTimesMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain n elements matching", matcher.element)
}
