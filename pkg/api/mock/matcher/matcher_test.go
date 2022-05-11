package matcher

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.anx.io/go-anxcloud/pkg/api/mock"
	"go.anx.io/go-anxcloud/pkg/utils/object/compare"
)

func TestMockAPIMatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock API Matcher")
}

var _ = Describe("custom gomega matcher", func() {
	var (
		a       mock.API
		objects []*testObject
	)

	BeforeEach(func() {
		a = mock.NewMockAPI()

		objects = make([]*testObject, 0, 8)

		for i := 0; i < 4; i++ {
			obj0 := testObject{TestFieldA: "text 0"}
			obj1 := testObject{TestFieldA: "text 1"}
			a.FakeExisting(&obj0)
			a.FakeExisting(&obj1)
			objects = append(objects, &obj0, &obj1)
		}
	})

	Context("ContainElementNTimes", func() {
		It("can test if the API list contains elements n times", func() {
			Expect(a.Existing()).To(ContainElementNTimes(Object(&testObject{}), 8))
			Expect(a.Existing()).To(ContainElementNTimes(Object(&testObject2{}), 0))
		})

		It("returns error when ACTUAL not of type array/slice/map", func() {
			matcher := ContainElementNTimes(nil, 5)
			_, err := matcher.Match(true)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("ContainElementNTimes matcher expects an array/slice/map."))
		})

		It("fallbacks to equal matcher when element is no types.GomegaMatcher", func() {
			matcher := ContainElementNTimes("test", 2)
			Expect(matcher.Match([]string{"test", "toast", "test"})).To(BeTrue())
			Expect(matcher.Match([]string{"test", "toast"})).To(BeFalse())
		})

		It("returns error when sub-match returns error", func() {
			matcher := ContainElementNTimes(ContainElement(1), 2)
			_, err := matcher.Match([]string{"test1", "test2"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("ContainElement matcher expects an array/slice/map."))
		})

		It("has (Negated)FailureMessage", func() {
			Expect(ContainElementNTimes(nil, 0).FailureMessage(nil)).To(ContainSubstring("to contain n elements matching"))
			Expect(ContainElementNTimes(nil, 0).NegatedFailureMessage(nil)).To(ContainSubstring("not to contain n elements matching"))
		})
	})

	Context("CRUD Matcher", func() {
		It("can match objects created n times", func() {
			// clouddnsv1.Zones have their name as Identifier and can be created and destroyed multiple times
			Expect(ContainElement(Created(5)).Match(a.All())).To(BeFalse())
			Expect(ContainElement(Created(0)).Match(a.All())).To(BeTrue())

			obj := testObject{}

			err := a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			err = a.Destroy(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			err = a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			Expect(ContainElementNTimes(Created(2), 1).Match(a.All())).To(BeTrue())
		})

		It("can match objects updated n times", func() {
			Expect(ContainElement(Updated(5)).Match(a.All())).To(BeFalse())
			Expect(ContainElement(Updated(0)).Match(a.All())).To(BeTrue())

			obj := testObject{}
			a.FakeExisting(&obj)
			err := a.Update(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())
			Expect(ContainElementNTimes(Updated(1), 1).Match(a.All())).To(BeTrue())
		})

		It("can match objects destroyed n times", func() {
			// clouddnsv1.Zones have their name as Identifier and can be created and destroyed multiple times
			Expect(ContainElement(Destroyed(5)).Match(a.All())).To(BeFalse())
			Expect(ContainElement(Destroyed(0)).Match(a.All())).To(BeTrue())

			obj := testObject{}
			a.FakeExisting(&obj)

			err := a.Destroy(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			err = a.Create(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			err = a.Destroy(context.TODO(), &obj)
			Expect(err).ToNot(HaveOccurred())

			Expect(ContainElementNTimes(Destroyed(2), 1).Match(a.All())).To(BeTrue())
		})

		DescribeTable("returns error when ACTUAL is not an APIObject pointer", func(m types.GomegaMatcher) {
			_, err := m.Match("not type *APIObject")
			Expect(err).To(MatchError(ErrMatcherExpectsAPIObjectPointer))
		},
			Entry("Created", Created(0)),
			Entry("Updated", Updated(0)),
			Entry("Destroyed", Destroyed(0)),
		)

		DescribeTable("has (Negated)FailureMessage", func(m types.GomegaMatcher, op string) {
			Expect(m.FailureMessage(nil)).To(Equal(fmt.Sprintf("Object has not been %s <0> times", op)))
			Expect(m.NegatedFailureMessage(nil)).To(Equal(fmt.Sprintf("Object has been %s <0> times", op)))
		},
			Entry("Created", Created(0), "created"),
			Entry("Updated", Updated(0), "updated"),
			Entry("Destroyed", Destroyed(0), "destroyed"),
		)
	})

	Context("Object", func() {
		It("can match object with same type", func() {
			a.FakeExisting(&testObject2{})
			Expect(a.All()).To(ContainElementNTimes(Object(&testObject{}), 8))
			Expect(a.All()).To(ContainElementNTimes(Object(&testObject2{}), 1))
		})

		It("returns error when match is called with uncomparable type", func() {
			var obj testObjectNotStruct
			id := a.FakeExisting(&obj)
			_, err := Object(&obj).Match(a.Inspect(id))
			Expect(err).To(MatchError(compare.ErrInvalidType))
		})

		It("has FailureMessage", func() {
			Expect(Object(&testObject{TestFieldA: "foo", TestFieldB: "bar"}, "TestFieldA", "TestFieldB").FailureMessage(a.Inspect(objects[0].Identifier))).To(HavePrefix("ACTUAL Object does not match EXPECTED Object:"))
			Expect(Object(&testObject{}).FailureMessage("not APIObject pointer")).To(Equal(ErrMatcherExpectsAPIObjectPointer.Error()))
		})

		It("has NegatedFailureMessage", func() {
			Expect(Object(&testObject{}).NegatedFailureMessage(a.Inspect(objects[0].Identifier))).To(HavePrefix("ACTUAL Object does match EXPECTED Object"))
		})
	})

	Context("Existing", func() {
		It("checks whether or not an object currently exists", func() {
			id := a.FakeExisting(&testObject{})
			Expect(Existing().Match(a.Inspect(id))).To(BeTrue())
			err := a.Destroy(context.TODO(), &testObject{Identifier: id})
			Expect(err).ToNot(HaveOccurred())
			Expect(Existing().Match(a.Inspect(id))).To(BeFalse())
		})

		It("returns error if ACTUAL is not *APIObject", func() {
			_, err := Existing().Match(nil)
			Expect(err).To(MatchError(ErrMatcherExpectsAPIObjectPointer))
		})

		It("has (Negated)FailureMessage", func() {
			Expect(Existing().FailureMessage(nil)).To(Equal("Object does not exist"))
			Expect(Existing().NegatedFailureMessage(nil)).To(Equal("Object does exist"))
		})
	})

	Context("TaggedWith", func() {
		It("checks whether or not an object was tagged with tags", func() {
			a.FakeExisting(&testObject{}, "tag-0", "tag-1", "tag-2", "tag-3")
			a.FakeExisting(&testObject{}, "tag-1", "tag-2")

			Expect(a.All()).To(ContainElementNTimes(TaggedWith("tag-1", "tag-2"), 2))
		})

		It("returns error if ACTUAL is not *APIObject", func() {
			_, err := TaggedWith().Match(nil)
			Expect(err).To(MatchError(ErrMatcherExpectsAPIObjectPointer))
		})

		It("has (Negated)FailureMessage", func() {
			id := a.FakeExisting(&testObject{}, "tag-2", "tag-3")
			Expect(TaggedWith("tag-0", "tag-1").FailureMessage(a.Inspect(id))).To(ContainSubstring("Object not tagged with all required tags"))
			Expect(TaggedWith("tag-2", "tag-3").NegatedFailureMessage(a.Inspect(id))).To(ContainSubstring("Object tagged with unexpected tags"))
		})
	})
})
