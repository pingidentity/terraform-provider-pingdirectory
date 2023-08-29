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

var _ resource.ConfigValidator = &ImpliesOtherValidatorValidator{}

// Create an ImpliesOtherValidatorValidator indicating that if the condition attribute is set to one of the given string values,
// then the implied attribute, if set, must have one of the allowed string values
func ImpliesOtherValidator(condition path.Expression, conditionValues []string, implied resource.ConfigValidator) resource.ConfigValidator {
	return ImpliesOtherValidatorValidator{
		Condition:       condition,
		ConditionValues: conditionValues,
		Implied:         implied,
	}
}

// ImpliesOtherAttributeOneOfString is the underlying struct implementing the config validator.
type ImpliesOtherValidatorValidator struct {
	Condition       path.Expression
	ConditionValues []string
	Implied         resource.ConfigValidator
}

func (v ImpliesOtherValidatorValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ImpliesOtherValidatorValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If the \"%s\" attribute is configured with one of the following values: %s, then the following validator check must pass: %s",
		v.Condition, stringSliceToReadableString(v.ConditionValues), v.Implied.MarkdownDescription(ctx))
}

func (v ImpliesOtherValidatorValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var conditionValue attr.Value

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

		for _, allowedValue := range v.ConditionValues {
			if allowedValue == conditionValueString.ValueString() {
				// Condition string value found, continue
				conditionValueMatch = true
				break
			}
		}
		if conditionValueMatch {
			break
		}
	}

	if !conditionValueMatch {
		return
	}

	// If we got here, the condition attribute was found with a condition value, so the implied validator must pass
	v.Implied.ValidateResource(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			fmt.Sprintf("\"%s\" attribute is present but implied validator check failed", v.Condition),
			v.Description(ctx),
		))
	}
}
