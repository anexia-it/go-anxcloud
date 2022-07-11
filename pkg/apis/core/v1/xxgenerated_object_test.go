package v1

import (
	. "github.com/onsi/ginkgo/v2"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object Location", func() {
	o := Location{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.FilterRequestURLHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Resource", func() {
	o := Resource{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.ResponseDecodeHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object ResourceWithTag", func() {
	o := ResourceWithTag{}

	ifaces := make([]interface{}, 0, 4)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.ResponseFilterHook
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}
	{
		var i types.FilterRequestURLHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})
