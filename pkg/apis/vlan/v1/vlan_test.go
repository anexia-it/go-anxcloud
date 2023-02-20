package v1

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	testutil "go.anx.io/go-anxcloud/pkg/utils/test"
)

var _ = Describe("VLAN Object", func() {
	DescribeTable("EndpointURL",
		func(status Status, locations []corev1.Location, expectedQuery url.Values, expectedErrors ...error) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationList)

			o := VLAN{
				Status:    status,
				Locations: locations,
			}

			url, err := o.EndpointURL(ctx)

			if len(expectedErrors) > 0 {
				Expect(err).To(HaveOccurred())

				for _, e := range expectedErrors {
					Expect(err).To(MatchError(e))
				}
			} else {
				query := url.Query()
				Expect(query).To(BeEquivalentTo(expectedQuery))
			}
		},
		Entry(
			"when no filters are set",
			StatusInvalid,
			nil,
			url.Values{},
		),
		Entry(
			"when Status is set to Active",
			StatusActive,
			nil,
			url.Values{
				"status": []string{"Active"},
			},
		),
		Entry(
			"when one Location is set",
			StatusInvalid,
			[]corev1.Location{{Identifier: "foo"}},
			url.Values{
				"location": []string{"foo"},
			},
		),
		Entry(
			"when Status and one Location is set",
			StatusActive,
			[]corev1.Location{
				{Identifier: "foo"},
			},
			url.Values{
				"status":   []string{"Active"},
				"location": []string{"foo"},
			},
		),
		Entry(
			"when two Location are set",
			StatusInvalid,
			[]corev1.Location{
				{Identifier: "foo"},
				{Identifier: "bar"},
			},
			nil,
			types.ErrInvalidFilter,
			ErrFilterMultipleLocations,
		),
	)

	DescribeTable("FilterAPIRequest",
		func(locations []corev1.Location, expErr error) {
			ctx := types.ContextWithOperation(context.TODO(), types.OperationCreate)

			vlan := VLAN{Locations: locations}

			data, err := vlan.FilterAPIRequestBody(ctx)

			if expErr != nil {
				Expect(err).To(MatchError(expErr))
				Expect(data).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(data).NotTo(BeNil())

				rval := reflect.ValueOf(data)
				locationIdentifier := rval.FieldByName("Location").Interface().(string)

				Expect(locationIdentifier).To(Equal(locations[0].Identifier))
			}
		},
		Entry(
			"errors without any location",
			[]corev1.Location{},
			ErrLocationCount,
		),
		Entry(
			"errors with two locations",
			[]corev1.Location{
				{Identifier: "foo"},
				{Identifier: "bar"},
			},
			ErrLocationCount,
		),
		Entry(
			"succeeds with a single location",
			[]corev1.Location{
				{Identifier: "foo"},
			},
			nil,
		),
	)
})

func TestVLAN(t *testing.T) {
	testutil.Seed(GinkgoRandomSeed())

	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for VLAN API definition")
}
