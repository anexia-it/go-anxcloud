package compare

import (
	"fmt"

	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleSearch() {
	lb := lbaasv1.LoadBalancer{Identifier: "some identifier string"}
	existingObjects := []lbaasv1.Frontend{
		{
			Name:         "Frontend A",
			Identifier:   "some identifier string A",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
		{
			Name:         "Frontend B",
			Identifier:   "some identifier string B",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
	}

	// positive example
	index, err := Search(lbaasv1.Frontend{
		Name:         "Frontend A",
		Mode:         lbaasv1.TCP,
		LoadBalancer: &lb,
	}, existingObjects, "Name", "Mode", "LoadBalancer.Identifier")
	if err != nil {
		fmt.Printf("Error comparing the objects: %v\n", err)
	} else {
		if index == -1 {
			fmt.Printf("Did not find our Object\n")
		} else {
			fmt.Printf("Index of our Object is %v\n", index)
		}
	}

	// negative example
	index, err = Search(lbaasv1.Frontend{
		Name: "Frontend C",
		Mode: lbaasv1.HTTP,
	}, existingObjects, "Name", "Mode", "LoadBalancer.Identifier")
	if err != nil {
		fmt.Printf("Error comparing the objects: %v\n", err)
	} else {
		if index == -1 {
			fmt.Printf("Did not find our Object\n")
		} else {
			fmt.Printf("Index of our Object is %v\n", index)
		}
	}

	// Output:
	// Index of our Object is 0
	// Did not find our Object
}

var _ = Describe("Search", func() {
	DescribeTable("returns the error from Compare",
		func(expectedError error, needle interface{}, haystack interface{}, fields ...string) {
			_, err := Search(needle, haystack, fields...)
			Expect(err).To(MatchError(expectedError))
		},
		Entry("Invalid type given", ErrInvalidType, "test", []string{"test", "test2"}),

		// attribute names are only checked when there is any entry in the haystack, I think that's ok
		// -- Mara @LittleFox94 Grosch, 2022-03-28
		Entry("Invalid attribute given", ErrKeyNotFound, lbaasv1.Frontend{}, []lbaasv1.Frontend{{}}, "Test"),

		// the types being the same is only checked when there is any entry in the haystack, I think that's ok
		// -- Mara @LittleFox94 Grosch, 2022-03-28
		Entry("Different types", ErrDifferentTypes, lbaasv1.Frontend{}, []lbaasv1.Backend{{}}, "Name"),
	)
})
