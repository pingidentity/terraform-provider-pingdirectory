package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Return true if this types.String represents an empty (but non-null and non-unknown) string
func IsEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() == ""
}

// Return true if this types.String represents a non-empty, non-null, non-unknown string
func IsNonEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() != ""
}

// Return true if this value represents a defined (non-null and non-unknown) value
func IsDefined(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// Check if a slice contains a value
func Contains(slice []attr.Value, value attr.Value) bool {
	for _, element := range slice {
		if element.Equal(value) {
			return true
		}
	}
	return false
}
