package v1_test

import (
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"

	. "github.com/onsi/ginkgo/v2"
)

const (
	locationIdentifier = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
)

var vlanIdentifier string

var _ = Describe("IPAM E2E tests", func() {
	testPrefix("working with private IPv4", ipamv1.FamilyIPv4, ipamv1.AddressSpacePrivate)
	testPrefix("working with private IPv6", ipamv1.FamilyIPv6, ipamv1.AddressSpacePrivate)
	testPrefix("working with public IPv4", ipamv1.FamilyIPv4, ipamv1.AddressSpacePublic)
	testPrefix("working with public IPv6", ipamv1.FamilyIPv6, ipamv1.AddressSpacePublic)
})

func netmaskForFamily(fam ipamv1.Family) int {
	switch fam {
	case ipamv1.FamilyIPv4:
		return 29
	case ipamv1.FamilyIPv6:
		return 64
	}

	panic("Invalid family")
}
