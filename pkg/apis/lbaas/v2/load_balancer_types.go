package v2

import (
	"encoding/json"

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// anxcloud:object

// LoadBalancer holds the information of a load balancing configuration within a Cluster
type LoadBalancer struct {
	gs.GenericService
	gs.HasState

	Identifier      string                  `json:"identifier,omitempty" anxcloud:"identifier"`
	Name            string                  `json:"name,omitempty"`
	Generation      int                     `json:"generation,omitempty"`
	Cluster         *common.PartialResource `json:"cluster,omitempty" anxcloud:"filterable"`
	FrontendIPs     *gs.PartialResourceList `json:"frontend_ips,omitempty"`
	SSLCertificates *gs.PartialResourceList `json:"ssl_certificates,omitempty"`
	// Definition configures the load balancer's frontends and backends.
	// This field is currently unstable and requires an update of go-anxcloud
	// in the near future.
	Definition *Definition `json:"definition,omitempty"`
}

type FrontendProtocol string

const (
	FrontendProtocolTCP FrontendProtocol = "TCP"
)

type BackendProtocol string

const (
	BackendProtocolTCP   BackendProtocol = "TCP"
	BackendProtocolPROXY BackendProtocol = "PROXY"
)

type definition struct {
	Frontends []Frontend `json:"frontends,omitempty"`
	Backends  []Backend  `json:"backends,omitempty"`
}

type Definition definition

func (d *Definition) MarshalJSON() ([]byte, error) {
	inner := definition(*d)

	data, err := json.Marshal(&inner)
	if err != nil {
		return nil, err
	}

	return json.Marshal(string(data))
}

func (d *Definition) UnmarshalJSON(data []byte) error {
	var innerString string
	if err := json.Unmarshal(data, &innerString); err != nil {
		return err
	}

	var inner definition
	if err := json.Unmarshal([]byte(innerString), &inner); err != nil {
		return err
	}

	*d = Definition(inner)

	return nil
}

// Define ports and protocols exposed to the public side of the Load Balancer
type Frontend struct {
	// Name of the Frontend
	Name string `json:"name"`
	// Frontend service protocol
	Protocol FrontendProtocol `json:"protocol"`
	// Configure frontend - backend relation
	Backend FrontendBackend `json:"backend,omitempty"`
	// TCP specific frontend configuration
	TCP *FrontendTCP `json:"tcp,omitempty"`
}

// TCP specific frontend configuration
type FrontendTCP struct {
	// Port for the frontend to listen to
	Port uint16 `json:"port,omitempty"`
}

// Configure frontend - backend relation
type FrontendBackend struct {
	// Backend service protocol
	Protocol BackendProtocol `json:"protocol"`
	// TCP specific backend configuration
	TCP *FrontendBackendTCP `json:"tcp,omitempty"`
}

// TCP specific backend configuration
type FrontendBackendTCP struct {
	Port uint16 `json:"port"`
}

// Define backend services and connect them to frontends
type Backend struct {
	// Name of the Backend
	Name string `json:"name"`
	// IP addresses of the backend service
	IPs []string `json:"ips,omitempty"`
	// List of frontends connected to the backend
	Frontends []BackendFrontend `json:"frontends,omitempty"`
}

type BackendFrontend struct {
	// Name of the frontend to be connected to the backend
	Name string `json:"name"`
}
