//go:build !integration
// +build !integration

package v1_test

import (
	"fmt"
	"net/http"

	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const (
	mockPrefixIdentifier = "prefix-foobarbaz4223691337"
)

var (
	// Engine takes some time to assign a name, emulate with this
	mockPrefixGetAssigned = false

	// Engine takes some time to have a new VLAN active, emulate with this
	mockPrefixGetActive = false

	// After destroying, the Engine needs some time to actually delete it - emulate with this
	mockPrefixGetDeleting = false

	// This is set to true to return a 404 response when retrieving our test VLAN
	mockPrefixGetDeleted = false
)

func mockPrefixResponseBody(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) map[string]interface{} {
	netmask := netmaskForFamily(fam)
	status := "Pending"
	prefix := "newPrefix422369"

	if mockPrefixGetAssigned {
		switch fam {
		case ipamv1.FamilyIPv4:
			prefix = "10.244.0.0"
		case ipamv1.FamilyIPv6:
			prefix = "2001:db8:1337::"
		}
	}

	name := fmt.Sprintf("%v/%v", prefix, netmask)

	if mockPrefixGetDeleting {
		status = "Marked for deletion"
	} else if mockPrefixGetActive {
		status = "Active"
	}

	return map[string]interface{}{
		"name":                 name,
		"description_customer": desc,
		"identifier":           mockPrefixIdentifier,
		"netmask":              netmask,
		"version":              fam,
		"type":                 space,
		"status":               status,
		"router_redundancy":    false,
		"locations": []map[string]interface{}{
			{
				"identifier": locationIdentifier,
				"name":       "ANX04",
			},
		},
		"vlans": []map[string]interface{}{
			{
				"identifier": vlanIdentifier,
			},
		},
	}
}

func preparePrefixCreate(desc string, createEmpty *bool, fam ipamv1.Family, space ipamv1.AddressSpace) {
	expected := map[string]interface{}{
		"description_customer": desc,
		"netmask":              netmaskForFamily(fam),
		"version":              fam,
		"type":                 space,
		"location":             locationIdentifier,
		"vlan":                 vlanIdentifier,
		"router_redundancy":    false,
	}

	if createEmpty != nil {
		expected["create_empty"] = *createEmpty
	}

	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("POST", "/api/ipam/v1/prefix.json"),
		VerifyJSONRepresenting(expected),
		RespondWithJSONEncoded(200, mockPrefixResponseBody(desc, fam, space)),
	))
}

func preparePrefixDelete() {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("DELETE", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"identifier":           nil,
			"name":                 nil,
			"description_customer": nil,
		}),
	))
}

func preparePrefixGet(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	var response http.HandlerFunc

	if mockPrefixGetDeleted {
		response = RespondWith(404, ``)
	} else {
		body := mockPrefixResponseBody(desc, fam, space)
		mockPrefixGetAssigned = true // after first request it's going to be assigned
		response = RespondWithJSONEncoded(200, body)
	}

	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier),
		response,
	))
}

func preparePrefixList(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	mockVLANs := []map[string]interface{}{
		{"identifier": "foo", "name": "10.0.0.0/8", "description_customer": "black lives matter", "role_text": "Default"},
		{"identifier": "bar", "name": "10.1.0.0/8", "description_customer": "trans rights are human rights", "role_text": "Default"},
		{"identifier": "blarz", "name": "10.2.0.0/8", "description_customer": "more good strings accepted via PR", "role_text": "Default"},
		{"identifier": "aölskdjasd", "name": "10.3.0.0/8", "description_customer": "aöäslkdjlsdkgjh.lfdknhdfg", "role_text": "Default"},
		{"identifier": "IShouldUsePwgen", "name": "fd00::/64", "description_customer": "I really should use pwgen for this.", "role_text": "Default"},
		{"identifier": "6 more to go", "name": "2001:db8::/32", "description_customer": "I need at least two pages, having our mock one on the second page to test if its iterating correctly", "role_text": "Default"},
		{"identifier": "5 more to go", "name": "2001:db8:1337/48", "description_customer": "booooooring", "role_text": "Default"},
		{"identifier": "4 more to go", "name": "192.168.0.0/24", "description_customer": "hey reviewer, are you reading this?", "role_text": "Default"},
		{"identifier": "3 more to go", "name": "192.168.1.0/24", "description_customer": "because, if you do, I hope you are less bored", "role_text": "Default"},
		{"identifier": "2 more to go", "name": "172.16.0.0/16", "description_customer": "google: how to have fun mocking things", "role_text": "Default"},
		{"identifier": "1 more to go", "name": "172.17.0.224/29", "description_customer": "This is the last random one!", "role_text": "Default"},
		mockPrefixResponseBody(desc, fam, space),
	}

	pages := [][]map[string]interface{}{
		mockVLANs[0:10],
		mockVLANs[10:],
	}

	Expect(pages[0]).To(HaveLen(10))
	Expect(pages[1]).To(HaveLen(2))

	for i, data := range pages {
		mockServer.AppendHandlers(CombineHandlers(
			VerifyRequest("GET", "/api/ipam/v1/prefix/filtered.json", fmt.Sprintf("page=%v&limit=10", i+1)),
			RespondWithJSONEncoded(200, map[string]interface{}{
				"page":        i + 1,
				"total_pages": len(pages),
				"total_items": len(mockVLANs),
				"limit":       len(data),
				"data":        data,
			}),
		))
	}
}

func preparePrefixEventuallyActive(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	preparePrefixGet(desc, fam, space)
	mockPrefixGetActive = true // Active on second request
}

func preparePrefixUpdate(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("PUT", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier),
		VerifyJSONRepresenting(map[string]interface{}{
			"identifier":           mockPrefixIdentifier,
			"description_customer": desc,
			"router_redundancy":    false,
		}),
		RespondWithJSONEncoded(200, mockPrefixResponseBody(desc, fam, space)),
	))
}

func preparePrefixEventuallyDeleted(desc string, fam ipamv1.Family, space ipamv1.AddressSpace) {
	mockPrefixGetDeleting = true
	preparePrefixGet(desc, fam, space)
	mockPrefixGetDeleted = true // Active on second request
}
