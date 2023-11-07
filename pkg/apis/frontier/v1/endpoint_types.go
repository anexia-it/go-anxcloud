package v1

// anxcloud:object

// Endpoint represents a path within an HTTP-based API and contains a collection of actions.
type Endpoint struct {
	omitResponseDecodeOnDestroy
	Identifier    string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name          string `json:"name,omitempty"`
	Path          string `json:"path,omitempty"`
	APIIdentifier string `json:"api_identifier,omitempty"`
}
