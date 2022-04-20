package v1

import "encoding/json"

// Type is part of info.
type Type struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// anxcloud:object:hooks=ResponseDecodeHook

// Resource contains all information about a resource.
type Resource struct {
	Identifier string          `json:"identifier" anxcloud:"identifier"`
	Name       string          `json:"name"`
	Type       Type            `json:"resource_type"`
	Tags       []string        `json:"tags"`
	Attributes json.RawMessage `json:"attributes"`
}

// anxcloud:object:hooks=ResponseFilterHook,RequestBodyHook,FilterRequestURLHook

// ResourceWithTag is a virtual Object used to add (Create) or remove (Destroy) a tag to/from a Resource.
type ResourceWithTag struct {
	// Identifier of the Resource which tags to change
	Identifier string `anxcloud:"identifier"`

	// Name of the Tag to add or remove from the resource
	Tag string
}
