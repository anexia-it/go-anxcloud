package v1

import (
	. "github.com/onsi/ginkgo/v2"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Object ACL", func() {
	o := ACL{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Backend", func() {
	o := Backend{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Bind", func() {
	o := Bind{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Frontend", func() {
	o := Frontend{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object LoadBalancer", func() {
	o := LoadBalancer{}

	ifaces := make([]interface{}, 0, 3)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}
	{
		var i types.ResponseFilterHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Rule", func() {
	o := Rule{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Server", func() {
	o := Server{}

	ifaces := make([]interface{}, 0, 2)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}
	{
		var i types.RequestBodyHook
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})
