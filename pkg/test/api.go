// Package test contains API functionality for "testing" the API.
package test

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/test/echo"
)

// API contains methods for "testing" the API.
type API interface {
	Echo() echo.API
}

type api struct {
	echo echo.API
}

func (a api) Echo() echo.API {
	return a.echo
}

// NewAPI creates a new test API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		echo.NewAPI(c),
	}
}
