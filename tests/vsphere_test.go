package tests_test

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/vsphere/info"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/vmlist"

	cpuperformancetype "github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/cpuperformancetypes"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/disktype"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/location"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/ips"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/progress"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/templates"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/vm"

	testUtils "github.com/anexia-it/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

const (
	templateType  = "templates"
	templateName  = "Flatcar Linux Stable"
	cpus          = 2
	sockets       = 1
	changedMemory = 4096
	memory        = 2048
	disk          = 10
)

var templateID string

// This versioning scheme that currently seems to be in place for template build numbers.
var buildNumberRegex = regexp.MustCompile(`[bB]?(\d+)`)

func extractBuildNumber(build string) int {
	matches := buildNumberRegex.FindStringSubmatch(build)
	if len(matches) != 2 {
		// panic here since someone needs to check on the regex
		panic("build does not match the buildNumberRegex")
	}

	number, err := strconv.ParseInt(matches[1], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("could not extract build for %s: %s", build, err.Error()))
	}
	return int(number)
}

func vsphereTestInit() {
	cli, err := client.New(client.AuthFromEnv(false))

	if err != nil {
		log.Fatalf("Error creating client for retrieving template ID: %v\n", err)
	}

	tplAPI := templates.NewAPI(cli)
	tpls, err := tplAPI.List(context.TODO(), locationID, templates.TemplateTypeTemplates, 1, 500)

	if err != nil {
		log.Fatalf("Error retrieving templates: %v\n", err)
	}

	selected := make([]templates.Template, 0, 1)
	for _, tpl := range tpls {
		if tpl.Name == templateName {
			selected = append(selected, tpl)
		}
	}

	sort.Slice(selected, func(i, j int) bool {
		return extractBuildNumber(selected[i].Build) > extractBuildNumber(selected[j].Build)
	})

	log.Printf("VSphere: selected template %v (build %v, ID %v)\n", selected[0].Name, selected[0].Build, selected[0].ID)

	templateID = selected[0].ID
}

var _ = Describe("vsphere API endpoint tests", func() {

	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("VMList Endpoint", func() {
		It("Should List VMs", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			vms, err := vmlist.NewAPI(cli).Get(ctx, 1, 1)
			if err != nil {
				return
			}
			Expect(vms).To(HaveLen(1))
		})
	})

	Context("provisioning endpoint", func() {

		Context("VM endpoint", func() {

			It("should create a VM and delete it later", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()

				By("reserving a new IP address")
				res, err := address.NewAPI(cli).ReserveRandom(ctx, address.ReserveRandom{
					LocationID: locationID,
					VlanID:     vlanID,
					Count:      1,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(res.Data)).To(Equal(1))

				networkInterfaces := []vm.Network{{NICType: "vmxnet3", IPs: []string{res.Data[0].Address}, VLAN: vlanID}}
				definition := vm.NewAPI(cli).NewDefinition(locationID, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
				definition.Sockets = sockets
				definition.SSH = randomPublicSSHKey()

				By("creating a new VM")
				base64Encoding := true
				provisionResponse, err := vm.NewAPI(cli).Provision(ctx, definition, base64Encoding)
				Expect(err).NotTo(HaveOccurred())

				By("waiting for the VM to be ready")
				vmID, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionResponse.Identifier)
				Expect(err).NotTo(HaveOccurred())

				By("updating the VM")
				change := vm.NewChange()
				change.MemoryMBs = changedMemory
				updateResponse, err := vm.NewAPI(cli).Update(ctx, vmID, change)
				Expect(err).NotTo(HaveOccurred())

				By("waiting for VM to be ready after an update")
				newVMid, err := progress.NewAPI(cli).AwaitCompletion(ctx, updateResponse.Identifier)
				Expect(err).NotTo(HaveOccurred())
				if newVMid != vmID {
					log.Fatalf("VM change resulted in a new ID: %v -> %v", vmID, newVMid)
				}

				By("deleting the VM")
				response, err := vm.NewAPI(cli).Deprovision(ctx, vmID, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Identifier).ToNot(BeEmpty())
				returnedIdent, err := progress.NewAPI(cli).AwaitCompletion(ctx, response.Identifier)
				Expect(err).NotTo(HaveOccurred())
				Expect(returnedIdent).To(BeEquivalentTo(vmID))
			})
		})

		Context("templates endpoint", func() {

			It("should list all available templates", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := templates.NewAPI(cli).List(ctx, locationID, templates.TemplateTypeTemplates, 1, 50)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("IPs endpoint", func() {

			It("should get a free IP address", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := ips.NewAPI(cli).GetFree(ctx, locationID, vlanID)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("disk type endpoint", func() {

			It("should list all available disk types", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()

				_, err := disktype.NewAPI(cli).List(ctx, locationID, 1, 1000)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("CPU Performance Type endpoint", func() {

			It("should list all cpu performance types", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				_, err := cpuperformancetype.NewAPI(cli).List(ctx)
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("VSphere Location endpoint", func() {

			It("should list all VSPhere locations", func() {
				ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
				defer cancel()
				locations, err := location.NewAPI(cli).List(ctx, 1, 50, "", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(locations)).To(BeNumerically(">", 0))
			})

		})

	})

	Context("info endpoint", func() {
		It("should create and retrieve a VM", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			By("reserving a new IP address")
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

			By("creating a new VM")
			base64Encoding := true
			provisionResponse, err := vm.NewAPI(cli).Provision(ctx, definition, base64Encoding)
			Expect(err).NotTo(HaveOccurred())

			By("waiting for the VM to be ready")
			vmID, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionResponse.Identifier)
			Expect(err).NotTo(HaveOccurred())

			By("retrieving the VM")
			vmInfo, err := info.NewAPI(cli).Get(ctx, vmID)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmInfo).NotTo(BeNil())
			Expect(vmInfo.Disks).To(Equal(1))
			expectedDiskSize := 10.00
			Expect(vmInfo.DiskInfo[0].DiskGB).To(Equal(expectedDiskSize))

		})
	})

	Context("progress Endpoint", func() {
		It("should handle 404 correctly", func() {
			By("using an identifiert which does not exist")
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()
			progress, err := progress.NewAPI(cli).AwaitCompletion(ctx, "this-id-does-not-exist")
			Expect(progress).To(BeEmpty())
			Expect(err).NotTo(BeNil())

		})
	})
})

var _ = Describe("build number parsing for templates", func() {
	It("extracting build number from string", func() {
		Expect(extractBuildNumber("b5555")).To(BeEquivalentTo(5555))
		Expect(extractBuildNumber("B111")).To(BeEquivalentTo(111))
		Expect(extractBuildNumber("123")).To(BeEquivalentTo(123))
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
	return fmt.Sprintf("go-test-%s", testUtils.RandomHostname())
}
