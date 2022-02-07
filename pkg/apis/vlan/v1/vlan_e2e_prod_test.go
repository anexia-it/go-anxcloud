//go:build integration
// +build integration

package v1

import (
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"
)

const (
	waitTimeout  = 5 * time.Minute
	retryTimeout = 15 * time.Second
)

func e2eApiClient() (api.API, error) {
	return api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
}

// below are the functions for setting up the mock, empty for the prod E2E version
func prepareCreate(desc string) {}

func prepareGet(desc string, vmProvisioning bool) {}

func prepareList(desc string, vmProvisioning bool) {}

func prepareUpdate(desc string, vmProvisioning bool) {}

func prepareDelete() {}

func prepareEventuallyActive(desc string, vmProvisioning bool) {}

func prepareEventuallyDeleted(desc string, vmProvisioning bool) {}

func prepareDeleting() {}
