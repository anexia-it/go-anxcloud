package tests_test

import (
	"context"
	"encoding/json"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"time"
)

var _ = Describe("CloudDNS API endpoint tests", func() {
	var cli client.Client

	const TestZone string = "go-sdk.test"

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Zone List Endpoint", func() {
		It("Should list all available zones", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).List(ctx)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Zone Get Endpoint", func() {
		It("Should return the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).Get(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("Definition Create Endpoint", func() {
		var createTestZoneName = "sdk-create-test.go-sdk.test"

		It("Should create the zone", func() {
			randRefresh := rand.Intn(10) * 100
			randRetry := rand.Intn(10) * 100
			randExpire := rand.Intn(10) * 1000
			randTTL := rand.Intn(10) * 100

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Zone{
					Definition: &zone.Definition{
						Name:       createTestZoneName,
						ZoneName:   createTestZoneName,
						IsMaster:   true,
						DNSSecMode: "unvalidated",
						AdminEmail: "admin@go-sdk.test",
						Refresh:    randRefresh,
						Retry:      randRetry,
						Expire:     randExpire,
						TTL:        randTTL,
					},
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			createDefinition := zone.Definition{
				ZoneName:   createTestZoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "admin@go-sdk.test",
				Refresh:    300,
				Retry:      300,
				Expire:     3600,
				TTL:        300,
			}
			response, err := zone.NewAPI(c).Create(ctx, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response.Name).To(Equal(createTestZoneName))
			Expect(response.AdminEmail).To(Equal("admin@go-sdk.test"))
		})
	})

	Context("Definition Update Endpoint", func() {
		updateTestZoneName := "sdk-update-test.go-sdk.test"

		It("Should update the zone", func() {
			randRefresh := rand.Intn(10) * 100
			randRetry := rand.Intn(10) * 100
			randExpire := rand.Intn(10) * 1000
			randTTL := rand.Intn(10) * 100

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Zone{
					Definition: &zone.Definition{
						Name:       updateTestZoneName,
						ZoneName:   updateTestZoneName,
						IsMaster:   true,
						DNSSecMode: "unvalidated",
						AdminEmail: "test@go-sdk.test",
						Refresh:    randRefresh,
						Retry:      randRetry,
						Expire:     randExpire,
						TTL:        randTTL,
					},
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			createDefinition := zone.Definition{
				ZoneName:   updateTestZoneName,
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "test@go-sdk.test",
				Refresh:    randRefresh,
				Retry:      randRetry,
				Expire:     randExpire,
				TTL:        randTTL,
			}
			response, err := zone.NewAPI(c).Update(ctx, updateTestZoneName, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.AdminEmail).To(Equal("test@go-sdk.test"))
			Expect(response.Refresh).To(Equal(randRefresh))
			Expect(response.Retry).To(Equal(randRetry))
			Expect(response.Expire).To(Equal(randExpire))
			Expect(response.TTL).To(Equal(randTTL))
		})
	})

	Context("Definition Delete Endpoint", func() {
		It("Should delete the zone", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal(http.MethodDelete))
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			err := zone.NewAPI(c).Delete(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Definition ChangeSet Endpoint", func() {
		changesetZoneName := "sdk-changeset-test.go-sdk.test"

		It("Should apply the changeset", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.ChangeSet
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := []zone.Record{{
					Identifier: uuid.NewV4(),
					Immutable:  false,
					Name:       "test2",
					RData:      "A",
					Region:     "default",
					TTL:        &ttl,
					Type:       "A",
				},
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			changeset := zone.ChangeSet{
				Delete: []zone.ResourceRecord{{
					Name:   "test1",
					Type:   "A",
					Region: "default",
					RData:  "127.0.0.1",
					TTL:    300,
				}},
				Create: []zone.ResourceRecord{{
					Name:   "test2",
					Type:   "A",
					Region: "default",
					RData:  "192.168.0.1",
					TTL:    600,
				}},
			}
			response, err := zone.NewAPI(c).Apply(ctx, changesetZoneName, changeset)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
		})
	})

	Context("Definition Import Endpoint", func() {
		importZoneName := "sdk-import-test.go-sdk.test"
		It("Should import the zone", func() {

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Import
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Revision{
					CreatedAt:  time.Now(),
					Identifier: uuid.NewV4(),
					ModifiedAt: time.Now(),
					Serial:     1,
					State:      "active",
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			importZone := zone.Import{
				ZoneData: "; Zone file for example.org. - region global\n$ORIGIN example.org.\n$TTL 600\n@ 600 IN NS acns01.local.\n@ 600 IN NS acns02.local.\n@ 600 IN SOA ns0.local. admin 7 3600 1800 604800 600\nwww 600 IN TXT \"2021-02-05 15:40:57.486411\"",
			}
			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			resp, err := zone.NewAPI(c).Import(ctx, importZoneName, importZone)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
		})
	})

	Context("Definition List Records Endpoint", func() {
		It("Should list all available records for the test zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			records, err := zone.NewAPI(cli).ListRecords(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
			Expect(records).NotTo(BeNil())
		})
	})

	Context("Definition Create Record Endpoint", func() {
		recordZoneName := "sdk-record-test.go-sdk.test"

		It("Should create the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := zone.Zone{
					Definition: &zone.Definition{
						Name:     recordZoneName,
						ZoneName: recordZoneName,
						IsMaster: true,
					},
					Revisions: []zone.Revision{{
						Identifier: uuid.NewV4(),
						Records: []zone.Record{
							{
								Identifier: uuid.NewV4(),
								Immutable:  false,
								Name:       "test1",
								RData:      "test record",
								Region:     "default",
								TTL:        &ttl,
								Type:       "TXT",
							},
						},
						Serial: 0,
						State:  "active",
					}},
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			record := zone.RecordRequest{
				Name:   "test1",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			}
			response, err := zone.NewAPI(c).NewRecord(ctx, recordZoneName, record)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
			Expect(response.Revisions).To(ContainElements())
		})
	})

	Context("Definition Update Record Endpoint", func() {
		recordZoneName := "sdk-record-test.go-sdk.test"

		It("Should update the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := zone.Zone{
					Definition: &zone.Definition{
						Name:     recordZoneName,
						ZoneName: recordZoneName,
						IsMaster: true,
					},
					Revisions: []zone.Revision{{
						Identifier: uuid.NewV4(),
						Records: []zone.Record{
							{
								Identifier: uuid.NewV4(),
								Immutable:  false,
								Name:       "test1",
								RData:      "test record",
								Region:     "default",
								TTL:        &ttl,
								Type:       "TXT",
							},
						},
						Serial: 0,
						State:  "active",
					}},
				}
				err = json.NewEncoder(w).Encode(&resp)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			record := zone.RecordRequest{
				Name:   "test1",
				Type:   "TXT",
				RData:  "test record",
				Region: "default",
				TTL:    300,
			}

			response, err := zone.NewAPI(c).UpdateRecord(ctx, recordZoneName, uuid.NewV4(), record)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
			Expect(response.Revisions).To(ContainElements())
		})
	})

	Context("Definition Delete Record Endpoint", func() {
		recordZoneName := "sdk-record-test.go-sdk.test"

		It("Should delete the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal(http.MethodDelete))
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			err := zone.NewAPI(c).DeleteRecord(ctx, recordZoneName, uuid.NewV4())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Zone Integration tests", func() {
		It("Create, update and delete a new zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout*3)
			defer cancel()

			zoneAPI := zone.NewAPI(cli)
			z, err := zoneAPI.Create(ctx, zone.Definition{
				ZoneName:   "sdk-create.test",
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "test@" + TestZone,
				Refresh:    100,
				Retry:      100,
				Expire:     600,
				TTL:        300,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(z.TTL).To(Equal(300))
			Expect(z.IsMaster).To(BeTrue())
			Expect(z.Name).To(Equal("sdk-create.test"))

			z.Definition.TTL = 600
			z.Definition.ZoneName = z.Name
			z, err = zoneAPI.Update(ctx, z.Name, *z.Definition)
			Expect(err).NotTo(HaveOccurred())
			Expect(z.TTL).To(Equal(600))

			found := false
			zones, err := zoneAPI.List(ctx)
			Expect(err).NotTo(HaveOccurred())
			for _, zoneResult := range zones {
				if zoneResult.Name == z.Name {
					Expect(z.TTL).To(BeEquivalentTo(600))
					found = true
				}
			}
			Expect(found).To(BeTrue())

			time.Sleep(time.Second * 1)

			err = zoneAPI.Delete(ctx, z.Name)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should apply a changeset to a fresh zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout*3)
			defer cancel()
			zoneAPI := zone.NewAPI(cli)

			z, err := zoneAPI.Create(ctx, zone.Definition{
				ZoneName:   "sdk-apply.test",
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "test@" + TestZone,
				Refresh:    100,
				Retry:      100,
				Expire:     600,
				TTL:        300,
			})
			Expect(err).NotTo(HaveOccurred())

			records, err := zoneAPI.Apply(ctx, z.Name, zone.ChangeSet{
				Create: []zone.ResourceRecord{{
					Name:   "test1",
					Type:   "A",
					Region: "default",
					RData:  "127.0.0.1",
					TTL:    300,
				}},
				Delete: []zone.ResourceRecord{},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(records).To(HaveLen(7))

			err = zoneAPI.Delete(ctx, z.Name)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("Should import a zone from zone data", func() {
		ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout*3)
		defer cancel()

		zoneName := "sdk-import.test"
		zoneAPI := zone.NewAPI(cli)
		z, err := zoneAPI.Create(ctx, zone.Definition{
			ZoneName:   zoneName,
			IsMaster:   true,
			DNSSecMode: "unvalidated",
			AdminEmail: "test@" + TestZone,
			Refresh:    100,
			Retry:      100,
			Expire:     600,
			TTL:        300,
		})
		Expect(err).NotTo(HaveOccurred())

		zoneImport := zone.Import{
			ZoneData: `; Zone file for sdk-import.test. - region global
$ORIGIN sdk-import.test.
$TTL 600
@ 600 IN NS acns01.local.
@ 600 IN NS acns02.local.
@ 600 IN SOA ns0.local. admin 7 3600 1800 604800 600
test 600 IN TXT \"go-anxcloud integration test generated data\"`,
		}

		_, err = zoneAPI.Import(ctx, zoneName, zoneImport)
		Expect(err).NotTo(HaveOccurred())

		err = zoneAPI.Delete(ctx, z.Name)
		Expect(err).NotTo(HaveOccurred())
	})

	// TODO Deactivated this test cause of ENGSUP-4782
	//It("Should create update and delete a record", func() {
	//	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout*3)
	//	defer cancel()
	//
	//	zoneName := "go-sdk.test"
	//	zoneAPI := zone.NewAPI(cli)
	//
	//	z, err := zoneAPI.NewRecord(ctx, zoneName, zone.RecordRequest{
	//		Name:  "test1",
	//		RData: "test record",
	//		TTL:   300,
	//		Type:  "TXT",
	//	})
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	zoneRecords, err := zoneAPI.ListRecords(ctx, z.Name)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	var foundRecord bool
	//	var record zone.Record
	//	for _, r := range zoneRecords {
	//		if r.Name == "test1" && r.Type == "TXT" {
	//			foundRecord = true
	//			record = r
	//			break
	//		}
	//	}
	//	Expect(foundRecord).To(BeTrue())
	//	foundRecord = false
	//
	//	upd8Ctx, upd8Cancel := context.WithTimeout(context.Background(), time.Minute*5)
	//	defer upd8Cancel()
	//	z, err = zoneAPI.UpdateRecord(upd8Ctx, z.Name, record.Identifier, zone.RecordRequest{
	//		Name:  "test1-updated",
	//		Type:  "TXT",
	//		RData: "updated test record",
	//		TTL:   600,
	//	})
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	get8Ctx, get8Cancel := context.WithTimeout(context.Background(), time.Minute*5)
	//	defer get8Cancel()
	//	zoneInfo, err := zoneAPI.Get(get8Ctx, zoneName)
	//	Expect(err).NotTo(HaveOccurred())
	//	for zoneInfo.DeploymentLevel < 100 {
	//		time.Sleep(time.Second * 1)
	//		zoneInfo, err = zoneAPI.Get(get8Ctx, zoneName)
	//		Expect(err).NotTo(HaveOccurred())
	//	}
	//
	//	listCtx, listCancel := context.WithTimeout(context.Background(), time.Minute*5)
	//	defer listCancel()
	//	zoneRecords, err = zoneAPI.ListRecords(listCtx, zoneName)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	for _, r := range zoneRecords {
	//		if r.Name == "test1-updated" {
	//			foundRecord = true
	//			record = r
	//			break
	//		}
	//	}
	//	Expect(foundRecord).To(BeTrue())
	//	Expect(record).NotTo(BeNil())
	//	Expect(record.Identifier).NotTo(BeNil())
	//	Expect(record.Type).To(Equal("TXT"))
	//	Expect(record.RData).To(Equal("\"updated test record\""))
	//
	//	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), time.Minute*5)
	//	defer deleteCancel()
	//	err = zoneAPI.DeleteRecord(deleteCtx, zoneName, record.Identifier)
	//	Expect(err).NotTo(HaveOccurred())
	//})

})
