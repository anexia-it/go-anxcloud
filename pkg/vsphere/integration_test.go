//go:build integration
// +build integration

package vsphere

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"go.anx.io/go-anxcloud/pkg/vsphere/info"
	"go.anx.io/go-anxcloud/pkg/vsphere/vmlist"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/ipam/address"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/progress"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/templates"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/vm"

	"go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gomegaTypes "github.com/onsi/gomega/types"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "16854ecb42af4fad89f9fcef26789d50"

	templateType       = "templates"
	templateNamePrefix = "Flatcar Linux"
	cpus               = 2
	sockets            = 1
	changedMemory      = 4096
	memory             = 2048
	disk               = 10
)

func BeBuiltFromTemplate(id string) gomegaTypes.GomegaMatcher {
	return WithTransform(
		func(i info.Info) string { return i.TemplateID },
		Equal(id),
	)
}

func HaveCPUs(cpus int) gomegaTypes.GomegaMatcher {
	return WithTransform(
		func(i info.Info) int { return i.Cores },
		Equal(cpus),
	)
}

func HaveSockets(sockets int) gomegaTypes.GomegaMatcher {
	return WithTransform(
		func(i info.Info) int { return i.Cores / i.CPU },
		Equal(sockets),
	)
}

func HaveMemory(memory int) gomegaTypes.GomegaMatcher {
	return WithTransform(
		func(i info.Info) int { return i.RAM },
		Equal(memory),
	)
}

func HaveDisks(diskSizes ...int) gomegaTypes.GomegaMatcher {
	floatSizes := make([]float64, 0, len(diskSizes))
	for _, s := range diskSizes {
		floatSizes = append(floatSizes, float64(s))
	}

	return SatisfyAll(
		WithTransform(
			func(i info.Info) int { return i.Disks },
			Equal(len(diskSizes)),
		),
		WithTransform(
			func(i info.Info) []float64 {
				ret := make([]float64, 0, len(i.DiskInfo))
				for _, di := range i.DiskInfo {
					ret = append(ret, di.DiskGB)
				}
				return ret
			},
			BeEquivalentTo(floatSizes),
		),
	)
}

func HaveIPv4Addresses(addresses ...[]string) gomegaTypes.GomegaMatcher {
	return WithTransform(
		func(i info.Info) [][]string {
			ret := make([][]string, 0, len(i.Network))
			for _, net := range i.Network {
				ret = append(ret, net.IPv4)
			}
			return ret
		},
		SatisfyAll(
			HaveLen(len(addresses)),
			ConsistOf(addresses),
		),
	)
}

// Best practice is to have tests independent from each other, but this would require us to spawn a new VM and
// wait for it to be ready for most of the tests in this block. Because this increases test runtime a lot, we
// opted to a ordered aproach, having the VM-create test create the VM we use for the other tests. This makes
// the tests depending on each other, but reduces runtime _a lot_.
// When adding new test cases, be especially careful to check if everything is fine after you modified the VM
// to ensure test cases following your changed one don't fail on things you didn't check and they didn't touch.
var _ = Describe("vsphere API client", Ordered, func() {
	var cli client.Client

	BeforeEach(func() {
		var err error
		cli, err = client.New(
			client.AuthFromEnv(false),
			client.LogWriter(GinkgoWriter),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	var ipAddress string
	var templateID string
	var vmID string
	var provisionID string

	verifyVMInfo := func(expectedMemory int) {
		It("eventually retrieves the test VM with expected data", func() {
			getVMInfo := func(g Gomega) info.Info {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				vmInfo, err := info.NewAPI(cli).Get(ctx, vmID)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(vmInfo).NotTo(BeNil())

				return vmInfo
			}

			Eventually(getVMInfo, 5*time.Minute, 30*time.Second).Should(SatisfyAll(
				BeBuiltFromTemplate(templateID),
				HaveCPUs(cpus),
				HaveSockets(sockets),
				HaveMemory(expectedMemory),
				HaveDisks(disk),
				HaveIPv4Addresses([]string{ipAddress}),
			))
		})
	}

	It("should find the template by name", func() {
		tplAPI := templates.NewAPI(cli)
		tpls, err := tplAPI.List(context.TODO(), locationID, templates.TemplateTypeTemplates, 1, 500)

		Expect(err).NotTo(HaveOccurred())

		selected := make([]templates.Template, 0, 1)
		for _, tpl := range tpls {
			if strings.HasPrefix(tpl.Name, templateNamePrefix) {
				selected = append(selected, tpl)
			}
		}

		sort.Slice(selected, func(i, j int) bool {
			return extractBuildNumber(selected[i].Build) > extractBuildNumber(selected[j].Build)
		})

		templateID = selected[0].ID
	})

	It("should reserve an IP address for our test VM", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := address.NewAPI(cli).ReserveRandom(ctx, address.ReserveRandom{
			LocationID: locationID,
			VlanID:     vlanID,
			Count:      1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Data).To(HaveLen(1))

		ipAddress = res.Data[0].Address
	})

	It("should create the test VM with the reserved IP address", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		definition := vm.NewAPI(cli).NewDefinition(
			locationID,
			templateType, templateID,
			randomHostname(),
			cpus, memory, disk,
			[]vm.Network{
				{
					NICType: "vmxnet3",
					IPs:     []string{ipAddress},
					VLAN:    vlanID,
				},
			},
		)
		definition.Sockets = sockets
		definition.SSH = randomPublicSSHKey()

		provisionResponse, err := vm.NewAPI(cli).Provision(ctx, definition, true)
		Expect(err).NotTo(HaveOccurred())

		provisionID = provisionResponse.Identifier
	})

	It("eventually retrieves the test VM provisioning being completed", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		id, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionID)
		Expect(err).NotTo(HaveOccurred())

		vmID = id
	})

	It("lists VMs including our test VM", func() {
		found := false
		page := 1

		for !found {
			vms, err := vmlist.NewAPI(cli).Get(context.TODO(), page, 20)
			Expect(err).NotTo(HaveOccurred())
			Expect(vms).NotTo(BeEmpty())

			for _, vm := range vms {
				if vm.Identifier == vmID {
					found = true
					break
				}
			}

			page++
		}
	})

	verifyVMInfo(memory)

	It("should update our test VM", func() {
		change := vm.NewChange()
		change.MemoryMBs = changedMemory
		updateResponse, err := vm.NewAPI(cli).Update(context.TODO(), vmID, change)
		Expect(err).NotTo(HaveOccurred())

		provisionID = updateResponse.Identifier
	})

	It("eventually retrieves the test VM updating being completed", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		id, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionID)
		Expect(err).NotTo(HaveOccurred())
		Expect(id).To(Equal(vmID))
	})

	verifyVMInfo(changedMemory)

	It("deletes the VM, waiting for it to be gone", func() {
		response, err := vm.NewAPI(cli).Deprovision(context.TODO(), vmID, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Identifier).ToNot(BeEmpty())

		provisionID = response.Identifier
	})

	It("eventually retrieves the test VM deletion being completed", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		returnedIdent, err := progress.NewAPI(cli).AwaitCompletion(ctx, provisionID)
		Expect(err).NotTo(HaveOccurred())
		Expect(returnedIdent).To(BeEquivalentTo(vmID))

		vmID = "" // we clear vmID to signal the VM already being deleted to the AfterAll block
	})

	// this block deletes the created VM if not already done by the test for that
	AfterAll(func() {
		if vmID == "" {
			return
		}

		By("deleting the VM")
		response, err := vm.NewAPI(cli).Deprovision(context.TODO(), vmID, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Identifier).ToNot(BeEmpty())
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
	return fmt.Sprintf("go-test-%s", test.RandomHostname())
}
