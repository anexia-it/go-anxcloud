//go:build integration
// +build integration

package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
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

	matchesANX04 := func(loc Location) {
		Expect(loc.Identifier).To(Equal("52b5f6b2fd3a4a7eaaedf1a7c019e9ea"))
		Expect(loc.Name).To(Equal("AT, Vienna, Datasix"))
		Expect(loc.CityCode).To(Equal("VIE"))
		Expect(loc.CountryCode).To(Equal("AT"))
		Expect(loc.Latitude).NotTo(BeNil())
		Expect(loc.Longitude).NotTo(BeNil())
		Expect(*loc.Latitude).To(BeNumerically("~", 48.19300))
		Expect(*loc.Longitude).To(BeNumerically("~", 16.35338))
	}

	matchesANX63 := func(loc Location) {
		Expect(loc.Identifier).To(Equal("c2e3d5caa3604cbf9ae49471fb027010"))
		Expect(loc.Name).To(Equal("NL, Amsterdam"))
		Expect(loc.CityCode).To(Equal("AMS"))
		Expect(loc.CountryCode).To(Equal("NL"))
		Expect(loc.Latitude).NotTo(BeNil())
		Expect(loc.Longitude).NotTo(BeNil())
		Expect(*loc.Latitude).To(BeNumerically("~", 52.3702157))
		Expect(*loc.Longitude).To(BeNumerically("~", 4.8951679))
	}

	DescribeTable("retrieves location with expected data",
		func(identifier string, check func(l Location)) {
			l := Location{Identifier: identifier}
			err := apiClient.Get(context.TODO(), &l)
			Expect(err).NotTo(HaveOccurred())

			check(l)
		},
		Entry("ANX04", "52b5f6b2fd3a4a7eaaedf1a7c019e9ea", matchesANX04),
		Entry("ANX63", "c2e3d5caa3604cbf9ae49471fb027010", matchesANX63),
	)

	It("lists locations, retrieving example locations with expected data", func() {
		var oc types.ObjectChannel
		err := apiClient.List(context.TODO(), &Location{}, api.ObjectChannel(&oc))
		Expect(err).ToNot(HaveOccurred())

		found := make([]string, 0)

		for r := range oc {
			var loc Location
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
})
