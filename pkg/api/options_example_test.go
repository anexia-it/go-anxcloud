package api_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/mock"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
)

func ExampleAutoTag() {
	a := mock.NewMockAPI()

	vlan := vlanv1.VLAN{DescriptionCustomer: "mocked VLAN"}
	if err := a.Create(context.TODO(), &vlan, api.AutoTag("foo", "bar", "baz")); err != nil {
		taggingErr := &api.ErrTaggingFailed{}
		if errors.As(err, taggingErr) {
			log.Fatalf("object successfully created but tagging failed: %s", taggingErr.Error())
		} else {
			log.Fatalf("unknown error occurred: %s", err)
		}
	}

	// Note that `a.Inspect` is only available when using the mock client implementation.
	tags := a.Inspect(vlan.Identifier).Tags()
	sort.Strings(tags)

	fmt.Println(strings.Join(tags, ", "))
	// Output: bar, baz, foo
}
