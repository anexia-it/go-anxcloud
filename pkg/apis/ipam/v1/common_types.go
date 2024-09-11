package v1

// Family defines an address family/IP version.
type Family int

const (
	// FamilyAll reserves both IPv4 and IPv6 addresses.
	FamilyAll Family = 0

	// FamilyIPv4 denotes a Prefix as being an IPv4 prefix.
	FamilyIPv4 Family = 4

	// FamilyIPv6 denotes a Prefix as being an IPv6 prefix.
	FamilyIPv6 Family = 6
)

// Status describes the status of an object, if it is being created, ready to be used, ...
type Status string

const (
	// StatusActive marks an address or prefix as being allocated and ready to be used.
	StatusActive Status = "Active"

	// StatusPending marks an address or prefix as being worked on and not yet ready to be used.
	StatusPending Status = "Pending"

	// StatusFailed marks an address or prefix as being failed, you have to contact support.
	StatusFailed Status = "Failed"

	// StatusMarkedForDeletion marks an address or prefix as being in the process of being deleted.
	StatusMarkedForDeletion Status = "Marked for deletion"

	// StatusInactive marks an address as inactive.
	StatusInactive Status = "Inactive"
)

// AddressType defines if an address or prefix is from public/internet routable or private address space.
type AddressType string

const (
	// TypePublic specifies an address or prefix as being from public/internet-routable address space.
	TypePublic AddressType = "0"

	// TypePrivate specifies an address or prefix as being from private (RFC1918) address space.
	TypePrivate AddressType = "1"
)

// MarshalJSON encodes Type into JSON and is required because the API expects a number, but we use a string
// to differentiate between TypePublic and no Type set for filtering.
func (t AddressType) MarshalJSON() ([]byte, error) {
	return []byte(t), nil
}

// UnmarshalJSON decodes the given JSON value into the Type it is called on. See MarshalJSON why it is needed.
func (t *AddressType) UnmarshalJSON(data []byte) error {
	*t = AddressType(string(data))
	return nil
}
