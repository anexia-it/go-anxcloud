package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega/format"
	gomegaTypes "github.com/onsi/gomega/types"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/compare"
)

// Object matches APIObject with provided Object
// It compares Objects by matchFields
func Object(obj types.Object, matchFields ...string) gomegaTypes.GomegaMatcher {
	return &objectMatcher{
		obj:         obj,
		matchFields: matchFields,
	}
}

type objectMatcher struct {
	obj         types.Object
	matchFields []string
}

var (
	// ErrObjectMatcherTypeMismatch is returned when ObjectMatcher is called with different Object types
	ErrObjectMatcherTypeMismatch = errors.New("Provided Objects have different types")
)

func (m *objectMatcher) Match(actual interface{}) (bool, error) {
	diff, err := m.compare(actual)
	if err == ErrObjectMatcherTypeMismatch {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return len(diff) == 0, nil
}

func (m *objectMatcher) FailureMessage(actual interface{}) string {
	diff, err := m.compare(actual)
	if err != nil {
		return err.Error()
	}

	formattedDiff := make([]string, 0, len(diff))
	for _, d := range diff {
		formattedDiff = append(formattedDiff, fmt.Sprintf(
			"Key %q\n%s\n",
			d.Key,
			format.Message(d.A, "to equal", d.B),
		))
	}

	return fmt.Sprintf("ACTUAL Object does not match EXPECTED Object:\n%s", strings.Join(formattedDiff, "\n"))
}

func (m *objectMatcher) NegatedFailureMessage(actual interface{}) string {
	return "ACTUAL Object does match EXPECTED Object"
}

func (m *objectMatcher) compare(actual interface{}) ([]compare.Difference, error) {
	obj, err := castToAPIObjectPointer(actual)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(m.obj) != reflect.TypeOf(obj.Unwrap()) {
		return nil, ErrObjectMatcherTypeMismatch
	}

	diff, err := compare.Compare(obj.Unwrap(), m.obj, m.matchFields...)
	if err != nil {
		return nil, err
	}

	return diff, nil
}
