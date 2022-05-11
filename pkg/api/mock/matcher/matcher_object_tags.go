package matcher

import (
	"fmt"

	"github.com/onsi/gomega"
	gomegaTypes "github.com/onsi/gomega/types"
)

// TaggedWith checks if APIObject is tagged with tags
func TaggedWith(tags ...string) gomegaTypes.GomegaMatcher {
	return &objectTaggedWith{
		tags: tags,
	}
}

type objectTaggedWith struct {
	tags []string
}

func (m *objectTaggedWith) Match(actual interface{}) (bool, error) {
	tags, err := extractObjectTags(actual)
	if err != nil {
		return false, err
	}

	return gomega.ContainElements(m.tags).Match(tags)
}

func (m *objectTaggedWith) FailureMessage(actual interface{}) string {
	return m.commonFailureMessage(actual, false, "Object not tagged with all required tags")
}

func (m *objectTaggedWith) NegatedFailureMessage(actual interface{}) string {
	return m.commonFailureMessage(actual, true, "Object tagged with unexpected tags")
}

func (m *objectTaggedWith) commonFailureMessage(actual interface{}, negated bool, failureMessage string) string {
	tags, err := extractObjectTags(actual)
	if err != nil {
		return fmt.Sprintf("Failed to extract object tags: %s", err)
	}

	subMatch := gomega.ContainElements(m.tags)

	var subFailureMessage string
	if !negated {
		subFailureMessage = subMatch.FailureMessage(tags)
	} else {
		subFailureMessage = subMatch.NegatedFailureMessage(tags)
	}

	return fmt.Sprintf("%s: %s", failureMessage, subFailureMessage)
}

func extractObjectTags(actual interface{}) ([]string, error) {
	obj, err := castToAPIObjectPointer(actual)
	if err != nil {
		return nil, err
	}

	return obj.Tags(), nil
}
