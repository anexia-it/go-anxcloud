package compare

import (
	"fmt"

	"go.anx.io/go-anxcloud/pkg/api/types"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleReconcile() {
	lb := lbaasv1.LoadBalancer{Identifier: "LoadBalancer identifier"}

	// array of objects that should exist, created by business logic
	targetObjects := []lbaasv1.Frontend{
		{
			Name:         "Frontend A",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
		{
			Name:         "Frontend B",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
	}

	// array of objects found in the Engine as currently existing
	existingObjects := []lbaasv1.Frontend{
		{
			Name:         "Frontend A",
			Identifier:   "Frontend identifier A",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
		{
			Name:         "Frontend C",
			Identifier:   "Frontend identifier C",
			Mode:         lbaasv1.TCP,
			LoadBalancer: &lb,
		},
	}

	// those will receive the objects to create and destroy
	var toCreate []types.Object
	var toDestroy []types.Object

	err := Reconcile(
		// array of objects we want in the end and we currently have
		targetObjects, existingObjects,
		// output arrays of objects to create and destroy to change reality into our desired state
		&toCreate, &toDestroy,
		// attributes to compare
		"Name", "Mode", "LoadBalancer.Identifier",
	)
	if err != nil {
		fmt.Printf("Error reconciling objects: %v\n", err)
		return
	}

	fmt.Printf("Found %v Objects to create and %v Objects to destroy\n", len(toCreate), len(toDestroy))

	for _, c := range toCreate {
		fmt.Printf("Going to create Object named %q\n", c.(*lbaasv1.Frontend).Name)
	}

	for _, c := range toDestroy {
		fmt.Printf("Going to destroy Object named %q\n", c.(*lbaasv1.Frontend).Name)
	}

	for _, f := range targetObjects {
		if f.Identifier != "" {
			fmt.Printf("Identifier of Object named %q was retrieved: %q\n", f.Name, f.Identifier)
		}
	}

	// Output:
	// Found 1 Objects to create and 1 Objects to destroy
	// Going to create Object named "Frontend B"
	// Going to destroy Object named "Frontend C"
	// Identifier of Object named "Frontend A" was retrieved: "Frontend identifier A"
}

var _ = Describe("Reconcile", func() {
	DescribeTable("errors out for invalid input",
		func(expectedError error, target, existing interface{}, compareAttributes ...string) {
			err := Reconcile(target, existing, nil, nil, compareAttributes...)
			Expect(err).To(MatchError(expectedError))
		},
		Entry("existing not an array", ErrInvalidType, []lbaasv1.Frontend{}, lbaasv1.Frontend{}),
		Entry("target not an array", ErrInvalidType, lbaasv1.Frontend{}, []lbaasv1.Frontend{}),

		// these errors are only checked when there are entries in the target and existing arrays, I think that's ok
		// -- Mara @LittleFox94 Grosch, 2022-03-28
		Entry("invalid attribute", ErrKeyNotFound, []lbaasv1.Frontend{{}}, []lbaasv1.Frontend{{}}, "Test"),
	)
})
