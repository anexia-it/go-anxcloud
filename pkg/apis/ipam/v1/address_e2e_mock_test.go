//go:build !integration
// +build !integration

package v1_test

import (
	"fmt"
	"math"
	"net"

	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const (
	mockAddressIdentifier = "address-foobarbaz4223691337"
)

func mockBasicAddressResponseBody(p ipamv1.Prefix, desc string, ip string) map[string]interface{} {
	ret := map[string]interface{}{
		"identifier":           testutils.TestResourceName(),
		"name":                 ip,
		"description_customer": desc,
		"version":              p.Version,
		"role_text":            "Default",
		"status":               "Active",
		"prefix":               p.Identifier,
	}

	if len(p.VLANs) == 1 {
		ret["vlan"] = p.VLANs[0].Identifier
	}

	return ret
}

func mockAddressResponseBody(p ipamv1.Prefix, desc string, ip net.IP) map[string]interface{} {
	status := "Pending"

	if mockPrefixGetDeleting {
		status = "Marked for deletion"
	} else if mockPrefixGetActive {
		status = "Active"
	}

	ret := mockBasicAddressResponseBody(p, desc, ip.String())
	ret["identifier"] = mockAddressIdentifier
	ret["status"] = status

	return ret
}

func prepareAddressCreate(p ipamv1.Prefix, desc string, ip net.IP) {
	exp := map[string]interface{}{
		"name":                 ip,
		"description_customer": desc,
		"version":              p.Version,
		"prefix":               p.Identifier,
	}

	if len(p.VLANs) == 1 {
		exp["vlan"] = p.VLANs[0].Identifier
	}

	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("POST", "/api/ipam/v1/address.json"),
		VerifyJSONRepresenting(exp),
		RespondWithJSONEncoded(200, mockAddressResponseBody(p, desc, ip)),
	))
}

func prepareAddressGet(p ipamv1.Prefix, desc string, ip net.IP) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/ipam/v1/address.json/"+mockAddressIdentifier),
		RespondWithJSONEncoded(200, mockAddressResponseBody(p, desc, ip)),
	))
}

func prepareAddressList(p ipamv1.Prefix, shouldEmpty bool, desc string, ip net.IP) {
	baseIP, _, err := net.ParseCIDR(p.Name)
	Expect(err).NotTo(HaveOccurred(), "expected a parsable prefix")

	var mockAddresses []map[string]interface{}

	if !shouldEmpty {
		if p.Version == ipamv1.FamilyIPv6 {
			// this would run out of memory very badly
			panic("never ever try to have an IPv6 prefix that is _not empty_")
		}

		if p.Netmask < 24 {
			panic("e2e mocks only supports netmask >= 24")
		}

		numAddresses := int(math.Exp2(float64(32 - p.Netmask)))
		mockAddresses = make([]map[string]interface{}, numAddresses)

		for i := range mockAddresses {
			ip := make(net.IP, len(baseIP))
			copy(ip, baseIP)
			ip[len(ip)-1] += byte(i)
			mockAddresses[i] = mockBasicAddressResponseBody(p, testutils.TestResourceName(), ip.String())
		}

		ipIdx := ip[len(ip)-1] - baseIP[len(baseIP)-1]
		mockAddresses[ipIdx] = mockAddressResponseBody(p, desc, ip)

	} else {
		mockAddresses = make([]map[string]interface{}, 0, 4)
		mockAddresses = append(mockAddresses, mockBasicAddressResponseBody(p, "Network address", baseIP.String()))

		gatewayIP := make(net.IP, len(baseIP))
		copy(gatewayIP, baseIP)
		gatewayIP[len(gatewayIP)-1]++
		mockAddresses = append(mockAddresses, mockBasicAddressResponseBody(p, "Gateway", gatewayIP.String()))

		mockAddresses = append(mockAddresses, mockAddressResponseBody(p, desc, ip))

		if p.Version == ipamv1.FamilyIPv4 {
			broadcastIP := make(net.IP, len(baseIP))
			copy(broadcastIP, baseIP)
			broadcastIP[len(broadcastIP)-1] += 7 // let's statically calc for a /29, good enough for the mock
			mockAddresses = append(mockAddresses, mockBasicAddressResponseBody(p, "Broadcast", broadcastIP.String()))
		}
	}

	pageCount := len(mockAddresses) / 10
	if pageCount*10 < len(mockAddresses) {
		pageCount++
	}

	Expect(pageCount).To(BeNumerically(">=", 1))

	expectedQuery := fmt.Sprintf(
		"prefix=%v&version=%v&private=%v",
		p.Identifier, p.Version, p.Type == ipamv1.AddressSpacePrivate,
	)

	for i := 0; i <= pageCount; i++ {
		pageQuery := fmt.Sprintf("%v&page=%v&limit=10", expectedQuery, i+1)

		var data []map[string]interface{}

		if i < pageCount {
			firstIdx := 10 * i
			count := 10

			if firstIdx+count > len(mockAddresses) {
				count = len(mockAddresses) - firstIdx
			}

			data = mockAddresses[firstIdx:count]
		} else {
			data = make([]map[string]interface{}, 0)
		}

		mockServer.AppendHandlers(CombineHandlers(
			VerifyRequest("GET", "/api/ipam/v1/address/filtered.json", pageQuery),
			RespondWithJSONEncoded(200, map[string]interface{}{
				"page":        i + 1,
				"total_pages": pageCount,
				"total_items": len(mockAddresses),
				"limit":       10,
				"data":        data,
			}),
		))
	}
}

func prepareAddressUpdate(p ipamv1.Prefix, desc string, ip net.IP) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("PUT", "/api/ipam/v1/address.json/"+mockAddressIdentifier),
		VerifyJSONRepresenting(map[string]interface{}{
			"identifier":           mockAddressIdentifier,
			"description_customer": desc,
		}),
		RespondWithJSONEncoded(200, mockAddressResponseBody(p, desc, ip)),
	))
}

func prepareAddressDelete(p ipamv1.Prefix, desc string, ip net.IP) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("DELETE", "/api/ipam/v1/address.json/"+mockAddressIdentifier),
		RespondWithJSONEncoded(200, mockAddressResponseBody(p, desc, ip)),
	))
}
