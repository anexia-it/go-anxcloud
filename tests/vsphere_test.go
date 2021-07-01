package tests_test

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/info"
	"log"
	"math/rand"
	"strings"
	"time"

	cpuperformancetype "github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/cpuperformancetypes"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/disktype"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/location"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/ips"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/progress"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/templates"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/vm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

const (
	hostnameCharset = "abcdefghijklmnopqrstuvwxyz"
	templateType    = "templates"
	templateID      = "12c28aa7-604d-47e9-83fb-5f1d1f1837b3"
	cpus            = 2
	changedMemory   = 4096
	memory          = 2048
	disk            = 10
)

var _ = Describe("Vsphere API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Provisioning endpoint", func() {

		Context("VM endpoint", func() {

			It("Should create a VM and delete it later", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()

				By("Reserving a new IP address")
				res, err := address.NewAPI(cli).ReserveRandom(ctx, address.ReserveRandom{
					LocationID: locationID,
					VlanID:     vlanID,
					Count:      1,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(res.Data)).To(Equal(1))

				networkInterfaces := []vm.Network{{NICType: "vmxnet3", IPs: []string{res.Data[0].Address}, VLAN: vlanID}}
				definition := vm.NewAPI(cli).NewDefinition(locationID, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
				definition.SSH = randomPublicSSHKey()

				By("Creating a new VM")
				base64Encoding := true
				provisionResponse, err := vm.NewAPI(cli).Provision(ctx, definition, base64Encoding)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the VM to be ready")
				vmID, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionResponse.Identifier)
				Expect(err).NotTo(HaveOccurred())

				By("Updating the VM")
				change := vm.NewChange()
				change.MemoryMBs = changedMemory
				updateResponse, err := vm.NewAPI(cli).Update(ctx, vmID, change)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for VM to be ready after an update")
				newVMid, err := progress.NewAPI(cli).AwaitCompletion(ctx, updateResponse.Identifier)
				Expect(err).NotTo(HaveOccurred())
				if newVMid != vmID {
					log.Fatalf("VM change resulted in a new ID: %v -> %v", vmID, newVMid)
				}

				By("Deleting the VM")
				err = vm.NewAPI(cli).Deprovision(ctx, vmID, false)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Templates endpoint", func() {

			It("Should list all available templates", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := templates.NewAPI(cli).List(ctx, locationID, templates.TemplateTypeTemplates, 1, 50)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("IPs endpoint", func() {

			It("Should get a free IP address", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := ips.NewAPI(cli).GetFree(ctx, locationID, vlanID)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("Disk type endpoint", func() {

			It("Should list all available disk types", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := disktype.NewAPI(cli).List(ctx, locationID, 1, 1000)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("CPU Performance Type endpoint", func() {

			It("Should list all cpu performance types", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				_, err := cpuperformancetype.NewAPI(cli).List(ctx)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("VSphere Location endpoint", func() {

			It("Should list all VSPhere locations", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				locations, err := location.NewAPI(cli).List(ctx, 1, 50, "", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(locations)).To(BeNumerically(">", 0))
			})

		})

	})

	Context("Info endpoint", func() {
		It("Should create and retrieve a VM", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			By("Reserving a new IP address")
			res, err := address.NewAPI(cli).ReserveRandom(ctx, address.ReserveRandom{
				LocationID: locationID,
				VlanID:     vlanID,
				Count:      1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Data)).To(Equal(1))

			networkInterfaces := []vm.Network{{NICType: "vmxnet3", IPs: []string{res.Data[0].Address}, VLAN: vlanID}}
			definition := vm.NewAPI(cli).NewDefinition(locationID, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
			definition.SSH = randomPublicSSHKey()

			By("Creating a new VM")
			base64Encoding := true
			provisionResponse, err := vm.NewAPI(cli).Provision(ctx, definition, base64Encoding)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for the VM to be ready")
			vmID, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionResponse.Identifier)
			Expect(err).NotTo(HaveOccurred())

			By("Retrieving the VM")
			vmInfo, err := info.NewAPI(cli).Get(ctx, vmID)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmInfo).NotTo(BeNil())
			Expect(vmInfo.Disks).To(Equal(1))
			expectedDiskSize := 10.00
			Expect(vmInfo.DiskInfo[0].DiskGB).To(Equal(expectedDiskSize))

		})
	})

})

func randomPublicSSHKey() string {
	private, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())

	public, err := ssh.NewPublicKey(&private.PublicKey)
	Expect(err).NotTo(HaveOccurred())

	return string(ssh.MarshalAuthorizedKey(public))
}

func randomHostname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // No crypto needed here.
	hostnameSuffix := make([]string, 8)
	for i := range hostnameSuffix {
		hostnameSuffix[i] = string(hostnameCharset[r.Intn(len(hostnameCharset))])
	}

	return fmt.Sprintf("go-test-%s", strings.Join(hostnameSuffix, ""))
}
