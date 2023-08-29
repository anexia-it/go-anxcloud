package filter

import (
	"net/url"
	"reflect"
	"testing"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("filter.Helper", func() {
	type aliasedString string
	type aliasedInt int
	type aliasedFloat64 float64
	type aliasedBool bool

	DescribeTable("supports filtering on type",
		func(val interface{}, expected string) {
			testStructType := reflect.StructOf([]reflect.StructField{
				{
					Name: "FilterValue",
					Type: reflect.TypeOf(val),
					Tag:  `anxcloud:"filterable,filter"`,
				},
			})

			testObject := reflect.New(testStructType)
			testObject.Elem().Field(0).Set(reflect.ValueOf(val))

			fh, err := NewHelper(testObject.Interface())
			Expect(err).NotTo(HaveOccurred())

			q := fh.BuildQuery()
			Expect(q).To(HaveKey("filter"))
			Expect(q.Get("filter")).To(Equal(expected))
		},
		Entry("string", "foo", "foo"),
		Entry("int", int(42), "42"),
		Entry("float64", float64(42.23), "42.23"),
		Entry("bool", true, "true"),

		Entry("type alias from string", aliasedString("foo"), "foo"),
		Entry("type alias from int", aliasedInt(42), "42"),
		Entry("type alias from float64", aliasedFloat64(42.23), "42.23"),
		Entry("type alias from bool", aliasedBool(true), "true"),

		Entry("example testObject", testObject{Identifier: "foobarbaz"}, "foobarbaz"),
	)

	DescribeTable("supports filtering on single array entry",
		func(data []string, expectedQuery string, expectedError error) {
			type t struct {
				// also tests empty name override field in tag
				Foo []string `json:"bar" anxcloud:"filterable,,single"`
			}

			fh, err := NewHelper(t{Foo: data})

			if expectedError != nil {
				Expect(err).To(MatchError(expectedError))
			} else {
				Expect(err).NotTo(HaveOccurred())

				q := fh.BuildQuery()

				if expectedQuery == "" {
					Expect(q).To(BeEmpty())
				} else {
					Expect(q).To(Equal(url.Values{"bar": []string{expectedQuery}}))
				}
			}
		},
		Entry("no filter value set", []string{}, "", nil),
		Entry("one filter value set", []string{"what does the fox say?"}, "what does the fox say?", nil),
		Entry("multiple filter values set", []string{"what does the fox say?", "eat the rich"}, "", types.ErrInvalidFilter),
	)

	Context("testing against the example testObject", func() {
		var o types.Object

		BeforeEach(func() {
			o = &testObject{}
		})

		It("returns the correct error when trying to retrieve an unknown field", func() {
			helper, err := NewHelper(o)
			Expect(err).NotTo(HaveOccurred())

			_, ok, err := helper.Get("fooooNotAConfiguredFilter")
			Expect(ok).NotTo(BeTrue())
			Expect(err).To(MatchError(ErrUnknownField))
		})

		Context("with no filters set", func() {
			It("gives an empty query string", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				query := helper.BuildQuery().Encode()
				Expect(query).To(Equal(""))
			})

			It("returns not-present for every filterable field", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				for field := range helper.(filterHelper).fields {
					_, ok, err := helper.Get(field)
					Expect(ok).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				}
			})
		})

		Context("with identifier set in referenced Object", func() {
			JustBeforeEach(func() {
				o.(*testObject).Parent.Identifier = "foo"
			})

			It("gives the expected query string", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				query := helper.BuildQuery().Encode()
				Expect(query).To(Equal("parent=foo"))
			})

			It("returns the filter value and it being set", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				val, present, err := helper.Get("parent")
				Expect(err).NotTo(HaveOccurred())
				Expect(present).To(BeTrue())
				Expect(val).To(Equal("foo"))
			})
		})

		Context("with identifier set in pointer-referenced Object", func() {
			JustBeforeEach(func() {
				o.(*testObject).PointerParent = &parentObject{
					Identifier: "pointerfoo",
				}
			})

			It("gives the expected query string", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				query := helper.BuildQuery().Encode()
				Expect(query).To(Equal("ptrParent=pointerfoo"))
			})

			It("returns the filter value and it being set", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				val, present, err := helper.Get("ptrParent")
				Expect(err).NotTo(HaveOccurred())
				Expect(present).To(BeTrue())
				Expect(val).To(Equal("pointerfoo"))
			})
		})

		Context("with identifier set in PartialResource", func() {
			JustBeforeEach(func() {
				o.(*testObject).Partial = common.PartialResource{
					Identifier: "partialfoo",
				}
			})

			It("gives the expected query string", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				query := helper.BuildQuery().Encode()
				Expect(query).To(Equal("partial=partialfoo"))
			})

			It("returns the filter value and it being set", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				val, present, err := helper.Get("partial")
				Expect(err).NotTo(HaveOccurred())
				Expect(present).To(BeTrue())
				Expect(val).To(Equal("partialfoo"))
			})
		})

		Context("with description set", func() {
			JustBeforeEach(func() {
				desc := "some random description"
				o.(*testObject).Description = &desc
			})

			It("gives the expected query string", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				query := helper.BuildQuery().Encode()
				Expect(query).To(Equal("desc=some+random+description"))
			})

			It("returns the filter value and it being set", func() {
				helper, err := NewHelper(o)
				Expect(err).NotTo(HaveOccurred())

				val, present, err := helper.Get("desc")
				Expect(err).NotTo(HaveOccurred())
				Expect(present).To(BeTrue())
				Expect(val).To(Equal("some random description"))
			})
		})
	})
})

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "filter test suite")
}
