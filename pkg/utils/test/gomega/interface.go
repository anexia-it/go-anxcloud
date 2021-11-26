package test

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"
)

type implementInterfaceMatcher struct {
	iface interface{}
}

// ImplementInterface succeeds if actual implements the given interface.
//
// The expected interface has to be passed as a pointer to variable of it, see the provided example.
func ImplementInterface(i interface{}) types.GomegaMatcher {
	return implementInterfaceMatcher{
		iface: i,
	}
}

func (ii implementInterfaceMatcher) Match(actual interface{}) (bool, error) {
	ifaceType := reflect.TypeOf(ii.iface)
	if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
		return false, errors.New("ImplementsInterface needs to have a pointer to a interface variable passed to it")
	}

	return reflect.TypeOf(actual).Implements(ifaceType.Elem()), nil
}

func (ii implementInterfaceMatcher) FailureMessage(actual interface{}) string {
	ifaceType := reflect.TypeOf(ii.iface).Elem()
	actualType := reflect.TypeOf(actual)

	name := actualType.Name()
	if actualType.Kind() == reflect.Ptr {
		name = "*" + actualType.Elem().Name()
	}

	return fmt.Sprintf("Type %v does not implement interface %v", name, ifaceType.Name())
}

func (ii implementInterfaceMatcher) NegatedFailureMessage(actual interface{}) string {
	ifaceType := reflect.TypeOf(ii.iface).Elem()
	actualType := reflect.TypeOf(actual)

	name := actualType.Name()
	if actualType.Kind() == reflect.Ptr {
		name = "*" + actualType.Elem().Name()
	}

	return fmt.Sprintf("Type %v implements interface %v", name, ifaceType.Name())
}
