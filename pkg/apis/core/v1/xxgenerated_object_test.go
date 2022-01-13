package v1

import (
	. "github.com/anexia-it/go-anxcloud/pkg/utils/test/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Info", func() {
	It("implements the interface types.Object", func() {
		var i types.Object
		o := Info{}
		Expect(&o).To(ImplementInterface(&i))
	})
})
