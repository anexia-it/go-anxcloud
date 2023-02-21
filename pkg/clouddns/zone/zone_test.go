package zone

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"
)

type randomTimes struct {
	refresh int
	retry   int
	expire  int
	ttl     int
}

func cleanupZone(c client.Client, zone string) error {
	if mock != nil {
		return nil
	}

	api := NewAPI(c)

	retryAllowance := 5
	retryTime := 5 * time.Second

	for {
		err := api.Delete(context.TODO(), zone)
		if err == nil {
			return nil
		}

		re := &client.ResponseError{}
		if isResponseError := errors.As(err, &re); isResponseError && re.Response.StatusCode == 404 {
			return nil
		}

		// Deleting the zone failed. Retry up to "initial value of retryAllowance" times,
		// sleeping retryTime between each retry, with retryTime being doubled after each try.
		retryAllowance--

		if retryAllowance >= 0 {
			time.Sleep(retryTime)

			retryTime = retryTime * 2
		} else {
			return fmt.Errorf("Error deleting the zone created for this test, delete zone %+v manually (last error was %+v)", zone, err)
		}
	}
}

func ensureTestRecord(c client.Client, zone string, record RecordRequest) uuid.UUID {
	if mock != nil {
		return uuid.Nil
	}

	response, err := NewAPI(c).NewRecord(context.TODO(), zone, record)
	Expect(err).NotTo(HaveOccurred())

	for _, revision := range response.Revisions {
		for _, r := range revision.Records {
			if r.Name == record.Name && r.Type == record.Type {
				return r.Identifier
			}
		}
	}

	Fail("Identifier of created record could not be found")
	return uuid.Nil
}

var _ = Describe("CloudDNS API client", Ordered, func() {
	var c client.Client

	var zoneName string
	var times randomTimes

	BeforeAll(func() {
		zoneName = test.RandomHostname() + ".go-anxcloud.test"
		c = getClient()

		rng := rand.New(rand.NewSource(GinkgoRandomSeed()))

		times = randomTimes{
			rng.Intn(10) * 100,
			rng.Intn(10) * 100,
			rng.Intn(10) * 1000,
			rng.Intn(10) * 100,
		}

		DeferCleanup(func() {
			if zoneName != "" {
				if err := cleanupZone(c, zoneName); err != nil {
					GinkgoLogr.Error(err, "Error deleting zone")
				}
			}
		})
	})

	It("should successfully create the zone", func() {
		zoneDefinition := Definition{
			Name:       zoneName,
			ZoneName:   zoneName,
			IsMaster:   true,
			DNSSecMode: "unvalidated",
			AdminEmail: "admin@" + zoneName,
			Refresh:    times.refresh,
			Retry:      times.retry,
			Expire:     times.expire,
			TTL:        times.ttl,
		}

		mock_create_zone(zoneDefinition)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		response, err := NewAPI(c).Create(ctx, zoneDefinition)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Not(BeNil()))
		Expect(response.Name).To(Equal(zoneName))

		mock_expect_request_count(1)
	})

	Context("with the zone existing", Ordered, func() {
		CheckZoneData := func(zone Zone) {
			Expect(zone.AdminEmail).To(Equal("admin@" + zoneName))
			Expect(zone.Refresh).To(Equal(times.refresh))
			Expect(zone.Retry).To(Equal(times.retry))
			Expect(zone.Expire).To(Equal(times.expire))
			Expect(zone.TTL).To(Equal(times.ttl))
		}

		It("should include our zone with expected data when listing all available zones", func() {
			mock_list_zones(zoneName, times)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			zones, err := NewAPI(c).List(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(zones).Should(WithTransform(func(zones []Zone) []string {
				names := make([]string, 0, len(zones))
				for _, zone := range zones {
					names = append(names, zone.Name)
				}
				return names
			}, ContainElement(zoneName)))

			for _, zone := range zones {
				if zone.Name == zoneName {
					CheckZoneData(zone)
				}
			}

			mock_expect_request_count(1)
		})

		It("should retrieve our zone with expected data", func() {
			mock_get_zone(zoneName, times)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			zone, err := NewAPI(c).Get(ctx, zoneName)
			Expect(err).NotTo(HaveOccurred())

			CheckZoneData(zone)

			mock_expect_request_count(1)
		})

		It("should make a valid update zone request", func() {
			zoneDefinition := Definition{
				ZoneName:   zoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "not-the-admin@" + zoneName,
				Refresh:    times.refresh,
				Retry:      times.retry,
				Expire:     times.expire,
				TTL:        times.ttl,
			}

			mock_update_zone(zoneDefinition)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := NewAPI(c).Update(ctx, zoneName, zoneDefinition)

			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response.AdminEmail).To(Equal("not-the-admin@" + zoneName))

			mock_expect_request_count(1)
		})

		It("should make a valid zone import request", func() {
			importZone := Import{
				ZoneData: `; Zone file for example.org. - region global
$ORIGIN ` + zoneName + `.
$TTL 600
@ 600 IN NS acns01.local.
@ 600 IN NS acns02.local.
@ 600 IN SOA ns0.local. admin 7 3600 1800 604800 600
www 600 IN TXT "2021-02-05 15:40:57.486411"
`,
			}

			mock_import_zone(zoneName, importZone)

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			resp, err := NewAPI(c).Import(ctx, zoneName, importZone)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())

			mock_expect_request_count(1)
		})

		It("should make a valid create record request", func() {
			record := RecordRequest{
				Name:   "test1",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			}

			mock_create_record(zoneName, record)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := NewAPI(c).NewRecord(ctx, zoneName, record)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())

			mock_expect_request_count(1)
		})

		It("should return an error when trying to create record with empty name", func() {
			record := RecordRequest{
				Name:   "",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			_, err := NewAPI(c).NewRecord(ctx, zoneName, record)
			Expect(err).To(MatchError(ErrEmptyRecordNameNotSupported))

			mock_expect_request_count(0)
		})

		Context("with an A record for test2 present", Ordered, func() {
			BeforeAll(func() {
				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "test2",
					Type:   "A",
					RData:  "127.0.0.1",
					Region: "default",
					TTL:    300,
				})
			})

			It("should make a valid changeset apply request", func() {
				changeset := ChangeSet{
					Delete: []ResourceRecord{{
						Name:   "test2",
						Type:   "A",
						Region: "default",
						RData:  "127.0.0.1",
						TTL:    300,
					}},
					Create: []ResourceRecord{{
						Name:   "test3",
						Type:   "A",
						Region: "default",
						RData:  "192.168.0.1",
						TTL:    600,
					}},
				}

				mock_apply_changeset(zoneName, changeset)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				response, err := NewAPI(c).Apply(ctx, zoneName, changeset)
				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())

				mock_expect_request_count(1)
			})
		})

		Context("with the record already existing", Ordered, func() {
			var identifier uuid.UUID
			BeforeEach(func() {
				identifier = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "test4",
					Type:   "TXT",
					RData:  "test record",
					Region: "default",
					TTL:    300,
				})
			})

			It("should make a valid delete record request", func() {
				mock_delete_record(zoneName, identifier)

				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				err := NewAPI(c).DeleteRecord(ctx, zoneName, identifier)
				Expect(err).NotTo(HaveOccurred())

				mock_expect_request_count(1)
			})

			// Have the update test after the delete test in our Ordered
			// context to ensure we don't have a record with the given name
			// anymore when creating it for this test and without having to
			// find the identifier of the record after it was updated to then
			// delete it.
			It("should make a valid update record request", func() {
				record := RecordRequest{
					Name:   "test4",
					Type:   "TXT",
					RData:  "test record",
					Region: "default",
					TTL:    300,
				}

				mock_update_record(zoneName, identifier, record)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				response, err := NewAPI(c).UpdateRecord(ctx, zoneName, identifier, record)
				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())

				mock_expect_request_count(1)
			})
		})

		Context("with some records existing", Ordered, func() {
			BeforeAll(func() {
				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "test5",
					Type:   "TXT",
					RData:  "test record",
					Region: "default",
					TTL:    300,
				})

				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "@",
					Type:   "A",
					RData:  "127.0.0.1",
					Region: "default",
					TTL:    300,
				})

				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "@",
					Type:   "AAAA",
					RData:  "::1",
					Region: "default",
					TTL:    300,
				})

				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "www",
					Type:   "A",
					RData:  "127.0.0.1",
					Region: "default",
					TTL:    300,
				})

				_ = ensureTestRecord(c, zoneName, RecordRequest{
					Name:   "www",
					Type:   "AAAA",
					RData:  "::1",
					Region: "default",
					TTL:    300,
				})
			})

			It("lists all records of the zone", func() {
				mock_list_records(zoneName)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				records, err := NewAPI(c).ListRecords(ctx, zoneName)
				Expect(err).NotTo(HaveOccurred())
				Expect(records).NotTo(BeNil())

				rmap := make(map[string]map[string]Record, 3)
				for _, r := range records {
					if _, ok := rmap[r.Name]; !ok {
						rmap[r.Name] = make(map[string]Record, 2)
					}

					rmap[r.Name][r.Type] = r
				}

				Expect(rmap).To(HaveKey("@"))
				Expect(rmap).To(HaveKey("www"))
				Expect(rmap).To(HaveKey("test5"))

				Expect(rmap["@"]).To(HaveKey("A"))
				Expect(rmap["@"]).To(HaveKey("AAAA"))
				Expect(rmap["www"]).To(HaveKey("A"))
				Expect(rmap["www"]).To(HaveKey("AAAA"))
				Expect(rmap["test5"]).To(HaveKey("TXT"))

				Expect(rmap["@"]["A"].RData).To(Equal("127.0.0.1"))
				Expect(rmap["www"]["A"].RData).To(Equal("127.0.0.1"))
				Expect(rmap["@"]["AAAA"].RData).To(Equal("::1"))
				Expect(rmap["www"]["AAAA"].RData).To(Equal("::1"))

				Expect(rmap["test5"]["TXT"].RData).To(Equal("\"test record\"")) // I love the engine. Mara @LittleFox94 Grosch

				mock_expect_request_count(1)
			})
		})

		It("should make a valid delete zone request", func() {
			mock_delete_zone(zoneName)

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()

			err := NewAPI(c).Delete(ctx, zoneName)

			Expect(err).NotTo(HaveOccurred())

			mock_expect_request_count(1)

			zoneName = ""
		})
	})
})
