// Package pointer provides helper functions for creating
// pointers of values and retrieving values of pointers
package pointer

// Bool converts bool value to pointer
func Bool(v bool) *bool {
	return &v
}

// BoolVal converts bool pointer to value
// with zero-value when nil
func BoolVal(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// String converts string value to pointer
func String(v string) *string {
	return &v
}

// StringVal converts string pointer to value
// with zero-value when nil
func StringVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// Int converts int value to pointer
func Int(v int) *int {
	return &v
}

// IntVal converts int pointer to value
// with zero-value when nil
func IntVal(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// UInt converts uint value to pointer
func UInt(v uint) *uint {
	return &v
}

// UIntVal converts uint pointer to value
// with zero-value when nil
func UIntVal(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}

// Float64 converts float64 value to pointer
func Float64(v float64) *float64 {
	return &v
}

// Float64Val converts float64 pointer to value
// with zero-value when nil
func Float64Val(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// Float32 converts float32 value to pointer
func Float32(v float32) *float32 {
	return &v
}

// Float32Val converts float32 pointer to value
// with zero-value when nil
func Float32Val(p *float32) float32 {
	if p == nil {
		return 0
	}
	return *p
}
