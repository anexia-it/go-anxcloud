//go:build integration
// +build integration

package v1_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/client"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"

	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	errPrefixNotReadyForDelete = errors.New("Prefix not yet ready to be deleted")
)

func e2eApiClient() (api.API, error) {
	SetDefaultEventuallyTimeout(15 * time.Minute)
	SetDefaultEventuallyPollingInterval(15 * time.Second)

	return api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
}

var _ = SynchronizedBeforeSuite(func() []byte {
	c, err := e2eApiClient()
	Expect(err).NotTo(HaveOccurred())

	vlan := vlanv1.VLAN{
		DescriptionCustomer: fmt.Sprintf("%s - go-anxcloud IPAM E2E", testutils.TestResourceName()),
		Locations: []corev1.Location{
			{Identifier: locationIdentifier},
		},
	}

	err = c.Create(context.TODO(), &vlan)
	Expect(err).NotTo(HaveOccurred())

	return []byte(vlan.Identifier)
}, func(identifier []byte) {
	vlanIdentifier = string(identifier)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	c, err := e2eApiClient()
	Expect(err).NotTo(HaveOccurred())

	Eventually(func(g Gomega) {
		err := cleanupVLAN(c)
		g.Expect(err).NotTo(HaveOccurred())
	}, 30*time.Minute, 15*time.Second).Should(Succeed())
})

func cleanupVLAN(c api.API) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var oc types.ObjectChannel
	if err := c.List(ctx, &ipamv1.Prefix{}, api.ObjectChannel(&oc), api.FullObjects(true)); err != nil {
		cancel()
		return fmt.Errorf("listing prefixes failed: %w", err)
	}

	for retriever := range oc {
		var p ipamv1.Prefix
		if err := retriever(&p); err != nil {
			cancel()
			return fmt.Errorf("retrieving prefix failed: %w", err)
		}

		// we cannot filter Prefixes by VLAN (ENGSUP-5702), so we list all Prefixes and do
		// client-side filtering here.
		if len(p.VLANs) != 1 || p.VLANs[0].Identifier != vlanIdentifier {
			continue
		}

		if err := cleanupPrefix(c, p); err != nil {
			cancel()
			return fmt.Errorf("prefix cleanup failed: %w", err)
		}
	}

	// delete the VLAN
	return c.Destroy(context.TODO(), &vlanv1.VLAN{Identifier: vlanIdentifier})
}

func cleanupPrefix(c api.API, p ipamv1.Prefix) error {
	if p.Status == ipamv1.StatusMarkedForDeletion {
		return nil
	} else if p.Status != ipamv1.StatusActive {
		return fmt.Errorf("error deleting prefix %q: %w", p.Identifier, errPrefixNotReadyForDelete)
	}

	return c.Destroy(context.TODO(), &p)
}

// below are the functions for setting up the mock, empty for the prod E2E version
func preparePrefixCreate(string, *bool, ipamv1.Family, ipamv1.AddressSpace)     {}
func preparePrefixGet(string, ipamv1.Family, ipamv1.AddressSpace)               {}
func preparePrefixList(string, ipamv1.Family, ipamv1.AddressSpace)              {}
func preparePrefixUpdate(string, ipamv1.Family, ipamv1.AddressSpace)            {}
func preparePrefixDelete()                                                      {}
func preparePrefixEventuallyActive(string, ipamv1.Family, ipamv1.AddressSpace)  {}
func preparePrefixEventuallyDeleted(string, ipamv1.Family, ipamv1.AddressSpace) {}
func prepareAddressCreate(ipamv1.Prefix, string, net.IP)                        {}
func prepareAddressGet(ipamv1.Prefix, string, net.IP)                           {}
func prepareAddressList(ipamv1.Prefix, bool, string, net.IP)                    {}
func prepareAddressUpdate(ipamv1.Prefix, string, net.IP)                        {}
func prepareAddressDelete(ipamv1.Prefix, string, net.IP)                        {}
