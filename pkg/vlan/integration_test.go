//go:build integration
// +build integration

package vlan

import (
	"context"
	"errors"
	"time"

	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
)

var _ = Describe("vlan client", func() {
	var cli client.Client
	var api API

	BeforeEach(func() {
		envCli, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())

		cli = envCli
		api = NewAPI(cli)
	})

	Context("with a VLAN created for testing", Ordered, func() {
		var vlanID string

		BeforeAll(func(ctx context.Context) {
			def := CreateDefinition{
				Location:            locationID,
				VMProvisioning:      false,
				CustomerDescription: "go SDK integration test",
			}
			summary, err := api.Create(ctx, def)
			Expect(err).NotTo(HaveOccurred())
			vlanID = summary.Identifier

			DeferCleanup(func(ctx context.Context) {
				err := api.Delete(ctx, vlanID)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("lists VLANs including test VLAN", func() {
			found := false
			page := 1

			for !found {
				vs, err := api.List(context.TODO(), page, 20, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(vs).NotTo(BeEmpty())

				for _, v := range vs {
					if v.Identifier == vlanID {
						found = true
						break
					}
				}

				page++
			}
		})

		It("eventually retrieves test VLAN with expected data and being active", func() {
			pollCheck := func() error {
				vlan, err := api.Get(context.TODO(), vlanID)
				if err != nil {
					return err
				}

				if len(vlan.Locations) != 1 {
					return errors.New("expected exactly one location")
				}
				if vlan.Locations[0].Identifier != locationID {
					return errors.New("location ID mismatch")
				}

				if vlan.VMProvisioning {
					return errors.New("expected VM provisioning to be false")
				}
				if vlan.CustomerDescription != "go SDK integration test" {
					return errors.New("customer description mismatch")
				}

				if vlan.Status != "Active" {
					return errors.New("VLAN not active")
				}
				return nil
			}

			Eventually(pollCheck, 5*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("updates test VLAN with changed data", func() {
			def := UpdateDefinition{
				CustomerDescription: "go SDK integration test updated",
				VMProvisioning:      true,
			}

			err := api.Update(context.TODO(), vlanID, def)
			Expect(err).NotTo(HaveOccurred())
		})

		It("eventually retrieves test VLAN with expected updated data", func() {
			pollCheck := func() error {
				vlan, err := api.Get(context.TODO(), vlanID)
				if err != nil {
					return err
				}

				if len(vlan.Locations) != 1 {
					return errors.New("expected exactly one location")
				}
				if vlan.Locations[0].Identifier != locationID {
					return errors.New("location ID mismatch")
				}

				if !vlan.VMProvisioning {
					return errors.New("expected VM provisioning to be true")
				}
				if vlan.CustomerDescription != "go SDK integration test updated" {
					return errors.New("updated customer description mismatch")
				}
				if vlan.Status != "Active" {
					return errors.New("VLAN not active")
				}
				return nil
			}

			Eventually(pollCheck, 5*time.Minute, 10*time.Second).Should(Succeed())
		})
	})
})
