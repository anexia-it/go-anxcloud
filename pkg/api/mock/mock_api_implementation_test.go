package mock

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
)

var _ = Describe("Mock API implementation", func() {
	var (
		a *mockAPI
	)

	BeforeEach(func() {
		a = NewMockAPI().(*mockAPI)
	})

	Context("FakeExisting", func() {
		It("can add objects to mock API", func() {
			obj := testObject{}
			a.FakeExisting(&obj, "tagname-1", "tagname-2")
			Expect(a.Existing().Unwrap()).To(ContainElement(&obj))
		})
	})

	Context("Inspect", func() {
		It("can retrieve an existing *APIObject", func() {
			obj := testObject{}
			id := a.FakeExisting(&obj)
			Expect(a.Inspect(id).Unwrap()).To(Equal(&obj))
		})

		It("returns nil if id not in API", func() {
			Expect(a.Inspect("not-in-api")).To(BeNil())
		})
	})

	Context("Get operations", func() {
		It("returns types.ErrUnidentifiedObject if object does not have an identifier set", func() {
			err := a.Get(context.TODO(), &testObject{})
			Expect(err).To(MatchError(types.ErrUnidentifiedObject))
		})

		It("returns api.ErrNotFound if the object does not exist in the mock API", func() {
			err := a.Get(context.TODO(), &testObject{Identifier: "not-in-api"})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("retrieve identified objects from mock API", func() {
			obj := testObject{}
			a.FakeExisting(&obj)
			err := a.Get(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("List operations", func() {
		var (
			objects []*testObject
		)

		BeforeEach(func() {
			objects = make([]*testObject, 0, 5)
			for i := 0; i < 5; i++ {
				obj := testObject{}
				objects = append(objects, &obj)
				a.FakeExisting(&obj)
			}
		})

		It("can list objects with types.ObjectChannel", func() {
			var oc types.ObjectChannel
			err := a.List(context.TODO(), &testObject{}, api.ObjectChannel(&oc))
			Expect(err).ToNot(HaveOccurred())

			var res []types.Object
			for retriever := range oc {
				var obj testObject
				err := retriever(&obj)
				Expect(err).ToNot(HaveOccurred())
				res = append(res, &obj)
			}

			Expect(res).To(ConsistOf(objects))
		})

		It("can list objects with types.PageInfo", func() {
			var pi types.PageInfo
			err := a.List(context.TODO(), &testObject{}, api.Paged(1, 1, &pi))
			Expect(err).ToNot(HaveOccurred())

			res := make([]types.Object, 0, pi.TotalItems())
			var page []types.Object
			for pi.Next(&page) {
				res = append(res, page...)
			}

			Expect(res).To(ConsistOf(objects))
		})

		It("returns api.ErrCannotListChannelAndPaged when both types.PageInfo and types.ObjectChannel requested", func() {
			var oc types.ObjectChannel
			var pi types.PageInfo
			err := a.List(context.TODO(), &testObject{}, api.ObjectChannel(&oc), api.Paged(1, 1, &pi))
			Expect(err).To(MatchError(api.ErrCannotListChannelAndPaged))
		})

		It("returns ErrPageSizeCannotBeZero when limit set to <0>", func() {
			var pi types.PageInfo
			err := a.List(context.TODO(), &testObject{}, api.Paged(1, 0, &pi))
			Expect(err).To(MatchError(ErrPageSizeCannotBeZero))
		})

		It("defaults to page <1> when page <0> was requested", func() {
			var pi types.PageInfo
			err := a.List(context.TODO(), &testObject{}, api.Paged(0, 3, &pi))
			Expect(err).ToNot(HaveOccurred())
			Expect(pi.CurrentPage()).To(BeEquivalentTo(1))

			res := make([]types.Object, 0, pi.TotalItems())
			var page []types.Object
			for pi.Next(&page) {
				res = append(res, page...)
			}

			Expect(res).To(ConsistOf(objects))
		})

		It("supports empty responses", func() {
			var oc types.ObjectChannel
			err := a.List(context.TODO(), &testObjectWithoutIdentifier{}, api.ObjectChannel(&oc))
			Expect(err).ToNot(HaveOccurred())
			Eventually(oc).Should(BeClosed())
		})

		It("supports cancelled context", func() {
			var oc types.ObjectChannel
			ctx, cancel := context.WithCancel(context.TODO())
			err := a.List(ctx, &testObject{}, api.ObjectChannel(&oc))
			Expect(err).ToNot(HaveOccurred())
			cancel()
			Eventually(oc).Should(BeClosed())
		})

		It("has limited support for corev1.Resource's", func() {
			a.FakeExisting(&lbaasv1.Bind{}, "tag-1")
			a.FakeExisting(&lbaasv1.Bind{}, "tag-1", "tag-2")
			a.FakeExisting(&lbaasv1.Bind{}, "tag-1", "tag-2")
			a.FakeExisting(&lbaasv1.Bind{}, "tag-1", "tag-2", "tag-3")

			var pi types.PageInfo
			err := a.List(context.TODO(), &corev1.Resource{Tags: []string{"tag-1", "tag-2"}}, api.Paged(1, 10, &pi))
			Expect(err).ToNot(HaveOccurred())

			var page []types.Object
			pi.Next(&page)
			Expect(page).To(HaveLen(3))
			r := page[0].(*corev1.Resource)
			Expect(r.Type.Identifier).To(Equal("bd24def982aa478fb3352cb5f49aab47"))
		})
	})

	Context("Create operations", func() {
		It("supports creating API resources", func() {
			obj := testObject{}
			err := a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())
			Expect(a.All().Unwrap()).To(ContainElement(&obj))
		})

		It("supports corev1.ResourceWithTag", func() {
			id := a.FakeExisting(&testObject{})
			err := a.Create(context.TODO(), &corev1.ResourceWithTag{Identifier: id, Tag: "tagname"})
			Expect(err).ToNot(HaveOccurred())
			Expect(a.Inspect(id).HasTags("tagname")).To(BeTrue())
		})

		It("handles errors on tag operation", func() {
			err := a.Create(context.TODO(), &corev1.ResourceWithTag{Identifier: "not-in-api", Tag: "tagname"})
			Expect(err).To(MatchError(api.ErrNotFound))
		})
	})

	Context("Update operations", func() {
		It("returns api.ErrNotFound if object to be updated not in mock API", func() {
			err := a.Update(context.TODO(), &testObject{Identifier: "not-in-api"})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("returns api.ErrObjectWithoutIdentifier if object does not have an identifier", func() {
			err := a.Update(context.TODO(), &testObjectWithoutIdentifier{})
			Expect(err).To(MatchError(types.ErrObjectWithoutIdentifier))
		})

		It("returns error when merge fails", func() {
			id := a.FakeExisting(&testObject{})
			err := a.Update(context.TODO(), &testObject2{Identifier: id})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("can update mock API objects", func() {
			id := a.FakeExisting(&testObject{TestFieldB: "some text B"})
			err := a.Update(context.TODO(), &testObject{Identifier: id, TestFieldA: "some text A"})
			Expect(err).ToNot(HaveOccurred())
			Expect(a.Existing().Unwrap()[id]).To(Equal(&testObject{Identifier: id, TestFieldA: "some text A", TestFieldB: "some text B"}))
		})
	})

	Context("Destroy operations", func() {
		It("can destroy objects", func() {
			id := a.FakeExisting(&testObject{})
			err := a.Destroy(context.TODO(), &testObject{Identifier: id})
			Expect(err).ToNot(HaveOccurred())
		})

		It("supports corev1.ResourceWithTag", func() {
			id := a.FakeExisting(&testObject{}, "tagname")
			err := a.Destroy(context.TODO(), &corev1.ResourceWithTag{Identifier: id, Tag: "tagname"})
			Expect(err).ToNot(HaveOccurred())
			Expect(a.Inspect(id).HasTags("tagname")).To(BeFalse())
		})

		It("returns 404 if object to be untagged does not have tag provided in corev1.ResourceWithTag", func() {
			id := a.FakeExisting(&testObject{})
			err := a.Destroy(context.TODO(), &corev1.ResourceWithTag{Identifier: id, Tag: "tagname"})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("returns api.ErrNotFound if object to be destroyed not in mock API", func() {
			err := a.Destroy(context.TODO(), &testObject{Identifier: "not-in-api"})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("returns api.ErrObjectWithoutIdentifier if object does not have an identifier", func() {
			err := a.Destroy(context.TODO(), &testObjectWithoutIdentifier{})
			Expect(err).To(MatchError(types.ErrObjectWithoutIdentifier))
		})
	})
})
