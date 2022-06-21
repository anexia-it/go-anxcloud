package v1

import (
	. "github.com/onsi/ginkgo/v2"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Template", func() {
	o := Template{}

	ifaces := make([]interface{}, 0, 4)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.PaginationSupportHook
		ifaces = append(ifaces, &i)
	}
	{
		var i types.FilterRequestURLHook
		ifaces = append(ifaces, &i)
	}
	{
		var i types.ResponseDecodeHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})
