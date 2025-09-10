package v1_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	clouddnsv1 "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"
)

type randomTimes struct {
	refresh int
	retry   int
	expire  int
	ttl     int
}

func cleanupZone(a api.API, zone string) error {
	if !isIntegrationTest {
		return nil
	}

	retryAllowance := 5
	retryTime := 5 * time.Second

	for {
		err := a.Destroy(context.TODO(), &clouddnsv1.Zone{Name: zone})
		if err == nil {
			return nil
		}

		re := api.HTTPError{}
		if isResponseError := errors.As(err, &re); isResponseError && re.StatusCode() == 404 {
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

func ensureTestRecord(a api.API, record clouddnsv1.Record) string {
	if !isIntegrationTest {
		return uuid.Nil.String()
	}

	err := a.Create(context.TODO(), &record)
	Expect(err).NotTo(HaveOccurred())

	channel := make(types.ObjectChannel)
	err = a.List(context.TODO(), &record, api.ObjectChannel(&channel))
	if err != nil {
		Fail("Couldn't get list of records")
		return uuid.Nil.String()
	}

	r := clouddnsv1.Record{}
	for res := range channel {
		err = res(&r)
		if err != nil {
			Fail("Couldn't get record from channel")
			return uuid.Nil.String()
		}

		if r.Name == record.Name && r.Type == record.Type {
			return r.Identifier
		}
	}

	Fail("Identifier of created record could not be found")
	return uuid.Nil.String()
}

var _ = Describe("CloudDNS API client", Ordered, func() {
	var zoneName string
	var times randomTimes
	var a api.API

	BeforeAll(func() {
		zoneName = test.RandomHostname() + "-genclient.go-anxcloud.test"

		rng := rand.New(rand.NewSource(GinkgoRandomSeed())) // #nosec G404 - test data generation doesn't need cryptographic randomness

		a, _ = getApi()

		times = randomTimes{
			rng.Intn(10) * 100,
			rng.Intn(10) * 100,
			rng.Intn(10) * 1000,
			rng.Intn(10) * 100,
		}

		DeferCleanup(func() {
			if zoneName != "" {
				if err := cleanupZone(a, zoneName); err != nil {
					GinkgoLogr.Error(err, "Error deleting zone")
				}
			}
		})
	})

	Context("without the zone existing", func() {
		It("should make a valid create zone request", func() {
			zone := clouddnsv1.Zone{
				Name:       zoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "admin@" + zoneName,
				Refresh:    times.refresh,
				Retry:      times.retry,
				Expire:     times.expire,
				TTL:        times.ttl,
			}

			mock_create_zone(zone)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			err := a.Create(ctx, &zone)

			Expect(err).NotTo(HaveOccurred())

			mock_expect_request_count(1)
		})
	})

	Context("with the zone existing", func() {
		CheckZoneData := func(zone clouddnsv1.Zone) {
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

			channel := make(types.ObjectChannel)
			err := a.List(ctx, &clouddnsv1.Zone{}, api.ObjectChannel(&channel))
			Expect(err).NotTo(HaveOccurred())

			zone := clouddnsv1.Zone{}
			for res := range channel {
				err = res(&zone)
				Expect(err).NotTo(HaveOccurred())
				if zone.Name == zoneName {
					CheckZoneData(zone)
				}
			}

			mock_expect_request_count(1)
		})

		It("should retrieve our zone with expected data", func() {
			mock_get_zone(zoneName, times, false)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			zone := clouddnsv1.Zone{Name: zoneName}
			err := a.Get(ctx, &zone)
			Expect(err).NotTo(HaveOccurred())

			CheckZoneData(zone)

			mock_expect_request_count(1)
		})

		It("should make a valid update zone request", func() {
			zone := clouddnsv1.Zone{
				Name:       zoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "not-the-admin@" + zoneName,
				Refresh:    times.refresh,
				Retry:      times.retry,
				Expire:     times.expire,
				TTL:        times.ttl,
			}

			mock_update_zone(zone)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			err := a.Update(ctx, &zone)

			Expect(err).NotTo(HaveOccurred())

			mock_get_zone(zoneName, times, true)

			retrievedZone := clouddnsv1.Zone{Name: zoneName}
			err = a.Get(ctx, &retrievedZone)
			Expect(err).NotTo(HaveOccurred())

			Expect(retrievedZone.AdminEmail).To(Equal("not-the-admin@" + zoneName))

			mock_expect_request_count(2)
		})

		It("should make a valid create record request", func() {
			record := clouddnsv1.Record{
				Name:     "test1",
				ZoneName: zoneName,
				Type:     "TXT",
				RData:    "test record",
				Region:   "default",
				TTL:      300,
			}

			mock_create_record(zoneName, record)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			err := a.Create(ctx, &record)
			Expect(err).NotTo(HaveOccurred())
			Expect(record.Identifier).NotTo(BeEmpty())
			Expect(record.RData).To(Equal(`"test record"`))

			mock_expect_request_count(1)
		})

		It("should return an error when trying to create record with empty name", func() {
			record := clouddnsv1.Record{
				Name:     "",
				ZoneName: zoneName,
				Type:     "TXT",
				RData:    "test record",
				Region:   "default",
				TTL:      300,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			err := a.Create(ctx, &record)
			Expect(err).To(MatchError(clouddnsv1.ErrEmptyRecordNameNotSupported))

			mock_expect_request_count(0)
		})

		Context("with the record already existing", Ordered, func() {
			var identifier string

			BeforeEach(func() {
				record := clouddnsv1.Record{
					Name:     "test2",
					ZoneName: zoneName,
					Type:     "TXT",
					RData:    "test record",
					Region:   "default",
					TTL:      300,
				}

				identifier = ensureTestRecord(a, record)
			})

			It("should make a valid delete record request", func() {
				mock_delete_record(zoneName, identifier)

				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				err := a.Destroy(ctx, &clouddnsv1.Record{Identifier: identifier, ZoneName: zoneName})
				Expect(err).NotTo(HaveOccurred())

				mock_expect_request_count(1)
			})

			// Have the update test after the delete test in our Ordered
			// context to ensure we don't have a record with the given name
			// anymore when creating it for this test and without having to
			// find the identifier of the record after it was updated to then
			// delete it.
			It("should make a valid update record request", func() {
				record := clouddnsv1.Record{
					Identifier: identifier,
					Name:       "test2",
					ZoneName:   zoneName,
					Type:       "TXT",
					RData:      "test record",
					Region:     "default",
					TTL:        300,
				}

				mock_update_record(zoneName, identifier, record)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				err := a.Update(ctx, &record)
				Expect(err).NotTo(HaveOccurred())

				mock_expect_request_count(1)
			})
		})

		Context("with some records existing", Ordered, func() {
			BeforeAll(func() {
				_ = ensureTestRecord(a, clouddnsv1.Record{
					Name:     "test3",
					ZoneName: zoneName,
					Type:     "TXT",
					RData:    "test record",
					Region:   "default",
					TTL:      300,
				})

				_ = ensureTestRecord(a, clouddnsv1.Record{
					Name:     "@",
					ZoneName: zoneName,
					Type:     "A",
					RData:    "127.0.0.1",
					Region:   "default",
					TTL:      300,
				})

				_ = ensureTestRecord(a, clouddnsv1.Record{
					Name:     "@",
					ZoneName: zoneName,
					Type:     "AAAA",
					RData:    "::1",
					Region:   "default",
					TTL:      300,
				})

				_ = ensureTestRecord(a, clouddnsv1.Record{
					Name:     "www",
					ZoneName: zoneName,
					Type:     "A",
					RData:    "127.0.0.1",
					Region:   "default",
					TTL:      300,
				})

				_ = ensureTestRecord(a, clouddnsv1.Record{
					Name:     "www",
					ZoneName: zoneName,
					Type:     "AAAA",
					RData:    "::1",
					Region:   "default",
					TTL:      300,
				})
			})

			It("lists all records of the zone", func() {
				mock_list_records(zoneName)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				channel := make(types.ObjectChannel)
				err := a.List(ctx, &clouddnsv1.Record{ZoneName: zoneName}, api.ObjectChannel(&channel))
				Expect(err).NotTo(HaveOccurred())

				r := clouddnsv1.Record{}
				rmap := make(map[string]map[string]clouddnsv1.Record, 3)
				for res := range channel {
					err := res(&r)
					Expect(err).NotTo(HaveOccurred())
					if _, ok := rmap[r.Name]; !ok {
						rmap[r.Name] = make(map[string]clouddnsv1.Record, 2)
					}

					rmap[r.Name][r.Type] = r
				}

				Expect(rmap).To(HaveKey("@"))
				Expect(rmap).To(HaveKey("www"))
				Expect(rmap).To(HaveKey("test3"))

				Expect(rmap["@"]).To(HaveKey("A"))
				Expect(rmap["@"]).To(HaveKey("AAAA"))
				Expect(rmap["www"]).To(HaveKey("A"))
				Expect(rmap["www"]).To(HaveKey("AAAA"))
				Expect(rmap["test3"]).To(HaveKey("TXT"))

				Expect(rmap["@"]["A"].RData).To(Equal("127.0.0.1"))
				Expect(rmap["www"]["A"].RData).To(Equal("127.0.0.1"))
				Expect(rmap["@"]["AAAA"].RData).To(Equal("::1"))
				Expect(rmap["www"]["AAAA"].RData).To(Equal("::1"))

				Expect(rmap["test3"]["TXT"].RData).To(Equal("\"test record\"")) // I love the engine. Mara @LittleFox94 Grosch

				Expect(rmap["@"]["A"].ZoneName).To(Equal(zoneName))
				Expect(rmap["www"]["A"].ZoneName).To(Equal(zoneName))
				Expect(rmap["@"]["AAAA"].ZoneName).To(Equal(zoneName))
				Expect(rmap["www"]["AAAA"].ZoneName).To(Equal(zoneName))

				Expect(rmap["test3"]["TXT"].ZoneName).To(Equal(zoneName))

				mock_expect_request_count(1)
			})

			It("searches for and finds specific records in the zone", func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				channel := make(types.ObjectChannel)

				mock_search_records_by_name(zoneName, "www")
				err := a.List(ctx, &clouddnsv1.Record{ZoneName: zoneName, Name: "www"}, api.ObjectChannel(&channel))
				Expect(err).NotTo(HaveOccurred())
				mock_expect_request_count(1)

				r := clouddnsv1.Record{}
				recordCount := 0
				for res := range channel {
					err := res(&r)
					recordCount += 1
					Expect(err).NotTo(HaveOccurred())

					Expect(r.Name).To(Equal("www"))
				}
				Expect(recordCount).To(Equal(2))

				mock_search_records_by_rdata(zoneName, "::1")
				channel = make(types.ObjectChannel)
				err = a.List(ctx, &clouddnsv1.Record{ZoneName: zoneName, RData: "::1"}, api.ObjectChannel(&channel))
				Expect(err).NotTo(HaveOccurred())
				mock_expect_request_count(1)

				recordCount = 0
				for res := range channel {
					err := res(&r)
					recordCount += 1
					Expect(err).NotTo(HaveOccurred())

					Expect(r.RData).To(Equal("::1"))
				}
				Expect(recordCount).To(Equal(2))

				mock_search_records_by_type(zoneName, "TXT")
				channel = make(types.ObjectChannel)
				err = a.List(ctx, &clouddnsv1.Record{ZoneName: zoneName, Type: "TXT"}, api.ObjectChannel(&channel))
				Expect(err).NotTo(HaveOccurred())
				mock_expect_request_count(1)

				recordCount = 0
				for res := range channel {
					err := res(&r)
					recordCount += 1
					Expect(err).NotTo(HaveOccurred())

					Expect(r.Type).To(Equal("TXT"))
				}
				Expect(recordCount).To(Equal(3))

				mock_search_records_by_all(zoneName, "www", "127.0.0.1", "A")
				channel = make(types.ObjectChannel)
				err = a.List(ctx, &clouddnsv1.Record{ZoneName: zoneName, Name: "www", RData: "127.0.0.1", Type: "A"}, api.ObjectChannel(&channel))
				Expect(err).NotTo(HaveOccurred())
				mock_expect_request_count(1)

				recordCount = 0
				for res := range channel {
					err := res(&r)
					recordCount += 1
					Expect(err).NotTo(HaveOccurred())

					Expect(r.Name).To(Equal("www"))
					Expect(r.RData).To(Equal("127.0.0.1"))
					Expect(r.Type).To(Equal("A"))
				}
				Expect(recordCount).To(Equal(1))
			})
		})

		It("should make a valid delete zone request", func() {
			mock_delete_zone(zoneName)

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()

			zone := clouddnsv1.Zone{Name: zoneName}
			err := a.Destroy(ctx, &zone)

			Expect(err).NotTo(HaveOccurred())

			mock_expect_request_count(1)

			zoneName = ""
		})
	})
})
