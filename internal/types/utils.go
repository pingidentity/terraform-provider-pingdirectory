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

// Check if an attribute slice contains a value
func Contains(slice []attr.Value, value attr.Value) bool {
	for _, element := range slice {
		if element.Equal(value) {
			return true
		}
	}
	return false
}

// Check if a string slice contains a value
func StringSliceContains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

// Check if two slices representing sets are equal
func SetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Assuming there are no duplicate elements since the slices represent sets
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if bElement == aElement {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ObjectsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func ObjectsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"type": types.StringType,
		},
	}
}
