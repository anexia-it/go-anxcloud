// DO NOT EDIT, auto generated

package v1

import (
	"context"
)

// GetIdentifier returns the primary identifier of a ACL object
func (x *ACL) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Backend object
func (x *Backend) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Bind object
func (x *Bind) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Frontend object
func (x *Frontend) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a LoadBalancer object
func (x *LoadBalancer) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Rule object
func (x *Rule) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a RuleInfo object
func (x *RuleInfo) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}

// GetIdentifier returns the primary identifier of a Server object
func (x *Server) GetIdentifier(ctx context.Context) (string, error) {
	return x.Identifier, nil
}
