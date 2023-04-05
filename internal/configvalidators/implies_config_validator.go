package configvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ resource.ConfigValidator = &ImpliesValidator{}

// Create an ImpliesValidator indicating that the condition path being configured implies the implied path is configured
func Implies(condition path.Expression, implied path.Expression) resource.ConfigValidator {
	return ImpliesValidator{
		Condition: condition,
		Implied:   implied,
	}
}

// ImpliesValidator is the underlying struct implementing Implies.
type ImpliesValidator struct {
	Condition path.Expression
	Implied   path.Expression
}

func (v ImpliesValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ImpliesValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If the \"%s\" attribute is configured, then the \"%s\" attribute must be configured", v.Condition, v.Implied)
}

func (v ImpliesValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ImpliesValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
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

		// Value is known and not null, it is configured, so this validator passes
		return diags
	}

	// If we got here, then the condition value is configured and the implied value is not, so
	// this validator should fail
	diags.Append(diag.NewErrorDiagnostic(
		"Missing Implied Attribute Configuration",
		v.Description(ctx),
	))

	return diags
}
