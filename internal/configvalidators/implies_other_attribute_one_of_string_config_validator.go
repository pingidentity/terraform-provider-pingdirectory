package configvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ConfigValidator = &ImpliesOtherAttributeOneOfStringValidator{}

// Create an ImpliesOtherAttributeOneOfString indicating that if the condition attribute is set,
// then the implied attribute, if set, must have one of the allowed string values
func ImpliesOtherAttributeOneOfString(condition, implied path.Expression, impliedAllowedValues []string) resource.ConfigValidator {
	return ImpliesOtherAttributeOneOfStringValidator{
		Condition:            condition,
		Implied:              implied,
		ImpliedAllowedValues: impliedAllowedValues,
	}
}

// ImpliesOtherAttributeOneOfString is the underlying struct implementing the config validator.
type ImpliesOtherAttributeOneOfStringValidator struct {
	Condition            path.Expression
	Implied              path.Expression
	ImpliedAllowedValues []string
}

func (v ImpliesOtherAttributeOneOfStringValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ImpliesOtherAttributeOneOfStringValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If the \"%s\" attribute is configured, then the \"%s\" attribute must have one of the following values if it is configured: %s",
		v.Condition, v.Implied, stringSliceToReadableString(v.ImpliedAllowedValues))
}

func (v ImpliesOtherAttributeOneOfStringValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ImpliesOtherAttributeOneOfStringValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var diags diag.Diagnostics
	var conditionValue, impliedValue attr.Value

	// Check for the condition attribute
	conditionMatchedPaths, conditionDiags := config.PathMatches(ctx, v.Condition)

	diags.Append(conditionDiags...)
	if conditionDiags.HasError() {
		return diags
	}

	conditionFound := false
	for _, matchedPath := range conditionMatchedPaths {
		getAttributeDiags := config.GetAttribute(ctx, matchedPath, &conditionValue)

		diags.Append(getAttributeDiags...)
		if getAttributeDiags.HasError() {
			return diags
		}

		// If value is unknown, it may be null or a value, so we cannot
		// know if the validator should succeed or not. Collect the path
		// path so we use it to skip the validation later and continue to
		// collect all path matching diagnostics.
		if conditionValue.IsUnknown() {
			return diags
		}

		// If the condition is null, then try the next match
		if conditionValue.IsNull() {
			continue
		}

		// Value is known and not null, it is configured.
		conditionFound = true
		break
	}

	// If the condition isn't configured, then this validator doesn't need to do anything
	if !conditionFound {
		return diags
	}

	// If we got here, the condition attribute was found, so the implied attribute must be present
	impliedMatchedPaths, impliedDiags := config.PathMatches(ctx, v.Implied)

	diags.Append(impliedDiags...)
	if impliedDiags.HasError() {
		return diags
	}

	for _, matchedPath := range impliedMatchedPaths {
		getAttributeDiags := config.GetAttribute(ctx, matchedPath, &impliedValue)

		diags.Append(getAttributeDiags...)
		if getAttributeDiags.HasError() {
			return diags
		}

		// If value is unknown, it may be null or a value, so we cannot
		// know if the validator should succeed or not. Collect the path
		// path so we use it to skip the validation later and continue to
		// collect all path matching diagnostics.
		if impliedValue.IsUnknown() {
			return diags
		}

		// If the condition is null, then try the next one
		if impliedValue.IsNull() {
			continue
		}

		// Value is known and not null, so we need to check if it is one of the allowed values
		impliedValueString, ok := impliedValue.(types.String)
		if !ok {
			diags.Append(diag.NewErrorDiagnostic(
				fmt.Sprintf("\"%s\" attribute has non-string value", v.Implied),
				v.Description(ctx),
			))
			return diags
		}

		allowedValueFound := false
		for _, allowedValue := range v.ImpliedAllowedValues {
			if allowedValue == impliedValueString.ValueString() {
				// Valid string value found, continue
				allowedValueFound = true
				break
			}
		}
		if allowedValueFound {
			continue
		}

		// If we reach here, then the value is not allowed
		diags.Append(diag.NewErrorDiagnostic(
			fmt.Sprintf("\"%s\" attribute is configured but \"%s\" attribute does not have an allowed value", v.Condition, v.Implied),
			v.Description(ctx),
		))
		return diags
	}

	return diags
}
