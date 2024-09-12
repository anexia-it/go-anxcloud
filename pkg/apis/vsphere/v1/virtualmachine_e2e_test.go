package v1_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/apis/vsphere/v1"
)

const (
	locationIdentifier             = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	provisioningLocationIdentifier = "b164595577114876af7662092da89f76"
)

var _ = Describe("Dynamic Compute E2E tests", func() {
	var apiClient api.API

	BeforeEach(func() {
		a, err := getApi()
		Expect(err).NotTo(HaveOccurred())
		apiClient = a
	})

	When("creating a Virtual Machine", func() {
		It("completes successfully", func() {
			prepareCreate("foo", "desc")

			adc := v1.VirtualMachine{
				Name: "foo", CustomName: "desc", RAM: 1024, CPU: 1, DiskInfo: []v1.DiskInfo{
					{DiskGB: 4},
					// TODO: disk_type can be specified as well
				},
				Location:   corev1.Location{Identifier: locationIdentifier},
				TemplateID: templateIdentifier,
			}

			err := apiClient.Create(context.TODO(), &adc)
			Expect(err).NotTo(HaveOccurred())
			Expect(adc.Name).To(Equal("foo"))
			Expect(adc.CustomName).To(Equal("desc"))
			Expect(adc.Identifier).To(Equal(mockADCIdentifier))
		})
	})

	When("retrieving Virtual Machine information", func() {
		It("completes successfully", func() {
			prepareGet("foo", "bar")

			adc := v1.VirtualMachine{Identifier: mockADCIdentifier}

			err := apiClient.Get(context.TODO(), &adc)
			Expect(err).NotTo(HaveOccurred())
			Expect(adc.Name).To(Equal("foo"))
			Expect(adc.CustomName).To(Equal("bar"))
			Expect(adc.Identifier).To(Equal(mockADCIdentifier))
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
				adc := v1.VirtualMachine{}
				err := r(&adc)
				Expect(err).NotTo(HaveOccurred())

				if adc.Identifier == mockADCIdentifier {
					Expect(adc.Name).To(Equal("foo"))
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})
	})

})
