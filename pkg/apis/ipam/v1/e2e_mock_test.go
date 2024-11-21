//go:build !integration

package v1_test

import (
	"net/netip"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"

	. "github.com/onsi/ginkgo/v2"
)

const (
	vlanIdentifier        = "but-why-lan"
	mockAddressIdentifier = "address-mock"
	mockPrefixIdentifier  = "prefix-foobarbaz4223691337"
)

var (
	mockPrefixV4 = netip.MustParsePrefix("10.244.0.0/29")
	mockPrefixV6 = netip.MustParsePrefix("2001:db8:1337::/64")
)

func GetTestPrefix(ipVersion ipamv1.Family, addrType ipamv1.AddressType) ipamv1.Prefix {
	GinkgoHelper()

	p := mockPrefixV4
	if ipVersion == ipamv1.FamilyIPv6 {
		p = mockPrefixV6
	}

	return ipamv1.Prefix{
		Identifier:          mockPrefixIdentifier,
		Name:                p.String(),
		DescriptionCustomer: "mockPrefix",
		Version:             ipVersion,
		Netmask:             p.Bits(),
		RoleText:            "Default",
		Status:              ipamv1.StatusActive,
		Type:                addrType,
		Locations:           []corev1.Location{{Identifier: locationIdentifier}},
		VLANs:               []vlanv1.VLAN{{Identifier: vlanIdentifier}},
	}
}
