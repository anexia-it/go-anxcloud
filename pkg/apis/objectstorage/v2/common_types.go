package v2

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
)

// AutomationRule represents an automation rule in the S3 API.
type AutomationRule struct {
	Identifier     string            `json:"identifier,omitempty"`
	Name           string            `json:"name,omitempty"`
	TriggerType    string            `json:"trigger_type,omitempty"`
	ProcessingType string            `json:"processing_type,omitempty"`
	Enabled        bool              `json:"enabled,omitempty"`
	Config         map[string]string `json:"config,omitempty"`
}

// AutomationRuleProcess represents the status of an automation rule process.
type AutomationRuleProcess struct {
	Identifier        string                      `json:"identifier,omitempty"`
	Rule              AutomationRule              `json:"rule,omitempty"`
	Status            AutomationRuleProcessStatus `json:"status,omitempty"`
	Message           *string                     `json:"message,omitempty"`
	Progress          int                         `json:"progress,omitempty"`
	CreatedAt         string                      `json:"created_at,omitempty"`
	UpdatedAt         string                      `json:"updated_at,omitempty"`
	ProcessTasks      []AutomationRuleProcessTask `json:"process_tasks,omitempty"`
	ResourceReference common.PartialResource      `json:"resource_reference,omitempty"`
}

// AutomationRuleProcessStatus represents the status of an automation rule process.
type AutomationRuleProcessStatus struct {
	StatusCode int    `json:"status_code"`
	StatusType string `json:"status_type"`
}

// AutomationRuleProcessTask represents a task within an automation rule process.
type AutomationRuleProcessTask struct {
	Identifier string                         `json:"identifier,omitempty"`
	Name       string                         `json:"name,omitempty"`
	Status     AutomationRuleProcessStatus    `json:"status,omitempty"`
	TaskInfo   *AutomationRuleProcessTaskInfo `json:"task_info,omitempty"`
}

// AutomationRuleProcessTaskInfo contains additional information about a process task.
type AutomationRuleProcessTaskInfo struct {
	Type   string      `json:"type,omitempty"`
	Config interface{} `json:"config,omitempty"`
}

// GenericAttributeState represents a state attribute in GS API responses.
type GenericAttributeState struct {
	Type  int    `json:"type,omitempty"`  // State type code from API
	ID    string `json:"id,omitempty"`    // State ID from API
	Title string `json:"title,omitempty"` // Human-readable state title from API
}

// String returns the string representation of the state value.
// Possible values: 0='OK', 1='Error', 2='Pending'
func (s *GenericAttributeState) String() string {
	if s == nil {
		return "Unknown"
	}

	// If the API provided a title, use it directly
	if s.Title != "" {
		return s.Title
	}

	// Fall back to mapping the Type code
	switch s.Type {
	case 0:
		return "OK"
	case 1:
		return "Error"
	case 2:
		return "Pending"
	default:
		return "Unknown"
	}
}

// GenericAttributeSelect represents a select attribute in GS API responses.
type GenericAttributeSelect struct {
	Identifier string `json:"identifier,omitempty"`
	Name       string `json:"name,omitempty"`
}

// OrganizationMinimal represents a minimal organization reference.
type OrganizationMinimal struct {
	Identifier string `json:"identifier,omitempty"`
	Name       string `json:"name,omitempty"`
}

// requestBody is a helper function that filters request bodies based on operation context.
func requestBody(_ interface{}, fn func() interface{}) interface{} {
	return fn()
}
