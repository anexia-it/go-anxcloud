//go:build integration
// +build integration

package v1_test

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("location E2E tests", func() {
	var apiClient api.API

	BeforeEach(func() {
		a, err := api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
		Expect(err).ToNot(HaveOccurred())
		apiClient = a
	})

	matchesANX04 := func(loc corev1.Location) {
		Expect(loc.Identifier).To(Equal("52b5f6b2fd3a4a7eaaedf1a7c019e9ea"))
		Expect(loc.Name).To(Equal("AT, Vienna, Datasix"))
		Expect(loc.CityCode).To(Equal("VIE"))
		Expect(loc.CountryCode).To(Equal("AT"))
		Expect(loc.Latitude).NotTo(BeNil())
		Expect(loc.Longitude).NotTo(BeNil())
		Expect(*loc.Latitude).To(Equal("48.1930000000000"))
		Expect(*loc.Longitude).To(Equal("16.3533800000000"))
	}

	matchesANX63 := func(loc corev1.Location) {
		Expect(loc.Identifier).To(Equal("c2e3d5caa3604cbf9ae49471fb027010"))
		Expect(loc.Name).To(Equal("NL, Amsterdam"))
		Expect(loc.CityCode).To(Equal("AMS"))
		Expect(loc.CountryCode).To(Equal("NL"))
		Expect(loc.Latitude).NotTo(BeNil())
		Expect(loc.Longitude).NotTo(BeNil())
		Expect(*loc.Latitude).To(Equal("52.3702157000000"))
		Expect(*loc.Longitude).To(Equal("4.8951679000000"))
	}

	DescribeTable("retrieves location with expected data",
		func(identifier string, check func(l corev1.Location)) {
			l := corev1.Location{Identifier: identifier}
			err := apiClient.Get(context.TODO(), &l)
			Expect(err).NotTo(HaveOccurred())

			check(l)
		},
		Entry("ANX04", "52b5f6b2fd3a4a7eaaedf1a7c019e9ea", matchesANX04),
		Entry("ANX63", "c2e3d5caa3604cbf9ae49471fb027010", matchesANX63),
	)

	It("lists locations, retrieving example locations with expected data", func() {
		var oc types.ObjectChannel
		err := apiClient.List(context.TODO(), &corev1.Location{}, api.ObjectChannel(&oc))
		Expect(err).ToNot(HaveOccurred())

		found := make([]string, 0)

		for r := range oc {
			var loc corev1.Location
			err := r(&loc)
			Expect(err).NotTo(HaveOccurred())

			found = append(found, loc.Code)

			if loc.Code == "ANX04" {
				matchesANX04(loc)
			} else if loc.Code == "ANX63" {
				matchesANX63(loc)
			}
		}

		Expect(found).To(ContainElement("ANX04"))
		Expect(found).To(ContainElement("ANX63"))
	})

	DescribeTable("Gets locations by code", func(code string, check func(l corev1.Location)) {
		loc := corev1.Location{Code: code}
		err := apiClient.Get(context.TODO(), &loc)
		Expect(err).ToNot(HaveOccurred())
		check(loc)
	},
		Entry("ANX04", "ANX04", matchesANX04),
		Entry("ANX63", "ANX63", matchesANX63),
	)
})
