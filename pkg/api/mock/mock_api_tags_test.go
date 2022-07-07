package mock

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
)

var _ = Describe("Mock API Tags", func() {
	var (
		a       *mockAPI
		objects []*testObject
	)
	BeforeEach(func() {
		a = NewMockAPI().(*mockAPI)

		objects = make([]*testObject, 0, 5)

		for i := 0; i < 5; i++ {
			obj := testObject{}
			objects = append(objects, &obj)
			err := a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("returns api.ErrNotFound when object to tag does not exist in mock API", func() {
		isTagOperation, err := a.tagOperation(&corev1.ResourceWithTag{Identifier: "not-in-api"}, tagOperationCreate)
		Expect(isTagOperation).To(BeTrue())
		Expect(err).To(MatchError(api.ErrNotFound))
	})

	It("is no-op when provided types.Object is not (*corev1.ResourceWithTag)", func() {
		isTagOperation, err := a.tagOperation(&testObject{}, tagOperationCreate)
		Expect(isTagOperation).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns api.HTTPError (422) when an object is tagged more than once with the same tag", func() {
		rwt := &corev1.ResourceWithTag{Identifier: objects[0].Identifier, Tag: "tagname"}

		isTagOperation, err := a.tagOperation(rwt, tagOperationCreate)
		Expect(isTagOperation).To(BeTrue())
		Expect(err).ToNot(HaveOccurred())

		isTagOperation, err = a.tagOperation(rwt, tagOperationCreate)
		Expect(isTagOperation).To(BeTrue())
		var he api.HTTPError
		Expect(errors.As(err, &he)).To(BeTrue())
		Expect(he.StatusCode()).To(Equal(422))
	})

	It("returns api.HTTPError (404) when an object is untagged from a tag it doesn't have", func() {
		rwt := &corev1.ResourceWithTag{Identifier: objects[0].Identifier, Tag: "tagname"}

		isTagOperation, err := a.tagOperation(rwt, tagOperationDestroy)
		Expect(isTagOperation).To(BeTrue())
		Expect(err).To(MatchError(api.ErrNotFound))
	})

	It("returns errTagOperationNotSupported if an unsupported tag operation was provided", func() {
		rwt := &corev1.ResourceWithTag{Identifier: objects[0].Identifier, Tag: "tagname"}

		isTagOperation, err := a.tagOperation(rwt, 3)
		Expect(isTagOperation).To(BeTrue())
		Expect(err).To(MatchError(errTagOperationNotSupported))
	})
})
