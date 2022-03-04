package v1

// Family defines an address family/IP version.
type Family int

const (
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
)

// AddressSpace defines if an address or prefix is from public/internet routable or private address space.
type AddressSpace string

const (
	// TypePublic specifies an address or prefix as being from public/internet-routable address space.
	AddressSpacePublic AddressSpace = "0"

	// TypePrivate specifies an address or prefix as being from private (RFC1918) address space.
	AddressSpacePrivate AddressSpace = "1"
)

// MarshalJSON encodes Type into JSON and is required because the API expects a number, but we use a string
// to differentiate between AddressSpacePublic and no AddressSpace set for filtering.
func (t AddressSpace) MarshalJSON() ([]byte, error) {
	return []byte(t), nil
}

// UnmarshalJSON decodes the given JSON value into the AddressSpace it is called on. See MarshalJSON why it is needed.
func (t *AddressSpace) UnmarshalJSON(data []byte) error {
	*t = AddressSpace(string(data))
	return nil
}
