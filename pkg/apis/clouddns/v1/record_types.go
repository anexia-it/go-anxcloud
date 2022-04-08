package v1

// anxcloud:object:hooks=ResponseDecodeHook,PaginationSupportHook,FilterRequestURLHook

type Record struct {
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	ZoneName   string `json:"-"`
	Immutable  bool   `json:"immutable,omitempty"`
	Name       string `json:"name"`
	RData      string `json:"rdata"`
	Region     string `json:"region"`
	TTL        int    `json:"ttl"`
	Type       string `json:"type"`
}
