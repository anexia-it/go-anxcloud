//go:build integration
// +build integration

package api_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api"
	clouddnsv1 "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/test"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"
)

const (
	waitTimeout  = 10 * time.Minute
	retryTimeout = 15 * time.Second
)

var _ = Describe("options", func() {
	testutils.Seed(GinkgoRandomSeed())

	var a api.API

	BeforeEach(func() {
		var err error
		a, err = api.NewAPI(
			api.WithClientOptions(client.AuthFromEnv(false)),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("AutoTag", func() {
		It("can automatically tag resources on api.Create", func() {
			vlan := vlanv1.VLAN{
				DescriptionCustomer: "go-anxcloud test api.Create AutoTag " + testutils.RandomHostname(),
				Locations: []corev1.Location{
					{Identifier: "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"},
				},
			}

			DeferCleanup(func() {
				if vlan.Identifier == "" {
					return
				}

				pollCheck := func() error {
					err := a.Get(context.TODO(), &vlan)
					if err != nil {
						return err
					}
					if vlan.Status != vlanv1.StatusActive {
						return errors.New("VLAN not yet active")
					}
					return nil
				}

				Eventually(pollCheck, waitTimeout, retryTimeout).Should(Succeed())

				err := a.Destroy(context.TODO(), &vlan)
				Expect(err).ToNot(HaveOccurred())
			})

			err := a.Create(context.TODO(), &vlan, api.AutoTag("foo", "bar", "baz"))
			Expect(err).ToNot(HaveOccurred())

			tags, err := corev1.ListTags(context.TODO(), a, &vlan)
			Expect(err).ToNot(HaveOccurred())
			Expect(tags).To(ContainElements("foo", "bar", "baz"))
		})

		It("returns an error when resource tagging failed", func() {
			zone := clouddnsv1.Zone{
				Name:       test.RandomHostname() + ".go-anxcloud.test",
				IsMaster:   true,
				DNSSecMode: "unvalidated",
				AdminEmail: "admin@go-anxcloud.test",
				Refresh:    14400,
				Retry:      3600,
				Expire:     604800,
				TTL:        3600,
			}

			DeferCleanup(func() {
				_ = a.Destroy(context.TODO(), &zone)
			})

			// Note: AutoTag fails, because CloudDNS zones cannot be tagged with the Name identifier
			// See ENGSUP-5900
			err := a.Create(context.TODO(), &zone, api.AutoTag("foo", "bar", "baz"))
			Expect(err).To(HaveOccurred())
			var e *api.ErrTaggingFailed
			Expect(errors.As(err, &e)).To(BeTrue())
		})
	})
})
