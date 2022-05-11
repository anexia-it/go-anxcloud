package matcher

import (
	gomegaTypes "github.com/onsi/gomega/types"
)

// Existing checks if APIObject does currently exist
func Existing() gomegaTypes.GomegaMatcher {
	return &objectExisting{}
}

type objectExisting struct{}

func (m *objectExisting) Match(actual interface{}) (bool, error) {
	obj, err := castToAPIObjectPointer(actual)
	if err != nil {
		return false, err
	}

	return obj.Existing(), nil
}

func (m *objectExisting) FailureMessage(actual interface{}) string {
	return "Object does not exist"
}

func (m *objectExisting) NegatedFailureMessage(actual interface{}) string {
	return "Object does exist"
}
