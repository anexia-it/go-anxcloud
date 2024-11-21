//go:build integration

package v1_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var errPrefixNotReadyForDelete = errors.New("Prefix not yet ready to be deleted")

func GetTestAPIClient() api.API {
	GinkgoHelper()
	SetDefaultEventuallyTimeout(2 * time.Hour) // please don't ask about my sanity with those timeouts.
	SetDefaultEventuallyPollingInterval(15 * time.Second)

	a, err := api.NewAPI(
		api.WithLogger(GinkgoLogr),
		api.WithClientOptions(client.TokenFromEnv(false)),
	)
	Expect(err).ToNot(HaveOccurred())

	return a
}

func GetTestPrefix(ipVersion ipamv1.Family, addrType ipamv1.AddressType) ipamv1.Prefix {
	GinkgoHelper()

	ctx := context.TODO()

	var netmask int = 29
	if ipVersion == ipamv1.FamilyIPv6 {
		netmask = 64
	}

	GinkgoLogr.Info("creating new test prefix")
	p := ipamv1.Prefix{
		DescriptionCustomer: fmt.Sprintf("%s - go-anxcloud IPAM E2E", testutils.TestResourceName()),
		Version:             ipVersion,
		Netmask:             netmask,
		Type:                addrType,
		Locations:           []corev1.Location{{Identifier: locationIdentifier}},
		VLANs:               []vlanv1.VLAN{{Identifier: vlanIdentifier}},
	}
	c := GetTestAPIClient()

	Expect(c.Create(ctx, &p, ipamv1.CreateEmpty(true))).To(Succeed())
	GinkgoLogr.Info("test prefix created", "identifier", p.Identifier, "url", "https://engine.anexia-it.com/service/generic-resource/detail/"+p.Identifier)
	Eventually(func(g types.Gomega) {
		GinkgoLogr.Info("waiting for status to be active", "identifier", p.Identifier, "current_status", p.Status)
		g.Expect(c.Get(ctx, &p)).To(Succeed())
		g.Expect(p.Status).To(Equal(ipamv1.StatusActive))
	}).Should(Succeed())

	DeferCleanup(func() {
		Expect(api.IgnoreNotFound(c.Destroy(ctx, &p))).To(Succeed())
		GinkgoLogr.Info("test prefix marked for deletion", "identifier", p.Identifier, "url", "https://engine.anexia-it.com/service/generic-resource/detail/"+p.Identifier)
	})

	return p
}

var vlanIdentifier = ""

var _ = SynchronizedBeforeSuite(func() []byte {
	c := GetTestAPIClient()
	vlan := vlanv1.VLAN{
		DescriptionCustomer: fmt.Sprintf("%s - go-anxcloud IPAM E2E", testutils.TestResourceName()),
		Locations: []corev1.Location{
			{Identifier: locationIdentifier},
		},
	}

	Expect(c.Create(context.Background(), &vlan)).To(Succeed())
	GinkgoLogr.Info("VLAN for tests created", "identifier", vlan.Identifier, "description", vlan.DescriptionCustomer)

	return []byte(vlan.Identifier)
}, func(identifier []byte) {
	vlanIdentifier = string(identifier)
})

// The first function would run on all parallel invocations, which is not
// necessary. The second function is only running on the first process and
// therefore the cleanup is only executed once.
var _ = SynchronizedAfterSuite(func() {}, func(ctx context.Context) {
	c := GetTestAPIClient()

	GinkgoLogr.Info("cleaning up test vlan", "identifier", vlanIdentifier)
	v := &vlanv1.VLAN{Identifier: vlanIdentifier}
	Expect(c.Get(ctx, v)).To(Succeed())

	Eventually(func(g Gomega) {
		c.Get(ctx, v)
		switch {
		case v.Status == vlanv1.StatusMarkedForDeletion:
			return
		default:
			// Cleanup of the prefixes is usually done by the individual GetTestPrefix invocations.
			// Yet, we're cleaning up all active prefixes in the VLAN, in case the cleanup got aborted or failed.
			var oc apiTypes.ObjectChannel
			g.Expect(c.List(ctx, &ipamv1.Prefix{Status: ipamv1.StatusActive}, api.ObjectChannel(&oc), api.FullObjects(true))).To(Succeed())
			for receiver := range oc {
				var p ipamv1.Prefix
				Expect(receiver(&p)).To(Succeed())

				// Unfortunately, we cannot filter for VLANs, therefore we filter other prefixes manually out.
				if len(p.VLANs) == 0 || p.VLANs[0].Identifier != v.Identifier {
					continue
				}

				switch p.Status {
				case ipamv1.StatusActive, ipamv1.StatusFailed:
					GinkgoLogr.Info("deleting test prefix", "name", p.Name, "identifier", p.Identifier)
					g.Expect(c.Destroy(ctx, &p)).To(Succeed())
				default:
					GinkgoLogr.Info("skipping deletion of test prefix, because it's not deletable", "status", p.Status, "identifier", p.Identifier)
					continue
				}
			}

			g.Expect(c.Destroy(ctx, &vlanv1.VLAN{Identifier: vlanIdentifier})).To(Succeed())
		}
	}).
		WithContext(ctx).
		Should(Succeed())
}, NodeTimeout(30*time.Minute))
