package v1

// anxcloud:object

// Applications are an easy way to bring more structure to your configured functions by grouping them.
type Application struct {
	omitResponseDecodeOnDestroy
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	Name       string `json:"name,omitempty"`
}
