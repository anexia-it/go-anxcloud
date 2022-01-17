package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go.anx.io/go-anxcloud/pkg/utils/test/gomega"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Record", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Record{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.ResponseDecodeHook", func() {
		var i types.ResponseDecodeHook
		o := Record{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.PaginationSupportHook", func() {
		var i types.PaginationSupportHook
		o := Record{}
		Expect(&o).To(ImplementInterface(&i))
	})
})

var _ = Describe("Object Zone", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Zone{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestFilterHook", func() {
		var i types.RequestFilterHook
		o := Zone{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Zone{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.ResponseFilterHook", func() {
		var i types.ResponseFilterHook
		o := Zone{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.PaginationSupportHook", func() {
		var i types.PaginationSupportHook
		o := Zone{}
		Expect(&o).To(ImplementInterface(&i))
	})
})
