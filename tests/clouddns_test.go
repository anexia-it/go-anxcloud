package tests_test

import (
	"context"
	"encoding/json"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"net/http"
	"time"
)

var _ = Describe("CloudDNS API endpoint tests", func() {
	var cli client.Client

	const TestZone string = "xocp.de"

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Definition List Endpoint", func() {
		It("Should list all available zones", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).List(ctx)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Definition Get Endpoint", func() {
		It("Should return the zone", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := zone.NewAPI(cli).Get(ctx, TestZone)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("Definition Create Endpoint", func() {
		var createTestZoneName = "sdk-create-test.xocp.de"

		It("Should create the zone", func() {
			randRefresh := rand.Intn(10) * 100
			randRetry := rand.Intn(10) * 100
			randExpire := rand.Intn(10) * 1000
			randTTL := rand.Intn(10) * 100

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Response{
					Definition: &zone.Definition{
						Name:       createTestZoneName,
						ZoneName:   createTestZoneName,
						IsMaster:   true,
						DNSSecMode: "unvalidated",
						AdminEmail: "admin@xocp.de",
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
				AdminEmail: "admin@xocp.de",
				Refresh:    300,
				Retry:      300,
				Expire:     3600,
				TTL:        300,
			}
			response, err := zone.NewAPI(c).Create(ctx, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.Name).To(Equal(createTestZoneName))
			Expect(response.AdminEmail).To(Equal("admin@xocp.de"))
		})
	})

	Context("Definition Update Endpoint", func() {
		updateTestZoneName := "sdk-update-test.xocp.de"

		It("Should update the zone", func() {
			randRefresh := rand.Intn(10) * 100
			randRetry := rand.Intn(10) * 100
			randExpire := rand.Intn(10) * 1000
			randTTL := rand.Intn(10) * 100

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Response{
					Definition: &zone.Definition{
						Name:       updateTestZoneName,
						ZoneName:   updateTestZoneName,
						IsMaster:   true,
						DNSSecMode: "unvalidated",
						AdminEmail: "test@xocp.de",
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
				AdminEmail: "test@xocp.de",
				Refresh:    randRefresh,
				Retry:      randRetry,
				Expire:     randExpire,
				TTL:        randTTL,
			}
			response, err := zone.NewAPI(c).Update(ctx, updateTestZoneName, createDefinition)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(Not(BeNil()))
			Expect(response).To(Not(BeNil()))
			Expect(response.AdminEmail).To(Equal("test@xocp.de"))
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
		changesetZoneName := "sdk-changeset-test.xocp.de"

		It("Should apply the changeset", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.ChangeSet
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := []zone.Record{{
					Identifier: "",
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
		importZoneName := "sdk-import-test.xocp.de"
		It("Should import the zone", func() {

			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Import
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				resp := zone.Revision{
					CreatedAt:  time.Now(),
					Identifier: "some-identifier",
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
		recordZoneName := "sdk-record-test.xocp.de"

		It("Should create the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := zone.Response{
					Definition: &zone.Definition{
						Name:       recordZoneName,
						ZoneName:   recordZoneName,
						IsMaster:   true,
					},
					Revisions: []zone.Revision{{
						Identifier: "test-uuid",
						Records:    []zone.Record{
							{
								Identifier: "record-identifier",
								Immutable:  false,
								Name:       "test1",
								RData:      "test record",
								Region:     "default",
								TTL:        &ttl,
								Type:       "TXT",
							},
						},
						Serial:     0,
						State:      "active",
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
		recordZoneName := "sdk-record-test.xocp.de"

		It("Should update the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request zone.Definition
				err := json.NewDecoder(r.Body).Decode(&request)
				Expect(err).NotTo(HaveOccurred())
				ttl := 300
				resp := zone.Response{
					Definition: &zone.Definition{
						Name:       recordZoneName,
						ZoneName:   recordZoneName,
						IsMaster:   true,
					},
					Revisions: []zone.Revision{{
						Identifier: "test-uuid",
						Records:    []zone.Record{
							{
								Identifier: "record-identifier",
								Immutable:  false,
								Name:       "test1",
								RData:      "test record",
								Region:     "default",
								TTL:        &ttl,
								Type:       "TXT",
							},
						},
						Serial:     0,
						State:      "active",
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

			response, err := zone.NewAPI(c).UpdateRecord(ctx, recordZoneName, "some=test-record-id", record)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
			Expect(response.Revisions).To(ContainElements())
		})
	})

	Context("Definition Delete Record Endpoint", func() {
		recordZoneName := "sdk-record-test.xocp.de"

		It("Should delete the record", func() {
			c, server := client.NewTestClient(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal(http.MethodDelete))
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			recordIdentifier := "some-test-record-id"
			ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
			defer cancel()
			err := zone.NewAPI(c).DeleteRecord(ctx, recordZoneName, recordIdentifier)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
