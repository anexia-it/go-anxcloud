package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go.anx.io/go-anxcloud/pkg/utils/test/gomega"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Backend", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Backend{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Backend{}
		Expect(&o).To(ImplementInterface(&i))
	})
})

var _ = Describe("Object Bind", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Bind{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Bind{}
		Expect(&o).To(ImplementInterface(&i))
	})
})

var _ = Describe("Object Frontend", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Frontend{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Frontend{}
		Expect(&o).To(ImplementInterface(&i))
	})
})

var _ = Describe("Object LoadBalancer", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := LoadBalancer{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := LoadBalancer{}
		Expect(&o).To(ImplementInterface(&i))
	})
})

var _ = Describe("Object Server", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Server{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Server{}
		Expect(&o).To(ImplementInterface(&i))
	})
})
