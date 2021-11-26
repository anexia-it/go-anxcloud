package v1

// Type is part of info.
type Type struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// anxcloud:object

// Info contains all information about a resource.
type Info struct {
	Identifier string   `json:"identifier" anxcloud:"identifier"`
	Name       string   `json:"name"`
	Type       Type     `json:"resource_type"`
	Tags       []string `json:"tags"`
}
