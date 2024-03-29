// Code generated by go.anx.io/go-anxcloud/tools object-generator - DO NOT EDIT!

package v1

import (
	"context"
)

// GetIdentifier returns the primary identifier of a Cluster object
func (o *Cluster) GetIdentifier(ctx context.Context) (string, error) {
	return o.Identifier, nil
}

// GetIdentifier returns the primary identifier of a kubeconfig object
func (o *kubeconfig) GetIdentifier(ctx context.Context) (string, error) {
	return o.Cluster, nil
}

// GetIdentifier returns the primary identifier of a NodePool object
func (o *NodePool) GetIdentifier(ctx context.Context) (string, error) {
	return o.Identifier, nil
}
