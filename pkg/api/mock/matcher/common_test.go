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
func (o *testObject) GetIdentifier(context.Context) (string, error)     { return o.Identifier, nil }

type testObject2 struct {
	Identifier string `anxcloud:"identifier"`
}

func (o *testObject2) EndpointURL(ctx context.Context) (*url.URL, error) { return nil, nil }
func (o *testObject2) GetIdentifier(context.Context) (string, error)     { return o.Identifier, nil }

type testObjectNotStruct string

func (o *testObjectNotStruct) EndpointURL(ctx context.Context) (*url.URL, error) { return nil, nil }
func (o testObjectNotStruct) GetIdentifier(ctx context.Context) (string, error) {
	return string(o), nil
}
