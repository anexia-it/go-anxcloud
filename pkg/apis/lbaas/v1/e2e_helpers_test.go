//go:build integration
// +build integration

package v1_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var apiClient api.API

type LBaaSE2ETestRun struct {
	Name string
	Port int
}

func waitObjectReady(ctx *context.Context, o *types.Object) {
	It("eventually is ready", func() {
		pollCheck := func() error {
			err := apiClient.Get(*ctx, *o)
			if err != nil {
				return err
			}

			hasState, ok := (*o).(gs.StateRetriever)
			if !ok {
				return errors.New("object does not implement StateRetriever")
			}

			if hasState.StateError() {
				return errors.New("object is in error state")
			}
			if !hasState.StateOK() {
				return errors.New("object is not ready")
			}
			return nil
		}

		Eventually(pollCheck, 5*time.Minute, 3*time.Second).Should(Succeed())
	})
}

func waitObjectGone(ctx *context.Context, o *types.Object) {
	It("eventually is gone", func() {
		pollCheck := func() error {
			err := apiClient.Get(*ctx, *o)
			if err == nil {
				return errors.New("object still exists")
			}

			var he api.HTTPError
			if !errors.As(err, &he) {
				return errors.New("expected HTTP error")
			}
			if he.StatusCode() != http.StatusNotFound {
				return errors.New("expected 404 Not Found")
			}
			return nil
		}

		Eventually(pollCheck, 5*time.Minute, 3*time.Second).Should(Succeed())
	})
}

func createObject(retriever func() types.Object, waitReady bool) func() {
	var ctx context.Context
	var identifier string

	var objectType reflect.Type

	var object types.Object
	var emptyObject types.Object
	var identifiedObject types.Object

	BeforeAll(func() {
		ctx = context.TODO()
		object = retriever()

		objectType = reflect.TypeOf(object).Elem()
		emptyObject = reflect.New(objectType).Interface().(types.Object)

		DeferCleanup(func() {
			if identifiedObject != nil {
				err := apiClient.Destroy(ctx, identifiedObject)
				if err != nil {
					GinkgoWriter.Printf("Error deleting Object %v: %v\n", identifiedObject, err)
				}
			}
		})
	})

	It("is created successfully", func() {
		err := apiClient.Create(ctx, object)
		Expect(err).NotTo(HaveOccurred())

		identifier, err = types.GetObjectIdentifier(object, true)
		Expect(err).NotTo(HaveOccurred())

		identifiedObjectValue := reflect.New(objectType)
		identifiedObjectValue.Elem().FieldByName("Identifier").SetString(identifier)
		identifiedObject = identifiedObjectValue.Interface().(types.Object)
	})

	if waitReady {
		waitObjectReady(&ctx, &identifiedObject)
	}

	It("is included when List-ing", func() {
		var oc types.ObjectChannel
		err := apiClient.List(ctx, emptyObject, api.ObjectChannel(&oc))
		Expect(err).NotTo(HaveOccurred())

		identifiers := make([]string, 0, 50)
		for retriever := range oc {
			err := retriever(emptyObject)
			Expect(err).NotTo(HaveOccurred())

			id, err := types.GetObjectIdentifier(emptyObject, true)
			Expect(err).NotTo(HaveOccurred())

			identifiers = append(identifiers, id)
		}

		Expect(identifiers).To(ContainElements(identifier))
	})

	return func() {
		It("is destroyed successfully", func() {
			err := apiClient.Destroy(ctx, identifiedObject)
			Expect(err).NotTo(HaveOccurred())
		})

		waitObjectGone(&ctx, &identifiedObject)

		It("marks Object as successfully destroyed", func() {
			identifiedObject = nil
		})
	}
}

func updateObject(retriever func() types.Object, waitReady bool, validate ...func(types.Object)) {
	Context("updating the Object", Ordered, func() {
		var ctx context.Context
		var obj types.Object

		BeforeAll(func() {
			ctx = context.TODO()
			obj = retriever()
		})

		It("is updated successfully", func() {
			err := apiClient.Update(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		if waitReady {
			waitObjectReady(&ctx, &obj)
		}

		if len(validate) > 0 {
			It("has the correct parameters", func() {
				for _, val := range validate {
					val(obj)
				}
			})
		}
	})
}
