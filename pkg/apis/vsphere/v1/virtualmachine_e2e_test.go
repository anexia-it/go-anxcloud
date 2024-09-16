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
	defaultWaitTimeout  = 30 * time.Millisecond
	defaultRetryTimeout = 10 * time.Millisecond
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
		// TODO: replace this with an API call
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

			// get the template
			// TODO use/add template mocks?
			template = v1.Template{Identifier: templateIdentifier, Type: v1.TypeTemplate, Location: corev1.Location{Identifier: locationIdentifier}}
			if isIntegrationTest {
				err := apiClient.Get(context.TODO(), &template)
				Expect(err).NotTo(HaveOccurred())
				Expect(template.Identifier).To(Equal(templateIdentifier))
			}

			//// get a free IP
			//waitTimeout := defaultWaitTimeout
			//retryTimeout := defaultRetryTimeout
			//if isIntegrationTest {
			//	// 30 ms => 30 seconds
			//	waitTimeout *= 1000
			//	retryTimeout *= 1000
			//}
			//
			//prepareGetIPs()
			//var channel types.ObjectChannel
			//ips := v1.IPs{} // LocationIdentifier: locationIdentifier, VLANIdentifier: vlanIdentifier}
			//err := apiClient.List(context.TODO(), &ips, api.ObjectChannel(&channel))
			//Expect(err).NotTo(HaveOccurred())
			//
			//Eventually(func(g Gomega) {
			//	prepareGetIPs()
			//
			//	//fmt.Printf("%+v\n", ips)
			//
			//	i := v1.IPs{}
			//	for ip := range channel {
			//		err = ip(&i)
			//		g.Expect(err).NotTo(HaveOccurred())
			//      // TODO: why is this always "empty"?
			//		fmt.Printf("%+v\n", i)
			//		g.Expect(i.Data).ToNot(BeEmpty())
			//		//g.Expect(i.Data[0].Text).To(Equal(mockFreeIPAddress))
			//		//vlanIPAddress = i.Data[0].Text
			//		break
			//	}
			//
			//	//g.Expect(ips.Data[0].Text).To(Equal(mockFreeIPAddress))
			//
			//}, waitTiemout, retryTimeout).Should(Succeed())

			prepareCreate(name, desc)
			vm := v1.VirtualMachine{
				Cores:              2,
				CPUPerformanceType: "performance-amd",
				CustomName:         desc,
				DiskInfo: []v1.DiskInfo{
					{DiskGB: 10, DiskType: "ENT6"},
					{DiskGB: 10, DiskType: "STD4"},
				},
				Location: corev1.Location{Identifier: locationIdentifier},
				Name:     name,
				Networks: []v1.Network{
					{
						NICType: "vmxnet3", VLAN: vlanIdentifier,
						BandwidthLimit: v1.Bandwidth1GBit, IPs: []string{vlanIPAddress},
					},
				},
				RAM:         2048,
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

			Expect(vm.Progress.Identifier).ToNot(BeEmpty())
			Expect(vm.Progress.Errors).To(BeEmpty())

			fmt.Printf("Creating virtual machine, progress identifier: %s\n", vm.Progress.Identifier)

			// After creation, we get the VM identifier from the progress endpoint
			waitTimeout := defaultWaitTimeout
			retryTimeout := defaultRetryTimeout
			if isIntegrationTest {
				// 30 ms => 30*60 seconds = 30 min
				waitTimeout *= 60000
				retryTimeout *= 60000
			}
			Eventually(func(g Gomega) {
				prepareEventuallyProvisioned()

				pp := v1.ProvisionProgress{Identifier: vm.Progress.Identifier}
				err := apiClient.Get(context.TODO(), &pp)

				fmt.Printf("Provisioning progress: %d\n", pp.Percent)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(pp.Errors).To(BeEmpty())
				g.Expect(pp.Percent).To(Equal(100))
				g.Expect(pp.Status).To(Equal("1"))
				g.Expect(pp.VMIdentifier).ToNot(BeEmpty())

				identifier = pp.VMIdentifier
			}, waitTimeout, retryTimeout).Should(Succeed())

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
			if isIntegrationTest {
				Expect(vm.Name).To(MatchRegexp("[0-9]*-" + name))
			} else {
				Expect(vm.Name).To(Equal(name))
			}
			Expect(vm.CustomName).To(Equal(desc))
			Expect(vm.Identifier).To(Equal(identifier))
			Expect(vm.Location.Identifier).To(Equal(locationIdentifier))
			Expect(vm.Disks).To(Equal(2))
			Expect(vm.DiskInfo).To(HaveLen(2))
			Expect(vm.Networks).To(HaveLen(1))
		})

		//It("retrieves Virtual Machine status", func() {})

		It("eventually is poweredOn", func() {
			waitTimeout := defaultWaitTimeout
			retryTimeout := defaultRetryTimeout
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
			}, waitTimeout, retryTimeout).Should(Succeed())
		})
	})

	When("retrieving a list of Virtual Machines", func() {
		// TODO: generic coverage for integration test
		if isIntegrationTest {
			return
		}

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
