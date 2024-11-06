package mock

import (
	"errors"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

var (
	errTagOperationNotSupported = errors.New("tag operation not supported")
)

type tagOperation int

const (
	tagOperationCreate tagOperation = iota
	tagOperationDestroy
)

func (a *mockAPI) tagOperation(o types.Object, op tagOperation) (bool, error) {
	rwt, ok := o.(*corev1.ResourceWithTag)
	if !ok {
		return false, nil
	}

	obj, ok := a.data[rwt.ResourceIdentifier]
	if !ok {
		return true, api.ErrNotFound
	}

	if err := a.applyTagOperation(op, obj, rwt); err != nil {
		return true, err
	}

	return true, nil
}

func (a *mockAPI) applyTagOperation(op tagOperation, obj *APIObject, rwt *corev1.ResourceWithTag) error {
	if op == tagOperationCreate {
		if _, ok := obj.tags[rwt.Tag]; ok {
			return api.NewHTTPError(422, "POST", nil, nil)
		}

		obj.tags[rwt.Tag] = true
	} else if op == tagOperationDestroy {
		if _, ok := obj.tags[rwt.Tag]; !ok {
			return api.NewHTTPError(404, "GET", nil, api.ErrNotFound)
		}

		delete(obj.tags, rwt.Tag)
	} else {
		return errTagOperationNotSupported
	}

	return nil
}
