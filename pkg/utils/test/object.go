package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type hookErrorCheck func(context.Context) error

// ObjectTests contains the logic to test any Object implementation with the hooks it implements, checking
// * if the Object actually implements the interface of the hook
// * if the Object has an identifier field
// * calling the hook function with incomplete contexts gives none or the correct error (meaning the error is handled correctly)
//
func ObjectTests(o types.Object, hooks ...interface{}) {
	ginkgo.It("has an identifier", func() {
		_, err := api.GetObjectIdentifier(o, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	for _, hookPtr := range hooks {
		hook := reflect.TypeOf(hookPtr).Elem()
		name := hook.Name()

		ginkgo.Context(fmt.Sprintf("implementing %v", name), func() {
			ginkgo.It("actually implements the interface", func() {
				implementsHook := reflect.TypeOf(o).Implements(hook)
				gomega.Expect(implementsHook).To(gomega.BeTrue())
			})

			testHookHandlingIncompleteContext(o, name)
		})
	}
}

// testHookHandlingIncompleteContext calls the function of a hook with incomplete contexts as argument,
// checking if it either returns the expected ErrContextKeyNotSet or no error at all.
//
// If this randomly fails for a single or few Objects, maybe the Hook implementation of the Object(s) needs
// actual data instead of an empty json Object - easiest workaround is to check all the context keys in the
// hook implementation, I have no idea how I could generate sensible test data to pass to the hooks and when
// I wrote this, all Objects were just fine with an empty json object. -- Mara @LittleFox94 Grosch, 2022-02-04
func testHookHandlingIncompleteContext(o types.Object, hook string) {
	// This is our map from hook name to a function calling its function and only returning the error returned by
	// it. When new hooks are supported by the generic client, they'd have to be added here to be checked.
	supportedHooks := map[string]hookErrorCheck{
		"Object": func(ctx context.Context) error {
			_, err := o.EndpointURL(ctx)
			return err
		},
		"PaginationSupportHook": func(ctx context.Context) error {
			_, err := o.(types.PaginationSupportHook).HasPagination(ctx)
			return err
		},
		"RequestBodyHook": func(ctx context.Context) error {
			_, err := o.(types.RequestBodyHook).FilterAPIRequestBody(ctx)
			return err
		},
		"RequestFilterHook": func(ctx context.Context) error {
			_, err := o.(types.RequestFilterHook).FilterAPIRequest(ctx, httptest.NewRequest("GET", "/", nil))
			return err
		},
		"ResponseDecodeHook": func(ctx context.Context) error {
			err := o.(types.ResponseDecodeHook).DecodeAPIResponse(ctx, bytes.NewBuffer([]byte(`{}`)))
			return err
		},
		"ResponseFilterHook": func(ctx context.Context) error {
			rec := httptest.NewRecorder()
			rec.WriteHeader(200)
			_, _ = rec.WriteString(`{}`)
			_, err := o.(types.ResponseFilterHook).FilterAPIResponse(ctx, rec.Result())
			return err
		},
	}

	if errorCheck, ok := supportedHooks[hook]; ok {
		args := []interface{}{
			func(ctx context.Context) {
				err := errorCheck(ctx)

				// It can be fine with incomplete context, but if it fails, than with the error indicating
				// it checked for it - the one OperationFromContext and co. return
				if err != nil {
					gomega.Expect(err).To(
						gomega.MatchError(types.ErrContextKeyNotSet),
					)
				}
			},

			// technically we have to omit "url" for checking EndpointURL, as it is not set there, yet
			// EndpointURL using URL from context should crash in other tests though, so we take the simplicity of only adding a test case
			// for others, instead of generating different sets of test cases.
			ginkgo.Entry("missing options", makeTestContext("operation", "url")),
			ginkgo.Entry("missing operation", makeTestContext("options", "url")),
		}

		if hook != "Object" {
			args = append(args,
				ginkgo.Entry("missing url", makeTestContext("options", "operation")),
			)
		}

		ginkgo.DescribeTable("handles being called with context", args...)
	}
}

func makeTestContext(elems ...string) context.Context {
	ctx := context.TODO()

	for _, elem := range elems {
		switch elem {
		case "operation":
			ctx = types.ContextWithOperation(ctx, types.OperationList)
		case "options":
			ctx = types.ContextWithOptions(ctx, &types.ListOptions{})
		case "url":
			u, _ := url.Parse("http://localhost:1312")
			ctx = types.ContextWithURL(ctx, *u)
		}
	}

	return ctx
}
