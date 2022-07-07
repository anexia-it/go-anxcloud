package mock

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mock API Object", func() {
	var (
		raw testObject
		obj APIObject
	)

	BeforeEach(func() {
		raw = testObject{}
		obj = APIObject{
			wrapped: &raw,
			tags:    map[string]interface{}{"test_1": true, "test_3": true},

			existing: true,

			createdCount:   4,
			updatedCount:   3,
			destroyedCount: 2,

			createdTime:   time.Now(),
			updatedTime:   time.Now(),
			destroyedTime: time.Now(),
		}
	})

	It("can unwrap raw object", func() {
		Expect(obj.Unwrap()).To(Equal(&raw))
	})

	It("can get tags of an object", func() {
		Expect(obj.Tags()).To(ConsistOf([]string{"test_1", "test_3"}))
	})

	It("can check if an object was tagged with tags", func() {
		Expect(obj.HasTags("test_1", "test_3")).To(BeTrue())
		Expect(obj.HasTags("test_1")).To(BeTrue())
		Expect(obj.HasTags("test_2")).To(BeFalse())
	})

	It("returns whether or not an object exists", func() {
		Expect(obj.Existing()).To(BeTrue())
		obj.existing = false
		Expect(obj.Existing()).To(BeFalse())
	})

	It("returns how often an object was created", func() {
		Expect(obj.CreatedCount()).To(Equal(4))
	})

	It("returns how often an object was updated", func() {
		Expect(obj.UpdatedCount()).To(Equal(3))
	})

	It("returns how often an object was destroyed", func() {
		Expect(obj.DestroyedCount()).To(Equal(2))
	})
})
