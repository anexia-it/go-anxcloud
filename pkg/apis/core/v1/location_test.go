package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var _ = Describe("Location Object", func() {
	var o *Location

	BeforeEach(func() {
		o = &Location{}
	})

	DescribeTable("gives ErrOperationNotSupported",
		func(op types.Operation) {
			_, err := o.EndpointURL(types.ContextWithOperation(context.TODO(), op))
			Expect(err).To(MatchError(api.ErrOperationNotSupported))
		},
		Entry("for Create operation", types.OperationCreate),
		Entry("for Update operation", types.OperationUpdate),
		Entry("for Destroy operation", types.OperationDestroy),
	)

	DescribeTable("ResponseDecodeHook",
		// expLat and expLon are passed as string to have nil ("") and a value without needing pointers,
		// constants of which are a pain in the ass to pass to a function
		func(expErr error, expLatS, expLonS, bodySnippet string) {
			var expLat *float64
			var expLon *float64

			if expLatS != "" {
				v, err := strconv.ParseFloat(expLatS, 64)
				if err != nil {
					panic(err)
				}

				expLat = &v
			}

			if expLonS != "" {
				v, err := strconv.ParseFloat(expLonS, 64)
				if err != nil {
					panic(err)
				}

				expLon = &v
			}

			loc := Location{}

			data := []byte(
				fmt.Sprintf(
					`{`+
						`"identifier":"foo",`+
						`"name":"DE, Chemnitz, dev home office",`+
						`"code":"ANX1337",`+
						`"city_code":"C",`+
						`"country":"DE",`+
						`%s}`,
					bodySnippet,
				),
			)
			err := loc.UnmarshalJSON(data)

			if expErr != nil {
				Expect(err).To(Or(
					MatchError(expErr),
					BeAssignableToTypeOf(expErr),
				))
			} else {
				Expect(err).NotTo(HaveOccurred())

				if expLat != nil {
					Expect(loc.Latitude).NotTo(BeNil())
					Expect(*loc.Latitude).To(BeNumerically("~", *expLat))
				} else {
					Expect(loc.Latitude).To(BeNil())
				}

				if expLon != nil {
					Expect(loc.Longitude).NotTo(BeNil())
					Expect(*loc.Longitude).To(BeNumerically("~", *expLon))
				} else {
					Expect(loc.Longitude).To(BeNil())
				}
			}
		},
		Entry("invalid json given",
			&json.SyntaxError{}, "", "",
			`asäldjaölskdjöalsnbf SDKGJNmäqRE`,
		),
		Entry("given values",
			nil, "42", "23",
			`"lat":"42.0000000","lon":"23.00000"`,
		),
		Entry("only latitude given",
			nil, "42", "",
			`"lat":"42.0000000","lon":null`,
		),
		Entry("only longitude given",
			nil, "", "23",
			`"lat":null,"lon":"23.00000"`,
		),
		Entry("none given",
			nil, "", "",
			`"lat":null,"lon":null`,
		),
		Entry("empty strings given",
			strconv.ErrSyntax, "", "",
			`"lat":"","lon":""`,
		),
		Entry("invalid string given for latitude",
			strconv.ErrSyntax, "", "",
			`"lat":"blurb","lon":null`,
		),
		Entry("invalid string given for longitude",
			strconv.ErrSyntax, "", "",
			`"lat":null,"lon":"blurb"`,
		),
	)
})
