package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (l *Location) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Locations can only be retrieved via the public engine, nothing else
	if op != types.OperationGet && op != types.OperationList {
		return nil, api.ErrOperationNotSupported
	}

	return url.Parse("/api/core/v1/location.json")
}

func (l *Location) DecodeAPIResponse(ctx context.Context, body io.Reader) error {
	type apiLocation struct {
		Location
		Lat *string `json:"lat"`
		Lon *string `json:"lon"`
	}

	loc := apiLocation{}
	if err := json.NewDecoder(body).Decode(&loc); err != nil {
		return err
	}

	if loc.Lat != nil {
		lat, err := strconv.ParseFloat(*loc.Lat, 64)
		if err != nil {
			return fmt.Errorf("error parsing latitude: %w", err)
		}

		loc.Location.Latitude = &lat
	}

	if loc.Lon != nil {
		lon, err := strconv.ParseFloat(*loc.Lon, 64)
		if err != nil {
			return fmt.Errorf("error parsing longitude: %w", err)
		}

		loc.Location.Longitude = &lon
	}

	*l = loc.Location
	return nil
}
