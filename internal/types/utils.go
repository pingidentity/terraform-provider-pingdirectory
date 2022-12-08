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

// Return true if this types.Bool represents a defined (non-null and non-unknown) boolean
func IsDefinedBool(b types.Bool) bool {
	return !b.IsNull() && !b.IsUnknown()
}

// Return true if this types.Set represents a defined (non-null and non-unknown) set
func IsDefinedSet(s types.Set) bool {
	return !s.IsNull() && !s.IsUnknown()
}

// Check if a slice contains a string value
func ContainsString(slice []attr.Value, value types.String) bool {
	for _, element := range slice {
		if element.(types.String).ValueString() == value.ValueString() {
			return true
		}
	}
	return false
}

// Check if a slice contains an int64 value
func ContainsInt64(slice []attr.Value, value types.Int64) bool {
	for _, element := range slice {
		if element.(types.Int64).ValueInt64() == value.ValueInt64() {
			return true
		}
	}
	return false
}
