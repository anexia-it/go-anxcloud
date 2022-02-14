//go:build !integration
// +build !integration

package v1

import (
	"fmt"
	"net/http"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const (
	mockVLANIdentifier = "foobarbaz4223691337"

	waitTimeout  = 30 * time.Millisecond
	retryTimeout = 10 * time.Millisecond
)

var (
	mockServer *Server

	// Engine takes some time to assign a name, emulate with this
	mockGetAssigned = false

	// Engine takes some time to have a new VLAN active, emulate with this
	mockGetActive = false

	// After destroying, the Engine needs some time to actually delete it - emulate with this
	mockGetDeleting = false

	// This is set to true to return a 404 response when retrieving our test VLAN
	mockGetDeleted = false
)

func e2eApiClient() (api.API, error) {
	if mockServer == nil {
		mockServer = NewServer()
	}

	return api.NewAPI(
		api.WithClientOptions(
			client.BaseURL(mockServer.URL()),
			client.IgnoreMissingToken(),
		),
	)
}

func mockVLANResponseBody(desc string, vmProvisioning bool) map[string]interface{} {
	name := "newVLAN-422369"
	status := "Pending"

	if mockGetAssigned {
		name = "VLAN1000"
	}

	if mockGetDeleting {
		status = "Marked for deletion"
	} else if mockGetActive {
		status = "Active"
	}

	return map[string]interface{}{
		"identifier":           mockVLANIdentifier,
		"name":                 name,
		"description_customer": desc,
		"description_internal": "",
		"role_text":            "SCND generated",
		"status":               status,
		"locations": []map[string]interface{}{
			{
				"identifier": locationIdentifier,
				"name":       "ANX04",
			},
		},
		"vm_provisioning": vmProvisioning,
	}
}

func prepareCreate(desc string) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("POST", "/api/vlan/v1/vlan.json"),
		VerifyJSONRepresenting(map[string]interface{}{
			"description_customer": desc,
			"vm_provisioning":      false,
			"location":             locationIdentifier,
		}),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"name":                 "newVLAN-422369",
			"identifier":           mockVLANIdentifier,
			"description_customer": desc,
			"vm_provisioning":      false,
			"locations": []map[string]interface{}{
				{
					"identifier": locationIdentifier,
					"name":       "ANX04",
				},
			},
		}),
	))
}

func prepareGet(desc string, vmProvisioning bool) {
	var response http.HandlerFunc

	if mockGetDeleted {
		response = RespondWith(404, ``)
	} else {
		body := mockVLANResponseBody(desc, vmProvisioning)
		mockGetAssigned = true // after first request it's going to be assigned
		response = RespondWithJSONEncoded(200, body)
	}

	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/vlan/v1/vlan.json/"+mockVLANIdentifier),
		response,
	))
}

func prepareList(desc string, vmProvisioning bool) {
	mockVLANs := []map[string]interface{}{
		{"identifier": "foo", "name": "VLAN1337", "description_customer": "black lives matter", "vm_provisioning": true},
		{"identifier": "bar", "name": "VLAN42", "description_customer": "trans rights are human rights", "vm_provisioning": false},
		{"identifier": "blarz", "name": "VLAN23", "description_customer": "more good strings accepted via PR", "vm_provisioning": true},
		{"identifier": "aölskdjasd", "name": "VLAN3324", "description_customer": "aöäslkdjlsdkgjh.lfdknhdfg", "vm_provisioning": false},
		{"identifier": "IShouldUsePwgen", "name": "VLAN3325", "description_customer": "I really should use pwgen for this.", "vm_provisioning": true},
		{"identifier": "6 more to go", "name": "VLAN3326", "description_customer": "I need at least two pages, having our mock one on the second page to test if its iterating correctly", "vm_provisioning": false},
		{"identifier": "5 more to go", "name": "VLAN3327", "description_customer": "booooooring", "vm_provisioning": true},
		{"identifier": "4 more to go", "name": "VLAN3328", "description_customer": "hey reviewer, are you reading this?", "vm_provisioning": false},
		{"identifier": "3 more to go", "name": "VLAN3329", "description_customer": "because, if you do, I hope you are less bored", "vm_provisioning": true},
		{"identifier": "2 more to go", "name": "VLAN3330", "description_customer": "google: how to have fun mocking things", "vm_provisioning": false},
		{"identifier": "1 more to go", "name": "VLAN3331", "description_customer": "This is the last random one!", "vm_provisioning": true},
		mockVLANResponseBody(desc, vmProvisioning),
	}

	pages := [][]map[string]interface{}{
		mockVLANs[0:10],
		mockVLANs[10:],
		{},
	}

	Expect(pages[0]).To(HaveLen(10))
	Expect(pages[1]).To(HaveLen(2))
	Expect(pages[len(pages)-1]).To(HaveLen(0))

	for i, data := range pages {
		mockServer.AppendHandlers(CombineHandlers(
			VerifyRequest("GET", "/api/vlan/v1/vlan.json/filtered", fmt.Sprintf("page=%v&limit=10", i+1)),
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

func prepareUpdate(desc string, vmProvisioning bool) {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("PUT", "/api/vlan/v1/vlan.json/"+mockVLANIdentifier),
		VerifyJSONRepresenting(map[string]interface{}{
			"identifier":           mockVLANIdentifier,
			"description_customer": desc,
			"vm_provisioning":      vmProvisioning,
		}),
		RespondWithJSONEncoded(200, mockVLANResponseBody(desc, vmProvisioning)),
	))
}

func prepareDelete() {
	mockServer.AppendHandlers(CombineHandlers(
		VerifyRequest("DELETE", "/api/vlan/v1/vlan.json/"+mockVLANIdentifier),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"identifier":           nil,
			"name":                 nil,
			"description_customer": nil,
		}),
	))
}

func prepareEventuallyActive(desc string, vmProvisioning bool) {
	prepareGet(desc, vmProvisioning)
	mockGetActive = true // Active on second request
}

func prepareEventuallyDeleted(desc string, vmProvisioning bool) {
	prepareGet(desc, vmProvisioning)
	mockGetDeleted = true // Active on second request
}

func prepareDeleting() {
	mockGetDeleting = true
	prepareGet("", true) // values not important anymore
}
