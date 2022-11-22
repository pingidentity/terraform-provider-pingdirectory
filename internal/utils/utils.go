package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

//TODO split out utils into more fine-grained packages with more descriptive names
// Error returned from PingDirectory config API
type pingDirectoryError struct {
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
	Detail  string   `json:"detail"`
}

// Report an HTTP error
func ReportHttpError(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	httpErrorPrinted := false
	var internalError error
	if httpResp != nil {
		body, internalError := io.ReadAll(httpResp.Body)
		if internalError == nil {
			var pdError pingDirectoryError
			internalError = json.Unmarshal(body, &pdError)
			if internalError == nil {
				diagnostics.AddError(errorSummary, err.Error()+" - Detail: "+pdError.Detail)
				httpErrorPrinted = true
			}
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to unmarshal HTTP response body: "+internalError.Error())
		}
		diagnostics.AddError(errorSummary, err.Error())
	}
}

// Log operations used during an update
func LogUpdateOperations(ctx context.Context, ops []client.Operation) {
	if len(ops) == 0 {
		return
	}

	tflog.Debug(ctx, "Update using the following operations:")
	for _, op := range ops {
		opJson, err := op.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update operation: "+string(opJson))
		}
	}
}

// Get BasicAuth context with a username and password
//TODO maybe cache this somehow so it doesn't need to be done so often?
func BasicAuthContext(ctx context.Context, providerConfig ProviderConfiguration) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: providerConfig.Username,
		Password: providerConfig.Password,
	})
}

//TODO any way to reduce duplication in these methods either?
// Add boolean operation if the plan doesn't match the state
func AddBoolOperationIfNecessary(ops *[]client.Operation, plan types.Bool, state types.Bool, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}

	if !plan.Equal(state) {
		var op *client.Operation
		if plan.IsNull() {
			op = client.NewOperation(client.ENUMOPERATION_REMOVE, path)
		} else {
			op = client.NewOperation(client.ENUMOPERATION_REPLACE, path)
			op.SetValue(strconv.FormatBool(plan.ValueBool()))
		}
		*ops = append(*ops, *op)
	}
}

// Add int64 operation if the plan doesn't match the state
func AddInt64OperationIfNecessary(ops *[]client.Operation, plan types.Int64, state types.Int64, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}

	if !plan.Equal(state) {
		var op *client.Operation
		if plan.IsNull() {
			op = client.NewOperation(client.ENUMOPERATION_REMOVE, path)
		} else {
			op = client.NewOperation(client.ENUMOPERATION_REPLACE, path)
			op.SetValue(strconv.FormatInt(plan.ValueInt64(), 10))
		}
		*ops = append(*ops, *op)
	}
}

// Add string operation if the plan doesn't match the state
func AddStringOperationIfNecessary(ops *[]client.Operation, plan types.String, state types.String, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}

	if !plan.Equal(state) {
		var op *client.Operation
		// Consider an empty string as null - allows removing values despite everything being Computed
		if plan.IsNull() || plan.ValueString() == "" {
			op = client.NewOperation(client.ENUMOPERATION_REMOVE, path)
		} else {
			op = client.NewOperation(client.ENUMOPERATION_REPLACE, path)
			op.SetValue(plan.ValueString())
		}
		*ops = append(*ops, *op)
	}
}

// Add set operation if the plan doesn't match the state
func AddSetOperationsIfNecessary(ops *[]client.Operation, plan types.Set, state types.Set, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}

	if !plan.Equal(state) {
		planElements := plan.Elements()
		stateElements := state.Elements()

		// Adds
		for _, planEl := range planElements {
			if !contains(stateElements, planEl.(types.String)) {
				op := client.NewOperation(client.ENUMOPERATION_ADD, path)
				op.SetValue(planEl.(types.String).ValueString())
				*ops = append(*ops, *op)
			}
		}

		// Removes
		for _, stateEl := range stateElements {
			if !contains(planElements, stateEl.(types.String)) {
				// Remove paths for multivalued attributes are formatted like this:
				// "[additional-tags eq \"five\"]"
				op := client.NewOperation(client.ENUMOPERATION_REMOVE, "["+path+" eq \""+stateEl.(types.String).ValueString()+"\"]")
				*ops = append(*ops, *op)
			}
		}
	}
}

// Check if a slice contains a value
func contains(slice []attr.Value, value types.String) bool {
	for _, element := range slice {
		if element.(types.String).ValueString() == value.ValueString() {
			return true
		}
	}
	return false
}

// Get a types.Set from a slice of strings
func GetSet(values []string) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(values[i])
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
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
