package mock

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("data views", func() {
	var a API
	var initialObjects []*testObject

	BeforeEach(func() {
		a = NewMockAPI()
		initialObjects = make([]*testObject, 0, 15)

		for i := 0; i < cap(initialObjects); i++ {
			obj := testObject{}
			a.FakeExisting(&obj)
			initialObjects = append(initialObjects, &obj)
		}
	})

	It("can list objects destroyed after timestamp", func() {
		for _, obj := range initialObjects[5:10] {
			err := a.Destroy(context.TODO(), obj)
			Expect(err).ToNot(HaveOccurred())
		}

		ts := time.Now()

		for _, obj := range initialObjects[10:15] {
			err := a.Destroy(context.TODO(), obj)
			Expect(err).ToNot(HaveOccurred())
		}

		Expect(a.DestroyedAfter(ts).Unwrap()).To(ConsistOf(initialObjects[10:15]))
	})

	It("can list objects updated after timestamp", func() {
		for _, obj := range initialObjects[5:10] {
			err := a.Update(context.TODO(), obj)
			Expect(err).ToNot(HaveOccurred())
		}

		ts := time.Now()

		for _, obj := range initialObjects[10:15] {
			err := a.Update(context.TODO(), obj)
			Expect(err).ToNot(HaveOccurred())
		}

		Expect(a.UpdatedAfter(ts, false).Unwrap()).To(ConsistOf(initialObjects[10:15]))
	})

	It("can list objects created after timestamp", func() {
		ts := time.Now()
		createdAfterTS := make([]*testObject, 0, 5)

		for i := 0; i < cap(createdAfterTS); i++ {
			obj := testObject{}
			err := a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())
			createdAfterTS = append(createdAfterTS, &obj)
		}

		Expect(a.CreatedAfter(ts, false).Unwrap()).To(ConsistOf(createdAfterTS))
	})
})
