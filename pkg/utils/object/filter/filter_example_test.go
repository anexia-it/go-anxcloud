package filter

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

type parentObject struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
}

func (o *parentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/v1/parent")
}

// testObject is the definition of our test object. We define Name and Parent to be filterable. The filter for
// Name is named "something" instead of "name". When having an Object marked as filterable, the identifier of
// the Object is used, if it is set at all.
type testObject struct {
	Identifier  string  `json:"identifier" anxcloud:"identifier"`
	Name        string  `json:"name" anxcloud:"filterable,something"`
	Description *string `json:"desc" anxcloud:"filterable"`

	Parent        parentObject  `json:"parent" anxcloud:"filterable"`
	PointerParent *parentObject `json:"ptrParent" anxcloud:"filterable"`

	// just a random non-special field
	SomeString        string `json:"someString"`
	SomeNumber        int    `json:"someNumber"`
	SomeNumberPointer *int   `json:"someNumberPtr"`
}

// EndpointURL builds the URL for interacting with our test Object. filter.Helper is most useful in EndpointURL
// when doing a List operation.
func (o *testObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse("/v1/test")

	if op == types.OperationList {
		filterHelper, err := NewHelper(o)
		if err != nil {
			return nil, err
		}

		query := filterHelper.BuildQuery()
		u.RawQuery = query.Encode()
	}

	return u, nil
}

func ExampleHelper_extended() {
	parent := parentObject{
		Identifier: "parentIdentifier",
	}

	test := testObject{
		Identifier: "testIdentifier",
		Name:       "test object",
		Parent:     parent,
	}

	// this is called by the generic client, here just to see what we do by using filter.Helper
	u, err := test.EndpointURL(types.ContextWithOperation(context.TODO(), types.OperationList))
	if err != nil {
		fmt.Printf("Error creating endpoint URL: %v\n", err)
		os.Exit(-1)
	}

	// code below is validating if the correct filters were set
	query := u.Query()

	if _, ok := query["something"]; !ok {
		fmt.Printf("Name filter not set but it should be\n")
		os.Exit(-1)
	}

	if _, ok := query["parent"]; !ok {
		fmt.Printf("Parent filter not set but it should be\n")
		os.Exit(-1)
	}

	fmt.Printf("Name filter configured for \"something\": %q\n", query.Get("something"))
	fmt.Printf("Parent filter configured for \"parent\": %q\n", query.Get("parent"))

	// Output:
	// Name filter configured for "something": "test object"
	// Parent filter configured for "parent": "parentIdentifier"
}
