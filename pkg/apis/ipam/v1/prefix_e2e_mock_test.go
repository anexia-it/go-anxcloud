//go:build !integration

package v1_test

import (
	"fmt"
	"net/http"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

// The functions below are wrappers around the API internal state. I implemented
// them to provide an easy way for changing the API behaviour without having to
// extend the actual API interface through embedding.
func preparePrefixCreate(a api.API, p ipamv1.Prefix) { a.(*mockAPI).preparePrefixCreate(p) }
func preparePrefixGet(a api.API, p ipamv1.Prefix)    { a.(*mockAPI).preparePrefixGet(p) }
func preparePrefixDelete(a api.API)                  { a.(*mockAPI).preparePrefixDelete() }
func preparePrefixUpdate(a api.API, p ipamv1.Prefix, newDescription string) {
	a.(*mockAPI).preparePrefixUpdate(p, newDescription)
}
func preparePrefixList(a api.API, p ipamv1.Prefix) { a.(*mockAPI).preparePrefixList(p) }

func (m *mockAPI) mockPrefixResponseBody(p ipamv1.Prefix) map[string]any {
	var prefix string
	switch p.Version {
	case ipamv1.FamilyIPv4:
		prefix = mockPrefixV4.String()
	case ipamv1.FamilyIPv6:
		prefix = mockPrefixV6.String()
	}

	name := fmt.Sprintf("%v/%v", prefix, p.Netmask)

	var status ipamv1.Status
	switch {
	case m.statusFlags.SetDeleting:
		status = ipamv1.StatusMarkedForDeletion
	case m.statusFlags.SetActive:
		status = ipamv1.StatusActive
	default:
		status = ipamv1.StatusPending
	}

	return map[string]any{
		"name":                 name,
		"description_customer": p.DescriptionCustomer,
		"identifier":           mockPrefixIdentifier,
		"netmask":              p.Netmask,
		"version":              p.Version,
		"type":                 p.Type,
		"status":               status,
		"router_redundancy":    false,
		"locations": []map[string]any{
			{
				"identifier": locationIdentifier,
				"name":       "ANX04",
			},
		},
		"vlans": []map[string]any{{"identifier": vlanIdentifier}},
	}
}

func (m *mockAPI) preparePrefixCreate(p ipamv1.Prefix) {
	expected := map[string]any{
		"description_customer": p.DescriptionCustomer,
		"netmask":              p.Netmask,
		"version":              p.Version,
		"type":                 p.Type,
		"location":             locationIdentifier,
		"vlan":                 vlanIdentifier,
	}

	m.srv.RouteToHandler("POST", "/api/ipam/v1/prefix.json", CombineHandlers(
		VerifyJSONRepresenting(expected),
		RespondWithJSONEncoded(200, m.mockPrefixResponseBody(p)),
	))
}

func (m *mockAPI) preparePrefixDelete() {
	m.srv.RouteToHandler("DELETE", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier,
		CombineHandlers(
			func(w http.ResponseWriter, r *http.Request) { m.statusFlags.SetDeleting = true },
			RespondWithJSONEncoded(200, map[string]any{
				"identifier":           nil,
				"name":                 nil,
				"description_customer": nil,
			}),
		),
	)
}

func (m *mockAPI) preparePrefixGet(p ipamv1.Prefix) {
	m.srv.RouteToHandler("GET", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier, func(w http.ResponseWriter, r *http.Request) {
		var response http.HandlerFunc

		switch {
		case m.statusFlags.SetDeleted:
			response = RespondWith(http.StatusNotFound, ``)
		case m.statusFlags.SetDeleting:
			m.statusFlags.SetDeleted = true
			body := m.mockPrefixResponseBody(p)
			response = RespondWithJSONEncoded(200, body)
		default:
			body := m.mockPrefixResponseBody(p)
			m.statusFlags.SetActive = true
			response = RespondWithJSONEncoded(200, body)
		}

		response.ServeHTTP(w, r)
	})
}

func (m *mockAPI) preparePrefixList(last ipamv1.Prefix) {
	mockVLANs := []map[string]any{
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
		m.mockPrefixResponseBody(last),
	}

	pages := [][]map[string]any{
		mockVLANs[0:10],
		mockVLANs[10:],
	}

	Expect(pages[0]).To(HaveLen(10))
	Expect(pages[1]).To(HaveLen(2))

	m.srv.RouteToHandler("GET", "/api/ipam/v1/prefix/filtered.json", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		Expect(err).ToNot(HaveOccurred())

		RespondWithJSONEncoded(http.StatusOK, map[string]any{
			"page":        page,
			"total_pages": len(pages),
			"total_items": len(mockVLANs),
			"limit":       len(pages[page]),
			"data":        pages[page],
		}).ServeHTTP(w, r)
	})
}

func (m *mockAPI) preparePrefixUpdate(response ipamv1.Prefix, expectedDescription string) {
	m.srv.RouteToHandler("PUT", "/api/ipam/v1/prefix.json/"+mockPrefixIdentifier, CombineHandlers(
		VerifyJSONRepresenting(map[string]any{
			"identifier":           mockPrefixIdentifier,
			"description_customer": expectedDescription,
		}),
		RespondWithJSONEncoded(200, m.mockPrefixResponseBody(response)),
	))
}
