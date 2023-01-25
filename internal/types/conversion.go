package types

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Convert a int64 to string
func Int64ToString(value types.Int64) string {
	return strconv.FormatInt(value.ValueInt64(), 10)
}

// Get a types.Set from a slice of strings
func GetStringSet(values []string) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(values[i])
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

// Get a types.Set from a slice of int32
func GetInt64Set(values []int32) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.Int64Value(int64(values[i]))
	}
	set, _ := types.SetValue(types.Int64Type, setValues)
	return set
}

// Convert Enums from the Go client to a types.Set.
// Maintain alphabetical order for ease of management
func GetEnumSet[V client.EnumglobalConfigurationAllowedInsecureTLSProtocolProp |
	client.EnumglobalConfigurationAttributesModifiableWithIgnoreNoUserModificationRequestControlProp |
	client.EnumglobalConfigurationDisabledPrivilegeProp |
	client.EnumglobalConfigurationInvalidAttributeSyntaxBehaviorProp |
	client.EnumglobalConfigurationJmxValueBehaviorProp |
	client.EnumglobalConfigurationSingleStructuralObjectclassBehaviorProp |
	client.EnumglobalConfigurationStartupErrorLoggerOutputLocationProp |
	client.EnumglobalConfigurationWritabilityModeProp |
	client.EnumglobalConfigurationUnrecoverableDatabaseErrorModeProp |
	client.EnumrootDnDefaultRootPrivilegeNameProp](values []V) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(string(values[i]))
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

// convert enum slice to string slice
func GetEnumStringSlice[V client.EnumglobalConfigurationAllowedInsecureTLSProtocolProp |
	client.EnumglobalConfigurationAttributesModifiableWithIgnoreNoUserModificationRequestControlProp |
	client.EnumglobalConfigurationDisabledPrivilegeProp |
	client.EnumglobalConfigurationInvalidAttributeSyntaxBehaviorProp |
	client.EnumglobalConfigurationJmxValueBehaviorProp |
	client.EnumglobalConfigurationSingleStructuralObjectclassBehaviorProp |
	client.EnumglobalConfigurationStartupErrorLoggerOutputLocationProp |
	client.EnumglobalConfigurationWritabilityModeProp |
	client.EnumglobalConfigurationUnrecoverableDatabaseErrorModeProp |
	client.EnumrootDnDefaultRootPrivilegeNameProp](values []V) []string {
	stringSlice := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		stringSlice[i] = string(values[i])
	}
	return stringSlice
}

// Get a types.String from the given string pointer, handling if the pointer is nil
func StringTypeOrNil(str *string, useEmptyStringForNil bool) types.String {
	if str == nil {
		// If a plan was provided and is using an empty string, we should use that for a nil string in the response.
		// To PingDirectory, nil and empty string is equivalent, but to Terraform they are distinct. So we
		// just want to match whatever is in the plan when we get a nil string back.
		if useEmptyStringForNil {
			// Use empty string instead of null to match the plan when resetting string properties.
			// This is useful for computed values being reset to null.
			return types.StringValue("")
		} else {
			return types.StringNull()
		}
	}
	return types.StringValue(*str)
}

// Get a types.Bool from the given bool pointer, handling if the pointer is nil
func BoolTypeOrNil(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// Get a types.Int64 from the given int32 pointer, handling if the pointer is nil
func Int64TypeOrNil(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

// Get a types.String from the given Stringer, handling if the pointer is nil
func StringerStringTypeOrNil(value fmt.Stringer) types.String {
	// Similar to the StringTypeOrNil method - use an empty string for nil values
	if value != nil {
		return types.StringValue(value.String())
	} else {
		return types.StringValue("")
	}
}
