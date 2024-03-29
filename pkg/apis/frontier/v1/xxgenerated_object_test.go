// Code generated by go.anx.io/go-anxcloud/tools object-generator - DO NOT EDIT!

package v1_test

import (
	. "github.com/onsi/ginkgo/v2"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	"go.anx.io/go-anxcloud/pkg/api/types"
	apipkg "go.anx.io/go-anxcloud/pkg/apis/frontier/v1"
)

var _ = Describe("Object Action", func() {
	o := apipkg.Action{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object API", func() {
	o := apipkg.API{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Deployment", func() {
	o := apipkg.Deployment{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})

var _ = Describe("Object Endpoint", func() {
	o := apipkg.Endpoint{}

	ifaces := make([]interface{}, 0, 1)
	{
		var i types.Object
		ifaces = append(ifaces, &i)
	}

	testutils.ObjectTests(&o, ifaces...)
})
