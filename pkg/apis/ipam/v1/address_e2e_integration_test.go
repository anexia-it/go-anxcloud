//go:build integration

package v1_test

import (
	"go.anx.io/go-anxcloud/pkg/api"

	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
)

func prepareAddressCreate(a api.API, addr *ipamv1.Address) {
}

func prepareAddressGet(a api.API, addr ipamv1.Address) {
}

func prepareAddressDelete(a api.API, addr ipamv1.Address) {
}

func prepareAddressUpdate(a api.API, addr ipamv1.Address, newDescription string) {
}

func prepareAddressList(a api.API, prefix ipamv1.Prefix, addr ipamv1.Address) {
}
