package pingdirectory

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// https://developer.hashicorp.com/terraform/plugin/framework/resources/plan-modification

// booleanDefaultModifier is a plan modifier that sets a default value for a
// types.BooleanType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type booleanDefaultModifier struct {
	Default bool
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m booleanDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m booleanDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m booleanDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// If the value is unknown or known, do not set default value.
	if !req.AttributePlan.IsNull() {
		return
	}

	// types.Bool must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var b types.Bool
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &b)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.AttributePlan = types.BoolValue(m.Default)
}

func booleanDefault(defaultValue bool) booleanDefaultModifier {
	return booleanDefaultModifier{
		Default: defaultValue,
	}
}
