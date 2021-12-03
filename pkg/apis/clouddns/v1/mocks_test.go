package v1

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var mock *mockserver

type mockserver struct {
	server *ghttp.Server
}

func initMockServer() {
	mock = &mockserver{
		server: ghttp.NewServer(),
	}
}

// nolint:golint,unparam
func mock_expect_request_count(count int) {
	if mock != nil {
		Expect(mock.server.ReceivedRequests()).Should(HaveLen(count))
	}
}

func mock_list_zones(zone string, times randomTimes) {
	if mock == nil {
		return
	}

	zones := make(map[string][]Zone)

	zones["results"] =
		[]Zone{
			{
				Name:       "example.com",
				AdminEmail: "root@" + zone,
				Refresh:    times.refresh * 2,
				Retry:      times.retry * 3,
				Expire:     times.expire * 4,
				TTL:        times.ttl * 4,

				Revisions: []Revision{
					{},
					{Records: []Record{}},
					{Records: []Record{
						{Name: "", Type: "A", RData: "127.0.0.2"},
					}},
					{}, {},
				},
			},
			{
				Name:       zone,
				AdminEmail: "admin@" + zone,
				Refresh:    times.refresh,
				Retry:      times.retry,
				Expire:     times.expire,
				TTL:        times.ttl,

				Revisions: []Revision{
					{}, {},
					{Records: []Record{
						{Name: "", Type: "A", RData: "127.0.0.1"},
						{Name: "", Type: "AAAA", RData: "::1"},
						{Name: "www", Type: "A", RData: "127.0.0.1"},
						{Name: "www", Type: "AAAA", RData: "::1"},
						{Name: "test1", Type: "TXT", RData: "\"test record\""},
					}},
					{},
				},
			},
		}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/api/clouddns/v1/zone.json"),
		ghttp.RespondWithJSONEncoded(200, zones),
	))
}

func mock_get_zone(zone string, times randomTimes) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s", zone)),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Name:       zone,
			AdminEmail: "admin@" + zone,
			Refresh:    times.refresh,
			Retry:      times.retry,
			Expire:     times.expire,
			TTL:        times.ttl,

			Revisions: []Revision{
				{}, {},
				{Records: []Record{
					{Name: "", Type: "A", RData: "127.0.0.1"},
					{Name: "", Type: "AAAA", RData: "::1"},
					{Name: "www", Type: "A", RData: "127.0.0.1"},
					{Name: "www", Type: "AAAA", RData: "::1"},
					{Name: "test1", Type: "TXT", RData: "\"test record\""},
				}},
				{},
			},
		}),
	))
}

func mock_create_zone(z Zone) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", "/api/clouddns/v1/zone.json"),
		ghttp.RespondWithJSONEncoded(200, z),
	))
}

func mock_update_zone(z Zone) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("PUT", "/api/clouddns/v1/zone.json"),
		ghttp.VerifyJSONRepresenting(z),
		ghttp.RespondWithJSONEncoded(200, z),
	))
}

func mock_delete_zone(zone string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/clouddns/v1/zone.json/%s", zone)),
		ghttp.RespondWith(204, nil),
	))
}

func mock_list_records(zone string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{Name: "", Type: "A", RData: "127.0.0.1"},
			{Name: "", Type: "AAAA", RData: "::1"},
			{Name: "www", Type: "A", RData: "127.0.0.1"},
			{Name: "www", Type: "AAAA", RData: "::1"},
			{Name: "test1", Type: "TXT", RData: "\"test record\""},
		}),
	))
}

func mock_search_records_by_name(zone string, name string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{Name: name, Type: "A", RData: "127.0.0.1"},
			{Name: name, Type: "AAAA", RData: "::1"},
		}),
	))
}

func mock_search_records_by_rdata(zone string, rdata string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{Name: "", Type: "AAAA", RData: rdata},
			{Name: "www", Type: "AAAA", RData: rdata},
		}),
	))
}

func mock_search_records_by_type(zone string, t string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{Name: "test1", Type: t, RData: "\"test record\""},
		}),
	))
}

func mock_search_records_by_all(zone string, name string, rdata string, t string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{Name: name, Type: t, RData: rdata},
		}),
	))
}

func mock_create_record(zone string, record Record) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.VerifyJSONRepresenting(record),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Name:     zone,
			IsMaster: true,

			Revisions: []Revision{{
				Identifier: uuid.NewV4().String(),
				Records: []Record{{
					Name:       record.Name,
					Type:       record.Type,
					RData:      record.RData,
					Region:     record.Region,
					TTL:        record.TTL,
					Identifier: uuid.NewV4().String(),
				}},
				ModifiedAt: time.Now(),
				Serial:     1,
				State:      "active",
			}},
		}),
	))
}

func mock_update_record(zone string, recordIdentifier string, record Record) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("PUT", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records/%s", zone, recordIdentifier)),
		ghttp.VerifyJSONRepresenting(record),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Name:     zone,
			IsMaster: true,
			Revisions: []Revision{{
				Identifier: uuid.NewV4().String(),
				Records: []Record{{
					Name:       record.Name,
					Type:       record.Type,
					RData:      record.RData,
					Region:     record.Region,
					TTL:        record.TTL,
					Identifier: record.Identifier,
				}},
				ModifiedAt: time.Now(),
				Serial:     1,
				State:      "active",
			}},
		}),
	))
}

func mock_delete_record(zone string, recordIdentifier string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records/%s", zone, recordIdentifier)),
		ghttp.RespondWith(204, nil),
	))
}
