package matcher

import "go.anx.io/go-anxcloud/pkg/api/mock"

func castToAPIObjectPointer(actual interface{}) (*mock.APIObject, error) {
	obj, ok := actual.(*mock.APIObject)
	if !ok {
		return nil, ErrMatcherExpectsAPIObjectPointer
	}
	return obj, nil
}
