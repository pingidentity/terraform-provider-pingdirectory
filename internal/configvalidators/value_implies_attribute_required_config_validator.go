// Copyright © 2025 Ping Identity Corporation

package configvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ConfigValidator = &ValueImpliesAttributeRequiredValidator{}

// Create a ValueImpliesAttributeRequiredValidator indicating that the implied attribute paths are required to be configured if the
// condition string attribute is configured with condition value
func ValueImpliesAttributeRequired(condition path.Expression, conditionValue string, implied []path.Expression) resource.ConfigValidator {
	return ValueImpliesAttributeRequiredValidator{
		Condition:      condition,
		ConditionValue: conditionValue,
		Implied:        implied,
	}
}

// ImpliesValidator is the underlying struct implementing Implies.
type ValueImpliesAttributeRequiredValidator struct {
	Condition      path.Expression
	ConditionValue string
	Implied        []path.Expression
}

func (v ValueImpliesAttributeRequiredValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ValueImpliesAttributeRequiredValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("The %s attribute(s) must be configured if the \"%s\" attribute is configured with the following value: \"%s\"", pathExpressionSliceToReadableString(v.Implied), v.Condition, v.ConditionValue)
}

func (v ValueImpliesAttributeRequiredValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var conditionValue, impliedValue attr.Value

	// Check for the condition attribute
	conditionMatchedPaths, conditionDiags := req.Config.PathMatches(ctx, v.Condition)

	resp.Diagnostics.Append(conditionDiags...)
	if conditionDiags.HasError() {
		return
	}

	conditionValueMatch := false
	for _, matchedPath := range conditionMatchedPaths {
		getAttributeDiags := req.Config.GetAttribute(ctx, matchedPath, &conditionValue)

		resp.Diagnostics.Append(getAttributeDiags...)
		if getAttributeDiags.HasError() {
			return
		}

		// If value is unknown, it may be null or a value, so we cannot
		// know if the validator should succeed or not. Collect the path
		// path so we use it to skip the validation later and continue to
		// collect all path matching diagnostics.
		if conditionValue.IsUnknown() {
			return
		}

		// If the condition is null, then try the next match
		if conditionValue.IsNull() {
			continue
		}

		// Value is known and not null, so we need to check if it is one of the condition values
		conditionValueString, ok := conditionValue.(types.String)
		if !ok {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic(
				fmt.Sprintf("\"%s\" attribute has non-string value", v.Implied),
				v.Description(ctx),
			))
			return
		}

		if v.ConditionValue == conditionValueString.ValueString() {
			// Condition string value found, continue
			conditionValueMatch = true
			break
		}
	}

	if !conditionValueMatch {
		return
	}

	// If we got here, the condition attribute was found with a condition value, so the implied attributes are required
	for _, impliedAttr := range v.Implied {
		impliedMatchedPaths, impliedDiags := req.Config.PathMatches(ctx, impliedAttr)

		resp.Diagnostics.Append(impliedDiags...)
		if impliedDiags.HasError() {
			return
		}

		valueFound := false
		for _, matchedPath := range impliedMatchedPaths {
			getAttributeDiags := req.Config.GetAttribute(ctx, matchedPath, &impliedValue)

			resp.Diagnostics.Append(getAttributeDiags...)
			if getAttributeDiags.HasError() {
				return
			}

			// If value is unknown, it may be null or a value, so we cannot
			// know if the validator should succeed or not.
			if impliedValue.IsUnknown() {
				return
			}

			// If the value is null, then try the next one
			if impliedValue.IsNull() {
				continue
			}

			// Value is known and not null, it is configured, so this attribute passes
			valueFound = true
		}

		// If we got here, then the condition value is configured and one of the implied values is not, so
		// this validator should fail
		if !valueFound {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic(
				"Missing Implied Attribute Configuration",
				v.Description(ctx),
			))
			return
		}
	}
}
