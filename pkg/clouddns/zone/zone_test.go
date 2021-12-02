package zone

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/anexia-it/go-anxcloud/pkg/client"

	testutils "github.com/anexia-it/go-anxcloud/pkg/utils/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type randomTimes struct {
	refresh int
	retry   int
	expire  int
	ttl     int
}

var createdZones = make([]string, 0)

func ensureTestZone(c client.Client, name string, times randomTimes) {
	if mock != nil {
		return
	}

	_, err := NewAPI(c).Create(context.TODO(), Definition{
		Name:       "object for " + name,
		ZoneName:   name,
		IsMaster:   true,
		DNSSecMode: "unvalidated",
		AdminEmail: "admin@" + name,
		Refresh:    times.refresh,
		Retry:      times.retry,
		Expire:     times.expire,
		TTL:        times.ttl,
	})

	Expect(err).NotTo(HaveOccurred())

	createdZones = append(createdZones, name)
}

func cleanupZones(c client.Client) error {
	if mock != nil {
		return nil
	}

	api := NewAPI(c)

	retryAllowance := 5
	retryTime := 5 * time.Second

	for len(createdZones) > 0 {
		notDone := make([]string, 0, len(createdZones))
		lastTryErrors := make([]error, 0, len(createdZones))

		for _, zone := range createdZones {
			err := api.Delete(context.TODO(), zone)

			if err != nil {
				re := &client.ResponseError{}
				if isResponseError := errors.As(err, &re); !isResponseError || re.Response.StatusCode != 404 {
					notDone = append(notDone, zone)
					lastTryErrors = append(lastTryErrors, err)
				}
			}
		}

		createdZones = notDone
		retryAllowance--

		// Deleting of at least one zone failed. Retry up to "initial value of retryAllowance" times,
		// sleeping retryTime between each retry, with retryTime being doubled after each try.
		if len(notDone) > 0 {
			if retryAllowance >= 0 {
				time.Sleep(retryTime)

				retryTime = retryTime * 2
			} else {
				messages := make([]string, 0, len(lastTryErrors))
				for _, err := range lastTryErrors {
					messages = append(messages, err.Error())
				}

				return fmt.Errorf("Error deleting the zones created for this request, remaining zones are %+v and last errors were %+v", notDone, messages)
			}
		}
	}

	return nil
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

var _ = Describe("CloudDNS API client", func() {
	var c client.Client

	var zoneName string
	var times randomTimes

	BeforeEach(func() {
		zoneName = testutils.RandomHostname() + ".go-anxcloud.test"
		c = getClient()

		rng := rand.New(rand.NewSource(GinkgoRandomSeed()))

		times = randomTimes{
			rng.Intn(10) * 100,
			rng.Intn(10) * 100,
			rng.Intn(10) * 1000,
			rng.Intn(10) * 100,
		}
	})

	Context("without the zone existing", func() {
		It("should make a valid create zone request", func() {
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

			createdZones = append(createdZones, zoneName)
		})
	})

	Context("with the zone existing", func() {
		JustBeforeEach(func() {
			ensureTestZone(c, zoneName, times)
		})

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

		// TODO: this is broken because of ENGSUP-5233
		PIt("should make a valid update zone request", func() {
			zoneDefinition := Definition{
				Name:       "Not the ZoneName",
				ZoneName:   zoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "admin@" + zoneName,
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
			Expect(response.Name).To(Equal("Not the ZoneName"))

			mock_expect_request_count(1)
		})

		It("should make a valid delete zone request", func() {
			mock_delete_zone(zoneName)

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()

			err := NewAPI(c).Delete(ctx, zoneName)

			Expect(err).NotTo(HaveOccurred())

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
	})

	Context("with an A record for test1 present", func() {
		JustBeforeEach(func() {
			ensureTestZone(c, zoneName, times)
			_ = ensureTestRecord(c, zoneName, RecordRequest{
				Name:   "test1",
				Type:   "A",
				RData:  "127.0.0.1",
				Region: "default",
				TTL:    300,
			})
		})

		It("should make a valid changeset apply request", func() {
			changeset := ChangeSet{
				Delete: []ResourceRecord{{
					Name:   "test1",
					Type:   "A",
					Region: "default",
					RData:  "127.0.0.1",
					TTL:    300,
				}},
				Create: []ResourceRecord{{
					Name:   "test2",
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

	Context("with the record already existing", func() {
		var identifier uuid.UUID

		JustBeforeEach(func() {
			ensureTestZone(c, zoneName, times)
			identifier = ensureTestRecord(c, zoneName, RecordRequest{
				Name:   "test1",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			})
		})

		It("should make a valid update record request", func() {
			record := RecordRequest{
				Name:   "test1",
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

		It("should make a valid delete record request", func() {
			mock_delete_record(zoneName, identifier)

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			err := NewAPI(c).DeleteRecord(ctx, zoneName, identifier)
			Expect(err).NotTo(HaveOccurred())

			mock_expect_request_count(1)
		})
	})

	Context("with some records existing", func() {
		JustBeforeEach(func() {
			ensureTestZone(c, zoneName, times)

			_ = ensureTestRecord(c, zoneName, RecordRequest{
				Name:   "test1",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			})

			_ = ensureTestRecord(c, zoneName, RecordRequest{
				Name:   "",
				Type:   "A",
				RData:  "127.0.0.1",
				Region: "default",
				TTL:    300,
			})

			_ = ensureTestRecord(c, zoneName, RecordRequest{
				Name:   "",
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

			Expect(rmap).To(HaveKey(""))
			Expect(rmap).To(HaveKey("www"))
			Expect(rmap).To(HaveKey("test1"))

			Expect(rmap[""]).To(HaveKey("A"))
			Expect(rmap[""]).To(HaveKey("AAAA"))
			Expect(rmap["www"]).To(HaveKey("A"))
			Expect(rmap["www"]).To(HaveKey("AAAA"))
			Expect(rmap["test1"]).To(HaveKey("TXT"))

			Expect(rmap[""]["A"].RData).To(Equal("127.0.0.1"))
			Expect(rmap["www"]["A"].RData).To(Equal("127.0.0.1"))
			Expect(rmap[""]["AAAA"].RData).To(Equal("::1"))
			Expect(rmap["www"]["AAAA"].RData).To(Equal("::1"))

			Expect(rmap["test1"]["TXT"].RData).To(Equal("\"test record\"")) // I love the engine. Mara @LittleFox94 Grosch

			mock_expect_request_count(1)
		})
	})
})
