package pointer

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOptionalPointerValues(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for pkg/utils/pointer")
}

var _ = Describe("Pointer helper", func() {
	Context("convert pointer to value", func() {
		DescribeTable("nil pointer", func(actual, expected interface{}) {
			Expect(actual).To(Equal(expected))
		},
			Entry("string", StringVal(nil), ""),
			Entry("bool", BoolVal(nil), false),
			Entry("int", IntVal(nil), 0),
			Entry("uint", UIntVal(nil), uint(0)),
			Entry("float32", Float32Val(nil), float32(0)),
			Entry("float64", Float64Val(nil), float64(0)),
		)

		DescribeTable("pointer is not nil", func(actual, expected interface{}) {
			Expect(actual).To(Equal(expected))
		},
			Entry("string", StringVal(String("test")), "test"),
			Entry("bool", BoolVal(Bool(true)), true),
			Entry("int", IntVal(Int(42)), 42),
			Entry("uint", UIntVal(UInt(42)), uint(42)),
			Entry("float32", Float32Val(Float32(42)), float32(42)),
			Entry("float64", Float64Val(Float64(42)), float64(42)),
		)
	})

	type test struct {
		Bool    *bool    `json:"bool,omitempty"`
		String  *string  `json:"string,omitempty"`
		Int     *int     `json:"int,omitempty"`
		UInt    *uint    `json:"uint,omitempty"`
		Float32 *float32 `json:"float32,omitempty"`
		Float64 *float64 `json:"float64,omitempty"`
	}

	DescribeTable("test omitempty json marshal", func(source *test, expected map[string]interface{}) {
		data, err := json.Marshal(source)
		Expect(err).ToNot(HaveOccurred())
		parsed := make(map[string]interface{})
		err = json.Unmarshal(data, &parsed)
		Expect(err).ToNot(HaveOccurred())
		Expect(parsed).To(Equal(expected))
	},
		Entry(
			"non-zero values",
			&test{
				Bool:    Bool(true),
				String:  String("string"),
				Int:     Int(123),
				UInt:    UInt(123),
				Float32: Float32(1.23),
				Float64: Float64(1.23),
			},
			map[string]interface{}{
				"bool":    true,
				"string":  "string",
				"int":     float64(123),
				"uint":    float64(123),
				"float32": float64(1.23),
				"float64": float64(1.23),
			},
		),
		Entry(
			"zero values",
			&test{
				Bool:    Bool(false),
				String:  String(""),
				Int:     Int(0),
				UInt:    UInt(0),
				Float32: Float32(0),
				Float64: Float64(0),
			},
			map[string]interface{}{
				"bool":    false,
				"string":  "",
				"int":     float64(0),
				"uint":    float64(0),
				"float32": float64(0),
				"float64": float64(0),
			},
		),
		Entry(
			"nil values",
			&test{},
			map[string]interface{}{},
		),
	)
})
