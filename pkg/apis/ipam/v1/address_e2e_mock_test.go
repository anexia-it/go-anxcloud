//go:build !integration

package v1_test

import (
	"cmp"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/gomega/ghttp"
)

func prepareAddressCreate(a api.API, addr *ipamv1.Address) {
	a.(*mockAPI).prepareAddressCreate(addr)
}

func prepareAddressGet(a api.API, addr ipamv1.Address) {
	a.(*mockAPI).prepareAddressGet(addr)
}

func prepareAddressDelete(a api.API, addr ipamv1.Address) {
	a.(*mockAPI).prepareAddressDelete(addr)
}

func prepareAddressUpdate(a api.API, addr ipamv1.Address, newDescription string) {
	a.(*mockAPI).prepareAddressUpdate(addr, newDescription)
}

func prepareAddressList(a api.API, prefix ipamv1.Prefix, addr ipamv1.Address) {
	a.(*mockAPI).prepareAddressList(prefix, addr)
}

func (m *mockAPI) mockAddressResponseBody(addr ipamv1.Address) map[string]any {
	var status ipamv1.Status
	switch {
	case m.statusFlags.SetActive:
		status = ipamv1.StatusActive
	case m.statusFlags.SetInactive:
		status = ipamv1.StatusInactive
	default:
		status = ipamv1.StatusPending
	}

	return map[string]any{
		"identifier":           cmp.Or(addr.Identifier, mockAddressIdentifier),
		"name":                 addr.Name,
		"description_customer": addr.DescriptionCustomer,
		"version":              ipamv1.FamilyIPv4,
		"role_text":            cmp.Or(addr.RoleText, "Default"),
		"status":               status,
		"prefix":               mockPrefixIdentifier,
		"vlan":                 vlanIdentifier,
	}
}

func (m *mockAPI) prepareAddressCreate(addr *ipamv1.Address) {
	m.srv.RouteToHandler("POST", "/api/ipam/v1/address.json", CombineHandlers(
		VerifyJSONRepresenting(map[string]any{
			"name":                 addr.Name,
			"description_customer": addr.DescriptionCustomer,
			"prefix":               mockPrefixIdentifier,
			"role":                 "Reserved",
		}),
		RespondWithJSONEncoded(200, m.mockAddressResponseBody(*addr)),
	))
}

func (m *mockAPI) prepareAddressGet(addr ipamv1.Address) {
	// We cannot use CombineHandlers here, as we have to compute the response body indivually on each request.
	m.srv.RouteToHandler("GET", "/api/ipam/v1/address.json/"+mockAddressIdentifier,
		func(w http.ResponseWriter, r *http.Request) {
			switch {
			case m.statusFlags.SetDeleted:
				RespondWithJSONEncoded(http.StatusNotFound, nil).ServeHTTP(w, r)
			default:
				RespondWithJSONEncoded(http.StatusOK, m.mockAddressResponseBody(addr)).ServeHTTP(w, r)
				m.statusFlags.SetActive = true
			}
		},
	)
}

func (m *mockAPI) prepareAddressList(prefix ipamv1.Prefix, testAddress ipamv1.Address) {
	p := netip.MustParsePrefix(prefix.Name)

	var (
		mockAddresses = []ipamv1.Address{testAddress}
		addr          = p.Addr()
	)
	for i := 0; i < 2; i++ {
		mockAddresses = append(mockAddresses, ipamv1.Address{
			Identifier:          testutils.RandomIdentifier(),
			Name:                addr.String(),
			DescriptionCustomer: fmt.Sprintf("Network address #%d", i),
			Version:             prefix.Version,
			RoleText:            "",
			Status:              ipamv1.StatusActive,
			VLAN:                vlanv1.VLAN{Identifier: vlanIdentifier},
			Prefix:              prefix,
			Location:            corev1.Location{Identifier: locationIdentifier},
			Type:                prefix.Type,
		})
		addr = addr.Next()
	}

	// Now we build the data that's going to be returned by mocking the responses for each address.
	type pageData map[string]any
	var data []pageData
	for _, a := range mockAddresses {
		data = append(data, m.mockAddressResponseBody(a))
	}

	const itemsPerPage = 10
	q := url.Values{}
	q.Set("prefix", prefix.Identifier)
	q.Set("version", strconv.Itoa(int(prefix.Version)))
	q.Set("private", strconv.FormatBool(prefix.Type == ipamv1.TypePrivate))
	q.Set("limit", strconv.Itoa(itemsPerPage))

	pageCount := cmp.Or(len(mockAddresses)/itemsPerPage, 1) // ensure to start at page 1
	pages := testutils.Chunk(data, itemsPerPage)
	m.srv.RouteToHandler("GET", "/api/ipam/v1/address/filtered.json",
		func(w http.ResponseWriter, r *http.Request) {
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			q.Set("page", strconv.Itoa(page)) // used for validation later

			var d []pageData
			if idx := page - 1; idx < len(pages) {
				d = pages[idx]
			}

			CombineHandlers(
				// We intentionally do not run a VerifyRequest here, since we use this handler
				// for prepareAddressGet *and* prepareAddressList.
				RespondWithJSONEncoded(http.StatusOK, pageData{
					"page":        page,
					"total_pages": pageCount,
					"total_items": len(mockAddresses),
					"limit":       itemsPerPage,
					"data":        d,
				}),
			).ServeHTTP(w, r)
		},
	)
}

func (m *mockAPI) prepareAddressUpdate(addr ipamv1.Address, desc string) {
	m.srv.RouteToHandler("PUT", "/api/ipam/v1/address.json/"+mockAddressIdentifier, CombineHandlers(
		VerifyJSONRepresenting(map[string]any{
			"identifier":           mockAddressIdentifier,
			"description_customer": desc,
			"role":                 "Reserved",
		}),
		RespondWithJSONEncoded(200, m.mockAddressResponseBody(addr)),
	))
}

func (m *mockAPI) prepareAddressDelete(addr ipamv1.Address) {
	m.srv.RouteToHandler("DELETE", "/api/ipam/v1/address.json/"+mockAddressIdentifier, CombineHandlers(
		func(w http.ResponseWriter, r *http.Request) { m.statusFlags.SetDeleted = true },
		RespondWithJSONEncoded(200, m.mockAddressResponseBody(addr)),
	))
}
