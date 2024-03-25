package v1

import (
	. "github.com/onsi/ginkgo/v2"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Address", func() {
	o := Address{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Prefix", func() {
	o := Prefix{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})
