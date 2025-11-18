package bind

import "go.anx.io/go-anxcloud/pkg/lbaas/common"

type Definition struct {
	Name               string       `json:"name"`
	State              common.State `json:"state"`
	Frontend           string       `json:"frontend"`
	Address            string       `json:"address"`
	Port               int          `json:"port"`
	SSL                bool         `json:"ssl"`
	SSLCertificatePath string       `json:"ssl_certificate_path"`
}
