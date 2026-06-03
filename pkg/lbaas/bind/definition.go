package bind

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Bind = v1.Bind
type Definition struct {
	Name               string       `json:"name,omitempty"`
	State              common.State `json:"state,omitempty"`
	Frontend           string       `json:"frontend,omitempty"`
	Address            string       `json:"address,omitempty"`
	Port               int          `json:"port,omitempty"`
	SSL                bool         `json:"ssl,omitempty"`
	SSLCertificatePath string       `json:"ssl_certificate_path,omitempty"`
}
