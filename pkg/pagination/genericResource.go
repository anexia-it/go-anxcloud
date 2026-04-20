package pagination

type IGenericResource interface {
	GetIdentifier() string
	GetName() string
}

type GenericResource struct {
	IGenericResource
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (g GenericResource) GetIdentifier() string {
	return g.Identifier
}

func (g GenericResource) GetName() string {
	return g.Name
}
