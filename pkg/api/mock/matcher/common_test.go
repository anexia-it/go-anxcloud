package matcher

import (
	"context"
	"net/url"
)

type testObject struct {
	Identifier string `anxcloud:"identifier"`
	TestFieldA string `json:",omitempty"`
	TestFieldB string `json:",omitempty"`
}

func (o *testObject) EndpointURL(ctx context.Context) (*url.URL, error) { return nil, nil }

type testObject2 struct {
	Identifier string `anxcloud:"identifier"`
}

func (o *testObject2) EndpointURL(ctx context.Context) (*url.URL, error) { return nil, nil }

type testObjectNotStruct string

func (o *testObjectNotStruct) EndpointURL(ctx context.Context) (*url.URL, error) { return nil, nil }