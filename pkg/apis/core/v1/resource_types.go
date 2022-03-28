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
