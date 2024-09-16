package v1

// anxcloud:object:hook=SupportPaginationHook

// IPs represents the endpoint to check for available IPs in a VLAN.
type IPs struct {
	Data []IP `json:"data"`
	//Data []json.RawMessage `json:"data"`

	CurrentPage    uint `json:"page"`
	TotalPages     uint `json:"total_pages"`
	TotalItems     uint `json:"total_items"`
	EntriesPerPage uint `json:"limit"`

	//// LocationIdentifier is required to query available IPs.
	//LocationIdentifier string `json:"location_identifier"`
	//
	//// VLANIdentifier is required to query available IPs.
	//VLANIdentifier string `json:"vlan_identifier"`
	//// Dummy?
	//Identifier string `json:"identifier" anxcloud:"identifier"`
}

type IP struct {
	Identifier string `json:"identifier"`
	Text       string `json:"text"`
	Prefix     string `json:"prefix"`
}
