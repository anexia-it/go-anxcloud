package tests_test

import (
	"context"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vlan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VLAN API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("VLAN endpoint", func() {

		It("Should list all available VLANs", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			_, err := vlan.NewAPI(cli).List(ctx, 1, 1000, "")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should create a VLAN, update it and delete it later", func() {
			v := vlan.NewAPI(cli)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()
			summary, err := v.Create(ctx, vlan.CreateDefinition{Location: locationID, VMProvisioning: false, CustomerDescription: "go SDK integration test"})
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for vlan to be 'Active'")
			Eventually(func() string {
				info, err := v.Get(ctx, summary.Identifier)
				Expect(err).NotTo(HaveOccurred())
				return info.Status
			}, 15*time.Minute, 5*time.Second).Should(Equal("Active"))

			defer func() {
				By("Deleting the vlan")
				err = v.Delete(ctx, summary.Identifier)
				Expect(err).NotTo(HaveOccurred())
			}()

			By("Update the vlan")
			err = v.Update(ctx, summary.Identifier, vlan.UpdateDefinition{
				CustomerDescription: "go SDK integration test updated",
				VMProvisioning:      true,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Fetching the vlan again")
			Eventually(func() bool {
				vlanInfo, err := v.Get(ctx, summary.Identifier)
				Expect(err).NotTo(HaveOccurred())
				return vlanInfo.VMProvisioning
			}, 5*time.Minute, 3*time.Second).Should(BeTrue())
		})

	})

})
