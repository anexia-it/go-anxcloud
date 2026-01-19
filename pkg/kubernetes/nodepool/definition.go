package nodepool

import "go.anx.io/go-anxcloud/pkg/apis/common/gs"

const (
	OSFlatcar = "Flatcar Linux"
	Gibibyte  = 1024 * 1024 * 1024
)

var (
	StatePending = gs.State{ID: "2", Text: "Pending", Type: gs.StateTypePending}
)

type Definition struct {
	Name            string   `json:"name"`
	State           gs.State `json:"state"`
	ClusterID       string   `json:"cluster"`
	Replicas        uint64   `json:"replicas"`
	CPUs            uint64   `json:"cpus"`
	MemoryBytes     uint64   `json:"memory"`
	DiskSizeBytes   uint64   `json:"disk_size"`
	OperatingSystem string   `json:"operating_system"`
}
