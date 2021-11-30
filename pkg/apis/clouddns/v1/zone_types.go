package v1

import "time"

type Revision struct {
	CreatedAt  time.Time `json:"created_at"`
	Identifier string    `json:"identifier"`
	ModifiedAt time.Time `json:"modified_at"`
	Records    []Record  `json:"records"`
	Serial     int       `json:"serial"`
	State      string    `json:"state"`
}

type DNSServer struct {
	// Required - DNS Server name (FQDN).
	Server string `json:"server"`

	// DNS Server alias
	Alias string `json:"alias"`
}

// anxcloud:object:hooks=ResponseDecodeHook,RequestFilterHook,RequestBodyHook,ResponseFilterHook,PaginationSupportHook

type Zone struct {
	// Zone name
	Name string `json:"name,omitempty" anxcloud:"identifier"`

	// Required - Is master flag
	// Flag designating if CloudDNS operates as master or slave.
	IsMaster bool `json:"master"`

	// Required - DNSSEC mode
	// DNSSEC mode (master-only) ["managed" or "unvalidated"].
	DNSSecMode string `json:"dnssec_mode"`

	// Required - Admin email address
	// Admin email address used in SOA record.
	AdminEmail string `json:"admin_email"`

	// Required - Refresh value
	// Refresh value used in SOA record.
	Refresh int `json:"refresh"`

	// Required - Retry value
	//Retry value used in SOA record.
	Retry int `json:"retry"`

	// Required - Expire value
	// Expire value used in SOA record.
	Expire int `json:"expire"`

	// Required - Time to live
	// Default TTL for NS records.
	TTL int `json:"ttl"`

	// Master Name Server
	MasterNS string `json:"master_ns,omitempty"`

	// IP addresses allowed to initiate domain transfer (DNS NOTIFY).
	NotifyAllowedIPs []string `json:"notify_allowed_ips,omitempty"`

	// Configured DNS servers (empty means default servers).
	DNSServers []DNSServer `json:"dns_servers,omitempty"`

	Customer        string     `json:"customer"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PublishedAt     time.Time  `json:"published_at"`
	IsEditable      bool       `json:"is_editable"`
	ValidationLevel int        `json:"validation_level"`
	DeploymentLevel int        `json:"deployment_level"`
	Revisions       []Revision `json:"revisions"`
	CurrentRevision string     `json:"current_revision,omitempty"`
}
