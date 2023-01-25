package gs

import (
	"encoding/json"
	"strings"

	"go.anx.io/go-anxcloud/pkg/apis/common"
)

// PartialResourceList represents a linked resource list in GS
type PartialResourceList []common.PartialResource

// MarshalJSON unfolds the resources to their identifier as a comma-separated string
func (prl PartialResourceList) MarshalJSON() ([]byte, error) {
	resourceIdentifiers := make([]string, 0, len(prl))
	for _, pr := range prl {
		resourceIdentifiers = append(resourceIdentifiers, pr.Identifier)
	}

	return json.Marshal(strings.Join(resourceIdentifiers, ","))
}
