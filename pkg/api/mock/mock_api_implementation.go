package mock

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/mitchellh/copystructure"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/mock/internal"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

type mockAPI struct {
	data mockDataView
	mu   sync.Mutex
}

// NewMockAPI creates a new MockAPI instance
func NewMockAPI() API {
	return &mockAPI{
		data: make(map[string]*APIObject),
	}
}

// Get retrieves an Object from MockAPIs local storage by its identifier
func (a *mockAPI) Get(ctx context.Context, o types.IdentifiedObject, opts ...types.GetOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	apiObject, err := a.findObject(o)
	if err != nil {
		return fmt.Errorf("couldn't find object in mock api: %w", err)
	}

	if apiObject.existing && reflect.TypeOf(apiObject.wrapped) == reflect.TypeOf(o) {
		copy, err := copystructure.Copy(apiObject.wrapped)
		if err != nil {
			return err
		}
		reflect.ValueOf(o).Elem().Set(reflect.ValueOf(copy).Elem())
		return nil
	}

	return api.ErrNotFound
}

// Lists Objects filtered by types.FilterObject
func (a *mockAPI) List(ctx context.Context, o types.FilterObject, opts ...types.ListOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	filteredData, err := listDataAggregation(o, a.data)
	if err != nil {
		return err
	}

	return listOutput(ctx, filteredData, opts)
}

func listDataAggregation(o types.Object, data mockDataView) ([]types.Object, error) {
	objectsWithSameType := make([]types.Object, 0)

	if coreResource, isCoreResource := o.(*corev1.Resource); isCoreResource {
		for id, obj := range data {
			if !obj.existing || !obj.HasTags(coreResource.Tags...) {
				continue
			}

			objectsWithSameType = append(objectsWithSameType, &corev1.Resource{
				Identifier: id,
				Type:       *internal.ObjectCoreType(obj.wrapped),
			})
		}
	} else {
		for _, obj := range data {
			if obj.existing && reflect.TypeOf(obj.wrapped) == reflect.TypeOf(o) {
				copy, err := copystructure.Copy(obj.wrapped)
				if err != nil {
					return nil, err
				}
				objectsWithSameType = append(objectsWithSameType, copy.(types.Object))
			}
		}
	}
	return objectsWithSameType, nil
}

func listOutput(ctx context.Context, objects []types.Object, opts []types.ListOption) error {
	options := types.ListOptions{}
	for _, opt := range opts {
		opt.ApplyToList(&options)
	}

	var channelPageIterator types.PageInfo
	if options.ObjectChannel != nil && !options.Paged {
		options.Paged = true
		options.Page = 1
		options.EntriesPerPage = api.ListChannelDefaultPageSize
		options.PageInfo = &channelPageIterator
	} else if options.ObjectChannel != nil && options.PageInfo != nil {
		return api.ErrCannotListChannelAndPaged
	}

	if options.Paged {
		if options.Page == 0 {
			options.Page = 1
		}

		pageInfo, err := newMockPageIter(objects, options.EntriesPerPage, options.Page)
		if err != nil {
			return err
		}
		*options.PageInfo = pageInfo
	}

	if options.ObjectChannel != nil {
		c := make(chan types.ObjectRetriever)
		*options.ObjectChannel = c

		go objectChannel(ctx, channelPageIterator, c)
	}

	return nil
}

func objectChannel(ctx context.Context, pi types.PageInfo, c chan types.ObjectRetriever) {
	objectRetrieved := make(chan bool)
	var pageData []types.Object

outer:
	for pi.Next(&pageData) {
		if len(pageData) == 0 {
			break outer
		}
		for _, o := range pageData {
			// since we are in a goroutine, we might already be in the next iteration of this loop
			// at the time the receiving end of this channel calls the closure. Having a loop-body
			// scoped variables makes the data for the closure perfectly identified.
			closureData := o
			c <- func(out types.Object) error {
				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(closureData).Elem())

				select {
				case <-ctx.Done():
				case objectRetrieved <- true:
				}

				return nil
			}

			select {
			case <-ctx.Done():
				break outer
			case <-objectRetrieved:
			}
		}
	}

	close(c)
}

// Create stores an types.Object to the MockAPIs local storage.
// When the provided Object has no Identifier set, a random one is set.
// An already set Identifier is kept as-is, without any validation.
func (a *mockAPI) Create(ctx context.Context, o types.Object, opts ...types.CreateOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if isTag, err := a.tagOperation(o, tagOperationCreate); isTag {
		return err
	}

	identifier := makeObjectIdentifiable(o)

	apiObject, exists := a.data[identifier]
	if !exists {
		apiObject = &APIObject{tags: make(map[string]interface{})}
		a.data[identifier] = apiObject
	}

	cloned, err := copystructure.Copy(o)
	if err != nil {
		return err
	}

	apiObject.wrapped = cloned.(types.Object)
	apiObject.existing = true
	apiObject.createdCount++
	apiObject.createdTime = time.Now()

	return nil
}

// Update overwrites a types.Object in MockAPIs local storage
func (a *mockAPI) Update(ctx context.Context, o types.IdentifiedObject, opts ...types.UpdateOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	apiObject, err := a.findObject(o)
	if err != nil {
		return fmt.Errorf("couldn't find object in mock api: %w", err)
	}

	merged, err := merge(apiObject.wrapped, o)
	if err != nil {
		if errors.Is(err, ErrMergeTypeMissmatch) {
			return api.ErrNotFound
		}
		return err
	}

	apiObject.wrapped = merged.(types.Object)
	apiObject.updatedCount++
	apiObject.updatedTime = time.Now()

	return updateSetMergedToInputObject(merged.(types.Object), o)
}

func updateSetMergedToInputObject(merged, o types.Object) error {
	cloned, err := copystructure.Copy(merged)
	if err != nil {
		return err
	}

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(cloned).Elem())
	return nil
}

// Destroy removes a types.Object from MockAPIs local storage
func (a *mockAPI) Destroy(ctx context.Context, o types.IdentifiedObject, opts ...types.DestroyOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if isTag, err := a.tagOperation(o, tagOperationDestroy); isTag {
		return err
	}

	apiObject, err := a.findObject(o)
	if err != nil {
		return fmt.Errorf("couldn't find object in mock api: %w", err)
	}

	apiObject.existing = false
	apiObject.destroyedCount++
	apiObject.destroyedTime = time.Now()

	return nil
}

// FakeExisting stores an object to the mock API without incrementing the created count
// and returns the created identifier
func (a *mockAPI) FakeExisting(o types.Object, tags ...string) string {
	a.mu.Lock()
	defer a.mu.Unlock()

	identifier := makeObjectIdentifiable(o)

	mao := APIObject{
		wrapped:  o,
		tags:     make(map[string]interface{}),
		existing: true,
	}

	for _, tag := range tags {
		mao.tags[tag] = true
	}

	a.data[identifier] = &mao

	return identifier
}

// Inspect retrieves an *APIObject by it's Identifier or nil if not found
func (a *mockAPI) Inspect(identifier string) *APIObject {
	if obj, ok := a.data[identifier]; ok {
		return obj
	}
	return nil
}

func (a *mockAPI) findObject(o types.Object) (*APIObject, error) {
	identifier, err := types.GetObjectIdentifier(o, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get object identifier: %w", err)
	}

	val, ok := a.data[identifier]
	if !ok || !val.existing {
		return nil, api.ErrNotFound
	}

	return val, nil
}
