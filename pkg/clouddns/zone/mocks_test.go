package zone

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

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/api/clouddns/v1/zone.json"),
		ghttp.RespondWithJSONEncoded(200, listResponse{
			Results: []Zone{
				{
					Definition: &Definition{
						Name:       "example.com",
						AdminEmail: "root@" + zone,
						Refresh:    times.refresh * 2,
						Retry:      times.retry * 3,
						Expire:     times.expire * 4,
						TTL:        times.ttl * 4,
					},
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
					Definition: &Definition{
						Name:       zone,
						AdminEmail: "admin@" + zone,
						Refresh:    times.refresh,
						Retry:      times.retry,
						Expire:     times.expire,
						TTL:        times.ttl,
					},
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
			},
		}),
	))
}

func mock_get_zone(zone string, times randomTimes) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", fmt.Sprintf("/api/clouddns/v1/zone.json/%s", zone)),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Definition: &Definition{
				Name:       zone,
				AdminEmail: "admin@" + zone,
				Refresh:    times.refresh,
				Retry:      times.retry,
				Expire:     times.expire,
				TTL:        times.ttl,
			},
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

func mock_create_zone(def Definition) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", "/api/clouddns/v1/zone.json"),
		ghttp.VerifyJSONRepresenting(def),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Definition: &def,
		}),
	))
}

func mock_update_zone(def Definition) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("PUT", "/api/clouddns/v1/zone.json"),
		ghttp.VerifyJSONRepresenting(def),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Definition: &def,
		}),
	))
}

func mock_delete_zone(zone string) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/clouddns/v1/zone.json/%s", zone)),
		ghttp.RespondWith(200, nil),
	))
}

func mock_apply_changeset(zone string, changes ChangeSet) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/changeset", zone)),
		ghttp.VerifyJSONRepresenting(changes),
		ghttp.RespondWithJSONEncoded(200, []Record{
			{
				Name:       changes.Create[0].Name,
				Type:       changes.Create[0].Type,
				Region:     changes.Create[0].Region,
				RData:      changes.Create[0].RData,
				TTL:        &changes.Create[0].TTL,
				Identifier: uuid.NewV4(),
				Immutable:  false,
			},
		}),
	))
}

func mock_import_zone(zone string, data Import) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/import", zone)),
		ghttp.VerifyJSONRepresenting(data),
		ghttp.RespondWithJSONEncoded(200, Revision{
			CreatedAt:  time.Now(),
			Identifier: uuid.NewV4(),
			ModifiedAt: time.Now(),
			Serial:     1,
			State:      "active",
		}),
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

func mock_create_record(zone string, record RecordRequest) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", zone)),
		ghttp.VerifyJSONRepresenting(record),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Definition: &Definition{
				Name:     zone,
				ZoneName: zone,
				IsMaster: true,
			},
			Revisions: []Revision{{
				Identifier: uuid.NewV4(),
				Records: []Record{{
					Name:       record.Name,
					Type:       record.Type,
					RData:      record.RData,
					Region:     record.Region,
					TTL:        &record.TTL,
					Identifier: uuid.NewV4(),
				}},
				ModifiedAt: time.Now(),
				Serial:     1,
				State:      "active",
			}},
		}),
	))
}

func mock_update_record(zone string, recordIdentifier uuid.UUID, record RecordRequest) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("PUT", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records/%s", zone, recordIdentifier)),
		ghttp.VerifyJSONRepresenting(record),
		ghttp.RespondWithJSONEncoded(200, Zone{
			Definition: &Definition{
				Name:     zone,
				ZoneName: zone,
				IsMaster: true,
			},
			Revisions: []Revision{{
				Identifier: uuid.NewV4(),
				Records: []Record{{
					Name:       record.Name,
					Type:       record.Type,
					RData:      record.RData,
					Region:     record.Region,
					TTL:        &record.TTL,
					Identifier: recordIdentifier,
				}},
				ModifiedAt: time.Now(),
				Serial:     1,
				State:      "active",
			}},
		}),
	))
}

func mock_delete_record(zone string, recordIdentifier uuid.UUID) {
	if mock == nil {
		return
	}

	mock.server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records/%s", zone, recordIdentifier)),
		ghttp.RespondWith(200, nil),
	))
}
