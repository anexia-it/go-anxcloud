package matcher

import (
	"context"
	"fmt"
	"log"

	clouddnsv1 "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/api/mock"
)

// Object is a custom gomega matcher which is used to test for matching objects
// on the MockAPI instance by comparing the type and provided field names
func ExampleObject() {
	// defer GinkgoRecover() is not needed when assertion is done within Ginkgo environment
	defer GinkgoRecover()

	a := mock.NewMockAPI()

	if err := a.Create(context.TODO(), &clouddnsv1.Record{ZoneName: "example.com", Name: "webmail"}); err != nil {
		log.Fatalf("failed to create DNS record: %s", err)
	}

	// a.Existing() retrieves all objects which aren't currently flagged as destroyed
	Expect(a.Existing()).To(
		ContainElement(
			Object(&clouddnsv1.Record{ZoneName: "example.com"}, "ZoneName"),
		),
	)

	fmt.Println("Success")

	// Output: Success
}

// Destroyed is a custom gomega matcher which is used to test
// if an object was destroyed
func ExampleDestroyed() {
	// defer GinkgoRecover() is not needed when assertion is done within Ginkgo environment
	defer GinkgoRecover()

	a := mock.NewMockAPI()

	record := clouddnsv1.Record{ZoneName: "example.com", Name: "webmail"}

	if err := a.Create(context.TODO(), &record); err != nil {
		log.Fatalf("failed to create record: %s", err)
	}

	if err := a.Destroy(context.TODO(), &clouddnsv1.Record{Identifier: record.Identifier}); err != nil {
		log.Fatalf("failed destroy record: %s", err)
	}

	Expect(a.All()).To(
		ContainElement(
			SatisfyAll(
				Destroyed(1),
				Object(&clouddnsv1.Record{Identifier: record.Identifier}, "Identifier"),
			),
		),
	)

	fmt.Println("Success")

	// Output: Success
}

// TaggedWith is a custom gomega matcher which is used to test
// if an object was tagged with the provided tags
func ExampleTaggedWith() {
	// defer GinkgoRecover() is not needed when assertion is done within Ginkgo environment
	defer GinkgoRecover()

	a := mock.NewMockAPI()

	zone := clouddnsv1.Zone{Name: "example.com"}

	if err := a.Create(context.TODO(), &zone); err != nil {
		log.Fatalf("failed to create DNS zone: %s", err)
	}

	if err := corev1.Tag(context.Background(), a, &zone, "Tag-0", "Tag-1", "Tag-2"); err != nil {
		log.Fatalf("failed to tag zone: %s", err)
	}

	Expect(a.Existing()).To(
		ContainElement(
			SatisfyAll(
				Object(&zone, "Name"),
				TaggedWith("Tag-0", "Tag-1"),
			),
		),
	)

	fmt.Println("Success")

	// Output: Success
}
