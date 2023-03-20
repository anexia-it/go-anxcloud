package compare

import (
	"fmt"

	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleCompare() {
	lb := lbaasv1.LoadBalancer{Identifier: "some identifier string"}
	a := lbaasv1.Frontend{
		Name:         "Frontend A",
		Mode:         lbaasv1.TCP,
		LoadBalancer: &lb,
	}

	b := lbaasv1.Frontend{
		Name:         "Frontend B",
		Mode:         lbaasv1.TCP,
		LoadBalancer: &lb,
	}

	differences, err := Compare(a, b, "Name", "Mode", "LoadBalancer.Identifier")
	if err != nil {
		fmt.Printf("Error comparing the objects: %v\n", err)
	} else if len(differences) > 0 {
		for _, difference := range differences {
			fmt.Printf("Difference between objects on attribute %q: %q != %q\n", difference.Key, difference.A, difference.B)
		}
	} else {
		fmt.Printf("Found no difference between A and B\n")
	}

	// Output:
	// Difference between objects on attribute "Name": "Frontend A" != "Frontend B"
}

func ExampleCompare_nestedNil() {
	lb := lbaasv1.LoadBalancer{Identifier: "some identifier string"}
	a := lbaasv1.Frontend{
		Name: "Frontend A",
		Mode: lbaasv1.TCP,
	}

	b := lbaasv1.Frontend{
		Name:         "Frontend B",
		Mode:         lbaasv1.TCP,
		LoadBalancer: &lb,
	}

	differences, err := Compare(a, b, "Name", "Mode", "LoadBalancer.Identifier")
	if err != nil {
		fmt.Printf("Error comparing the objects: %v\n", err)
	} else if len(differences) > 0 {
		for _, difference := range differences {
			fmt.Printf("Difference between objects on attribute %q: '%v' != '%v'\n", difference.Key, difference.A, difference.B)
		}
	} else {
		fmt.Printf("Found no difference between A and B\n")
	}

	// Output:
	// Difference between objects on attribute "Name": 'Frontend A' != 'Frontend B'
	// Difference between objects on attribute "LoadBalancer.Identifier": '<nil>' != 'some identifier string'
}

var _ = Describe("Compare", func() {
	It("works fine with one arg being a pointer and the other not", func() {
		diff, err := Compare(lbaasv1.Frontend{Name: "test"}, &lbaasv1.Frontend{Name: "test"}, "Name")
		Expect(err).NotTo(HaveOccurred())
		Expect(diff).To(BeEmpty())

		diff, err = Compare(&lbaasv1.Frontend{Name: "test"}, lbaasv1.Frontend{Name: "test"}, "Name")
		Expect(err).NotTo(HaveOccurred())
		Expect(diff).To(BeEmpty())
	})

	It("errors out on anything not a struct or pointer to a struct", func() {
		_, err := Compare("test", "test")
		Expect(err).To(MatchError(ErrInvalidType))

		_, err = Compare(true, false)
		Expect(err).To(MatchError(ErrInvalidType))

		s := "test"
		_, err = Compare(&s, &s)
		Expect(err).To(MatchError(ErrInvalidType))
	})

	It("errors out when trying to compare different types", func() {
		_, err := Compare(lbaasv1.Frontend{Name: "test"}, lbaasv1.Backend{Name: "test"}, "Name")
		Expect(err).To(MatchError(ErrDifferentTypes))
	})

	It("errors out when given non-existing attribute", func() {
		_, err := Compare(lbaasv1.Frontend{Name: "test"}, lbaasv1.Frontend{Name: "test"}, "Test")
		Expect(err).To(MatchError(ErrKeyNotFound))
	})
})
