package matcher

import (
	"fmt"

	gomegaTypes "github.com/onsi/gomega/types"
	"go.anx.io/go-anxcloud/pkg/api/mock"
)

// Created checks if *APIObject was created n times (ignoring FakeExisting)
func Created(n int) gomegaTypes.GomegaMatcher {
	return &crudMatcher{
		n:  n,
		op: "created",
		validator: func(o *mock.APIObject) bool {
			return o.CreatedCount() == n
		},
	}
}

// Destroyed checks if *APIObject was destroyed n times
func Destroyed(n int) gomegaTypes.GomegaMatcher {
	return &crudMatcher{
		n:  n,
		op: "destroyed",
		validator: func(o *mock.APIObject) bool {
			return o.DestroyedCount() == n
		},
	}
}

// Updated checks if *APIObject was updated n times
func Updated(n int) gomegaTypes.GomegaMatcher {
	return &crudMatcher{
		n:  n,
		op: "updated",
		validator: func(o *mock.APIObject) bool {
			return o.UpdatedCount() == n
		},
	}
}

type crudMatcher struct {
	validator func(*mock.APIObject) bool
	op        string
	n         int
}

func (m *crudMatcher) Match(actual interface{}) (bool, error) {
	obj, err := castToAPIObjectPointer(actual)
	if err != nil {
		return false, err
	}

	return m.validator(obj), nil
}

func (m *crudMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Object has not been %s <%d> times", m.op, m.n)
}

func (m *crudMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Object has been %s <%d> times", m.op, m.n)
}
