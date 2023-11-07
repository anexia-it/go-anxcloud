package v1

// anxcloud:object

// Deployment represents a published version of a Frontier API with all its endpoints
// and actions exactly as it was at the time it was deployed.
// Note that using api.Create on a `Deployment` only sets the Deployments `Identifier`.
// Use api.Get on the same struct instance to retrieve all data.
type Deployment struct {
	omitResponseDecodeOnDestroy
	Identifier    string `json:"identifier,omitempty" anxcloud:"identifier"`
	APIIdentifier string `json:"api_identifier,omitempty"`
	Name          string `json:"name,omitempty"`
	Slug          string `json:"slug,omitempty"`
	State         string `json:"state,omitempty"`
}
