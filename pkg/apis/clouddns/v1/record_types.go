package v1

// anxcloud:object:hooks=ResponseDecodeHook,PaginationSupportHook,RequestFilterHook

type Record struct {
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	ZoneName   string `json:"-"`
	Immutable  bool   `json:"immutable,omitempty"`
	// Name of the DNS record.
	// Use "@" to select the domain root. Creation of records with an empty Name field is not supported.
	Name   string `json:"name"`
	RData  string `json:"rdata"`
	Region string `json:"region"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
}
