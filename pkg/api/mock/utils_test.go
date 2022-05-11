package mock

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mock utils", func() {
	Context("merge", func() {
		var (
			res interface{}
			err error
		)

		It("can merge primitives", func() {
			type test struct {
				String string `json:"string,omitempty"`
				Int    int    `json:"int,omitempty"`
			}

			res, err = merge(&test{}, &test{String: "test", Int: 1})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: "test", Int: 1}))

			res, err = merge(&test{Int: 1}, &test{String: "test"})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: "test", Int: 1}))

			res, err = merge(&test{String: "test", Int: 1}, &test{})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: "test", Int: 1}))

			res, err = merge(&test{String: "test", Int: 1}, &test{Int: 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: "test", Int: 2}))
		})

		It("can merge primitive pointer", func() {
			type test struct {
				String *string `json:"string,omitempty"`
				Int    *int    `json:"int,omitempty"`
			}

			s := "test"
			i := 1

			iChanged := 2

			res, err = merge(&test{}, &test{String: &s, Int: &i})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: &s, Int: &i}))

			res, err = merge(&test{Int: &i}, &test{String: &s})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: &s, Int: &i}))

			res, err = merge(&test{String: &s, Int: &i}, &test{})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: &s, Int: &i}))

			res, err = merge(&test{String: &s, Int: &i}, &test{Int: &iChanged})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{String: &s, Int: &iChanged}))

			Expect(s).To(Equal("test"))
			Expect(i).To(Equal(1))
		})

		It("can merge nested structs", func() {
			type nested struct {
				String string `json:"string,omitempty"`
			}
			type test struct {
				Nested nested `json:"nested,omitempty"`
			}

			res, err = merge(&test{}, &test{Nested: nested{String: "test"}})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{Nested: nested{String: "test"}}))

			res, err = merge(&test{Nested: nested{String: "foo"}}, &test{Nested: nested{String: "test"}})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{Nested: nested{String: "test"}}))

			res, err = merge(&test{Nested: nested{String: "foo"}}, &test{})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(&test{Nested: nested{String: "foo"}}))

		})
	})

	Context("makeObjectIdentifiable helper", func() {
		It("sets and returns identifier when identifier not already set", func() {
			obj := &testObject{}
			id := makeObjectIdentifiable(obj)

			Expect(id).ToNot(BeZero())
			Expect(id).To(Equal(obj.Identifier))
		})

		It("returns already set identifier", func() {
			obj := &testObject{Identifier: "test-identifier"}
			id := makeObjectIdentifiable(obj)
			Expect(id).To(Equal("test-identifier"))
			Expect(obj.Identifier).To(Equal("test-identifier"))
		})
	})
})
