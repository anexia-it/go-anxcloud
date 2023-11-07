package v1

type TransferProtocol string

const (
	TransferProtocolHTTP TransferProtocol = "http"
)

// anxcloud:object

// API represents Frontier's root object and contains a collection of endpoints.
// The API defines the transfer protocol, such as HTTP and HTTPS, for all containing endpoints.
type API struct {
	omitResponseDecodeOnDestroy
	Identifier           string  `json:"identifier,omitempty" anxcloud:"identifier"`
	Name                 string  `json:"name,omitempty"`
	Description          *string `json:"description,omitempty"`
	TransferProtocol     string  `json:"transfer_protocol,omitempty"`
	DeploymentIdentifier string  `json:"deployment_identifier,omitempty"`
}
