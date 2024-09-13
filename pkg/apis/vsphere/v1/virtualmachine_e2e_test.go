package v1_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/apis/vsphere/v1"
	testutil "go.anx.io/go-anxcloud/pkg/utils/test"
)

const (
	waitTimeout  = 30 * time.Millisecond
	retryTimeout = 10 * time.Millisecond
)

var _ = Describe("Dynamic Compute E2E tests", func() {
	var apiClient api.API
	var name string
	var desc string
	var identifier string
	var template v1.Template

	locationIdentifier := mockLocationIdentifier
	vlanIdentifier := mockVLANIdentifier
	templateIdentifier := mockTemplateIdentifier
	sshKey := mockSSHKey
	vlanIPAddress := mockIPAddress

	if isIntegrationTest {
		if os.Getenv("ANEXIA_LOCATION_ID") != "" {
			locationIdentifier = os.Getenv("ANEXIA_LOCATION_ID")
		}
		if os.Getenv("ANEXIA_VLAN_ID") != "" {
			vlanIdentifier = os.Getenv("ANEXIA_VLAN_ID")
		}
		if os.Getenv("ANEXIA_VLAN_IP_ADDRESS") != "" {
			vlanIPAddress = os.Getenv("ANEXIA_VLAN_IP_ADDRESS")
		}
		if os.Getenv("ANEXIA_TEMPLATE_ID") != "" {
			templateIdentifier = os.Getenv("ANEXIA_TEMPLATE_ID")
		}
		if os.Getenv("ANEXIA_SSH_KEY") != "" {
			sshKey = os.Getenv("ANEXIA_SSH_KEY")
		}
	}

	BeforeEach(func() {
		a, err := getApi()
		Expect(err).NotTo(HaveOccurred())
		apiClient = a
	})

	Context("creating a Virtual Machine from template", Ordered, func() {
		BeforeAll(func() {
			name = "go-anxcloud-test-" + testutil.RandomHostname()
			desc = "go-anxcloud test"
			prepareCreate(name, desc)

			template = v1.Template{Identifier: templateIdentifier, Type: v1.TypeTemplate, Location: corev1.Location{Identifier: locationIdentifier}}
			if isIntegrationTest {
				err := apiClient.Get(context.TODO(), &template)
				Expect(err).NotTo(HaveOccurred())
				Expect(template.Identifier).To(Equal(templateIdentifier))
			}

			vm := v1.VirtualMachine{
				Cores:              1,
				CPUPerformanceType: "performance-amd",
				CustomName:         desc,
				DiskInfo: []v1.DiskInfo{
					{DiskGB: 10, DiskType: "ENT2"},
					{DiskGB: 10, DiskType: "ENT2"},
				},
				Location: corev1.Location{Identifier: locationIdentifier},
				Name:     name,
				Networks: []v1.Network{
					{
						NICType: "vmxnet3", VLAN: vlanIdentifier,
						BandwidthLimit: v1.Bandwidth1GBit, IPs: []string{vlanIPAddress},
					},
				},
				RAM:         1024,
				SSHKey:      sshKey,
				StartScript: "#/bin/sh\n",
				TemplateID:  template.Identifier,
			}

			err := apiClient.Create(context.TODO(), &vm)
			if err != nil {
				fmt.Printf("Failed to create virtual machine: %s\n", err)
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(vm.Name).To(Equal(name))
			Expect(vm.CustomName).To(Equal(desc))

			identifier = vm.Identifier
			if !isIntegrationTest {
				Expect(identifier).To(Equal(mockVMIdentifier))
			}

			DeferCleanup(func() {
				prepareDelete()
				err := apiClient.Destroy(context.TODO(), &v1.VirtualMachine{Identifier: identifier})
				if err != nil {
					Fail(fmt.Sprintf("Error destroying Virtual Machine %q created for testing: %v", identifier, err))
				}
			})
		})

		It("retrieves Virtual Machine information", func() {
			prepareGetInfo(name, desc)

			vm := v1.VirtualMachine{Identifier: identifier}
			err := apiClient.Get(context.TODO(), &vm)

			Expect(err).NotTo(HaveOccurred())
			Expect(vm.Name).To(Equal(name))
			Expect(vm.CustomName).To(Equal(desc))
			Expect(vm.Identifier).To(Equal(identifier))
		})

		//It("retrieves Virtual Machine status", func() {})

		It("eventually is poweredOn", func() {
			waitTimeout := waitTimeout
			retryTimeout := retryTimeout
			if isIntegrationTest {
				// ms to seconds for integration test
				waitTimeout *= 1000
				retryTimeout *= 1000
			}
			Eventually(func(g Gomega) {
				prepareEventuallyActive(name, desc)

				vm := v1.VirtualMachine{Identifier: identifier}
				err := apiClient.Get(context.TODO(), &vm)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(vm.Status).To(Equal(v1.StatusPoweredOn))
				g.Expect(vm.DiskInfo).To(HaveLen(vm.Disks))
				g.Expect(vm.DiskInfo[0].DiskType).To(Equal("ENT2"))
				g.Expect(vm.Networks).To(HaveLen(1))
				g.Expect(vm.Networks[0].IPsv4).To(HaveLen(1))
			}, waitTimeout, retryTimeout).Should(Succeed())
		})
	})

	When("retrieving a list of Virtual Machines", func() {
		It("completes successfully", func() {
			prepareList("foo", "bar")

			var oc types.ObjectChannel
			err := apiClient.List(context.TODO(), &v1.VirtualMachine{}, api.ObjectChannel(&oc))
			Expect(err).NotTo(HaveOccurred())

			found := false
			for r := range oc {
				vm := v1.VirtualMachine{}
				err := r(&vm)
				Expect(err).NotTo(HaveOccurred())

				if vm.Identifier == mockVMIdentifier {
					Expect(vm.Name).To(Equal("foo"))
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})
	})

})
