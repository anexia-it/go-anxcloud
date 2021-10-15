package types

// Operation to do on the engine with an object. Users are expected to compare values
// of this type to the Operation(Get|Create|...) constants in this package.
type Operation string

const (
	// OperationGet is used to retrieve the given identified object from the engine.
	OperationGet Operation = "Get"

	// OperationCreate is used to create the given object on the engine.
	OperationCreate Operation = "Create"

	// OperationUpdate is used to update the given identified object on the engine with the newly given data.
	OperationUpdate Operation = "Update"

	// OperationDestroy is used to destroy the identified object from the engine.
	OperationDestroy Operation = "Destroy"

	// OperationList is used to retrieve objects with attributes matching the ones in the given object from
	// the engine.
	OperationList Operation = "List"
)

// GetOptions contains options valid for Get operations.
type GetOptions struct {
	commonOptions
}

// ListOptions contains options valid for List operations.
type ListOptions struct {
	commonOptions
	ObjectChannel *ObjectChannel

	Paged          bool
	Page           uint
	EntriesPerPage uint
	PageInfo       *PageInfo
}

// CreateOptions contains options valid for Create operations.
type CreateOptions struct {
	commonOptions
}

// UpdateOptions contains options valid for Update operations.
type UpdateOptions struct {
	commonOptions
}

// DestroyOptions contains options valid for Destroy operations.
type DestroyOptions struct {
	commonOptions
}
