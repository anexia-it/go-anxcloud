package pagination

import (
	"context"
	"errors"
	"fmt"
	"go.anx.io/go-anxcloud/pkg/utils/param"
	"reflect"
)

var (
	ErrConditionNeverMet = errors.New("looped all items and the condition was never met")
)

type Page interface {
	Num() int
	Size() int
	Total() int
	Options() []param.Parameter
	Content() interface{}
}

// Pageable should be implemented by evey struct that supports pagination.
type Pageable interface {
	GetPage(ctx context.Context, page, limit int, opts ...param.Parameter) (Page, error)
	NextPage(ctx context.Context, page Page) (Page, error)
}

// HasNext is a helper function which checks whether there are more pages to fetch
func HasNext(page Page) bool {
	return page.Num() < page.Total()
}

type UntilTrueFunc func(interface{}) (bool, error)

// LoopUntil takes a pageable and loops over it until untilFunc returns true or an error.
func LoopUntil(ctx context.Context, pageable Pageable, untilFunc UntilTrueFunc, opts ...param.Parameter) error {
	page, err := pageable.GetPage(ctx, 1, 10, opts...)
	if err != nil {
		return err
	}

	content := reflect.ValueOf(page.Content())
	switch content.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		stopped := false
		for i := 0; i < content.Len(); i++ {
			content.Interface()
			stop, err := untilFunc(content.Index(i).Interface())
			if err != nil {
				return err
			}
			if stop {
				stopped = true
				break
			}
		}
		if !stopped {
			return ErrConditionNeverMet
		}
	default:
		panic(fmt.Sprintf("The page content is supposed to be of type slice or array but was %T", page.Content()))
	}
	return nil
}

type CancelFunc func()

// AsChan takes a Pageable and returns its Pageable.Content via a channel until there are no more pages or
// CancelFunc gets called by the consumer.
func AsChan(ctx context.Context, pageable Pageable, opts ...param.Parameter) (chan interface{}, CancelFunc) {
	consumer := make(chan interface{})
	done := make(chan interface{})
	cancel := func() {
		close(done)
	}

	go func() {
		defer close(consumer)
		_ = LoopUntil(ctx, pageable, func(i interface{}) (bool, error) {
			select {
			case consumer <- i:
			case _, ok := <-done:
				if !ok {
					return true, nil
				}
			}
			return false, nil
		}, opts...)
	}()
	return consumer, cancel
}
