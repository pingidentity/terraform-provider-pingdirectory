package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

var _ planmodifier.String = toLowercasePlanModifier{}

// toLowercasePlanModifier implements the plan modifier.
type toLowercasePlanModifier struct{}

// Helper method to create the plan modifier.
func ToLowercasePlanModifier() planmodifier.String {
	return toLowercasePlanModifier{}
}

// Description returns a human-readable description of the plan modifier.
func (m toLowercasePlanModifier) Description(_ context.Context) string {
	return "Forces this string to lowercase in the plan."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m toLowercasePlanModifier) MarkdownDescription(_ context.Context) string {
	return "Forces this string to lowercase in the plan."
}

// PlanModifyBool implements the plan modification logic.
func (m toLowercasePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no defined value
	if !internaltypes.IsDefined(req.PlanValue) {
		return
	}

	resp.PlanValue = types.StringValue(strings.ToLower(req.PlanValue.ValueString()))
}
