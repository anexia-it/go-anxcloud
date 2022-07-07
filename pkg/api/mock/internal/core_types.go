package internal

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
	clouddnsv1 "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	vlanv1 "go.anx.io/go-anxcloud/pkg/apis/vlan/v1"
)

// ObjectCoreType maps types.Object to *corev1.Type
// returns nil if no mapping was found
// TODO: (SYSENG-1333) have generic client bindings register themselves with their Type Identifier so we can use that for other purposes, too
func ObjectCoreType(o types.Object) *corev1.Type {
	switch o.(type) {
	case *lbaasv1.Backend:
		return &corev1.Type{Identifier: "33164a3066a04a52be43c607f0c5dd8c"}
	case *lbaasv1.Bind:
		return &corev1.Type{Identifier: "bd24def982aa478fb3352cb5f49aab47"}
	case *lbaasv1.Frontend:
		return &corev1.Type{Identifier: "da9d14b9d95840c08213de67f9cee6e2"}
	case *lbaasv1.Server:
		return &corev1.Type{Identifier: "01f321a4875446409d7d8469503a905f"}
	case *vlanv1.VLAN:
		return &corev1.Type{Identifier: "cf8e4dac56894afaa3244f6911bb62be"}
	case *clouddnsv1.Zone:
		return &corev1.Type{Identifier: "9190d73d8f4f42b5ad29e1a057f184fc"}
	}
	return nil
}
