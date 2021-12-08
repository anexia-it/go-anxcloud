package loadbalancer

import (
	. "github.com/anexia-it/go-anxcloud/pkg/utils/test/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Loadbalancer", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Loadbalancer{}
		Expect(&o).To(ImplementInterface(&i))
	})
	It("implements the interface types.RequestBodyHook", func() {
		var i types.RequestBodyHook
		o := Loadbalancer{}
		Expect(&o).To(ImplementInterface(&i))
	})
})
