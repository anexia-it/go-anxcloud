package common

import (
	"encoding/json"
)

// PartialResource represents a linked resource
type PartialResource struct {
	Identifier string `json:"identifier,omitempty"`
	Name       string `json:"name,omitempty"`
}

// MarshalJSON unfolds the resource to its identifier
func (pr PartialResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(pr.Identifier)
}
