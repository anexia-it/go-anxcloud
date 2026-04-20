package genericResource

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/pagination"
)

type IGenericResource interface {
	GetIdentifier() string
	GetName() string
}

type API[R any, D any] interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]Identity, error)
	GetByID(ctx context.Context, identifier string) (R, error)
	Create(ctx context.Context, definition D) (R, error)
	Update(ctx context.Context, identifier string, definition D) (R, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type Identity struct {
	IGenericResource
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (g Identity) GetIdentifier() string {
	return g.Identifier
}

func (g Identity) GetName() string {
	return g.Name
}
