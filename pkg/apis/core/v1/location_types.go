package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// anxcloud:object

// Location describes a Anexia site where resources can be deployed.
type Location struct {
	Identifier  string  `json:"identifier" anxcloud:"identifier"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	CountryCode string  `json:"country"`
	CityCode    string  `json:"city_code"`
	Latitude    *string `json:"lat"`
	Longitude   *string `json:"lon"`
}

// GetIdentifier returns the objects identifier
func (l *Location) GetIdentifier(ctx context.Context) (string, error) {
	if l.Identifier != "" {
		return l.Identifier, nil
	}

	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return "", err
	}

	if op == types.OperationGet {
		return l.Code, nil
	}

	return "", types.ErrUnidentifiedObject
}
