package mock

import (
	"time"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

type mockDataView map[string]*APIObject

func (a *mockAPI) Existing() mockDataView {
	return a.filteredDataMap(func(o *APIObject) bool {
		return o.existing
	})
}

func (a *mockAPI) All() mockDataView {
	return a.filteredDataMap(func(o *APIObject) bool {
		return true
	})
}

func (a *mockAPI) CreatedAfter(t time.Time, onlyExisting bool) mockDataView {
	return a.filteredDataMap(func(o *APIObject) bool {
		return o.createdTime.After(t) && (!onlyExisting || o.existing)
	})
}

func (a *mockAPI) UpdatedAfter(t time.Time, onlyExisting bool) mockDataView {
	return a.filteredDataMap(func(o *APIObject) bool {
		return o.updatedTime.After(t) && (!onlyExisting || o.existing)
	})
}

func (a *mockAPI) DestroyedAfter(t time.Time) mockDataView {
	return a.filteredDataMap(func(o *APIObject) bool {
		return o.destroyedTime.After(t)
	})
}

func (a *mockAPI) filteredDataMap(filter func(o *APIObject) bool) mockDataView {
	out := make(mockDataView)
	for key, obj := range a.data {
		if filter(obj) {
			out[key] = obj
		}
	}
	return out
}

type rawObjectMap map[string]types.Object

func (m mockDataView) Unwrap() rawObjectMap {
	out := make(rawObjectMap)
	for key, obj := range m {
		out[key] = obj.Unwrap()
	}
	return out
}
