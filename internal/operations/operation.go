package operations

import (
	"context"
	"strconv"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

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

// Validate that the path for a given operation is valid
func validateOperationPath(path string) {
	// Paths must only contain lowercase letters, dashes and digits
	for _, c := range path {
		if !unicode.IsLower(c) && c != '-' && !unicode.IsDigit(c) {
			//lintignore:R009
			panic("Non-lowercase, non-dash character and non-digit '" + string(c) + "' included in Operation path: '" + path + "'")
		}
	}
}

// Add boolean operation if the plan doesn't match the state
func AddBoolOperationIfNecessary(ops *[]client.Operation, plan types.Bool, state types.Bool, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}
	validateOperationPath(path)

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
	validateOperationPath(path)

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

// Add float64 operation if the plan doesn't match the state
func AddFloat64OperationIfNecessary(ops *[]client.Operation, plan types.Float64, state types.Float64, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}
	validateOperationPath(path)

	if !plan.Equal(state) {
		var op *client.Operation
		if plan.IsNull() {
			op = client.NewOperation(client.ENUMOPERATION_REMOVE, path)
		} else {
			op = client.NewOperation(client.ENUMOPERATION_REPLACE, path)
			op.SetValue(strconv.FormatFloat(plan.ValueFloat64(), 'f', -1, 64))
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
	validateOperationPath(path)

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

// Get a path to remove a value from a multi-valued attribute
func removeMultiValuedAttributePath(attributePath string, toRemove string) string {
	// Remove paths for multivalued attributes are formatted like this:
	// "[additional-tags eq \"five\"]"
	return "[" + attributePath + " eq \"" + toRemove + "\"]"
}

// Add set operation if the plan doesn't match the state
func AddStringSetOperationsIfNecessary(ops *[]client.Operation, plan types.Set, state types.Set, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}
	validateOperationPath(path)

	if !plan.Equal(state) {
		planElements := plan.Elements()
		stateElements := state.Elements()

		// Adds
		for _, planEl := range planElements {
			if !internaltypes.Contains(stateElements, planEl) {
				op := client.NewOperation(client.ENUMOPERATION_ADD, path)
				op.SetValue(planEl.(types.String).ValueString())
				*ops = append(*ops, *op)
			}
		}

		// Removes
		for _, stateEl := range stateElements {
			if !internaltypes.Contains(planElements, stateEl) {
				op := client.NewOperation(client.ENUMOPERATION_REMOVE, removeMultiValuedAttributePath(path, stateEl.(types.String).ValueString()))
				*ops = append(*ops, *op)
			}
		}
	}
}

// Add int64 set operation if the plan doesn't match the state
func AddInt64SetOperationsIfNecessary(ops *[]client.Operation, plan types.Set, state types.Set, path string) {
	// If plan is unknown, then just take whatever's in the state - no operation needed
	if plan.IsUnknown() {
		return
	}
	validateOperationPath(path)

	if !plan.Equal(state) {
		planElements := plan.Elements()
		stateElements := state.Elements()

		// Adds
		for _, planEl := range planElements {
			if !internaltypes.Contains(stateElements, planEl.(types.Int64)) {
				op := client.NewOperation(client.ENUMOPERATION_ADD, path)
				op.SetValue(internaltypes.Int64ToString(planEl.(types.Int64)))
				*ops = append(*ops, *op)
			}
		}

		// Removes
		for _, stateEl := range stateElements {
			if !internaltypes.Contains(planElements, stateEl.(types.Int64)) {
				op := client.NewOperation(client.ENUMOPERATION_REMOVE, removeMultiValuedAttributePath(path, internaltypes.Int64ToString(stateEl.(types.Int64))))
				*ops = append(*ops, *op)
			}
		}
	}
}
