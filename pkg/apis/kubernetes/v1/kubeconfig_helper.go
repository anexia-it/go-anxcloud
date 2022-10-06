package v1

import (
	"context"
	"fmt"
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

const getKubeConfigCheckInterval = 5 * time.Second

type kubeconfig struct {
	Cluster string `json:"cluster" anxcloud:"identifier"`
}

// GetKubeConfig retrieves the kubeconfig
func GetKubeConfig(ctx context.Context, a api.API, clusterID string) (string, error) {
	kubeconfigRequested := false

	ticker := time.NewTicker(getKubeConfigCheckInterval)
	defer ticker.Stop()

	cluster := Cluster{Identifier: clusterID}

	for {
		if err := a.Get(ctx, &cluster); err != nil {
			return "", fmt.Errorf("failed to get cluster: %w", err)
		}

		if cluster.KubeConfig != nil {
			return pointer.StringVal(cluster.KubeConfig), nil
		}

		if !kubeconfigRequested {
			if err := RequestKubeConfig(ctx, a, clusterID); err != nil {
				return "", fmt.Errorf("failed to request kubeconfig: %w", err)
			}
			kubeconfigRequested = true
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

// RequestKubeConfig triggers the "Request kubeconfig" automation rule
func RequestKubeConfig(ctx context.Context, a api.API, clusterID string) error {
	return a.Create(ctx, &kubeconfig{Cluster: clusterID})
}

// RemoveKubeConfig triggers the "Remove kubeconfig" automation rule
func RemoveKubeConfig(ctx context.Context, a api.API, clusterID string) error {
	return a.Destroy(ctx, &kubeconfig{Cluster: clusterID})
}
