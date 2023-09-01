package matchingrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &matchingRuleResource{}
	_ resource.ResourceWithConfigure   = &matchingRuleResource{}
	_ resource.ResourceWithImportState = &matchingRuleResource{}
)

// Create a Matching Rule resource
func NewMatchingRuleResource() resource.Resource {
	return &matchingRuleResource{}
}

// matchingRuleResource is the resource implementation.
type matchingRuleResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *matchingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_matching_rule"
}

// Configure adds the provider configured client to the resource.
func (r *matchingRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type matchingRuleResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *matchingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Matching Rule.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Matching Rule resource. Options are ['ordering', 'approximate', 'equality', 'substring', 'generic']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ordering", "approximate", "equality", "substring", "generic"}...),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Matching Rule is enabled for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a OrderingMatchingRuleResponse object into the model struct
func readOrderingMatchingRuleResponse(ctx context.Context, r *client.OrderingMatchingRuleResponse, state *matchingRuleResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ordering")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a ApproximateMatchingRuleResponse object into the model struct
func readApproximateMatchingRuleResponse(ctx context.Context, r *client.ApproximateMatchingRuleResponse, state *matchingRuleResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("approximate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a EqualityMatchingRuleResponse object into the model struct
func readEqualityMatchingRuleResponse(ctx context.Context, r *client.EqualityMatchingRuleResponse, state *matchingRuleResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("equality")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a SubstringMatchingRuleResponse object into the model struct
func readSubstringMatchingRuleResponse(ctx context.Context, r *client.SubstringMatchingRuleResponse, state *matchingRuleResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("substring")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a GenericMatchingRuleResponse object into the model struct
func readGenericMatchingRuleResponse(ctx context.Context, r *client.GenericMatchingRuleResponse, state *matchingRuleResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createMatchingRuleOperations(plan matchingRuleResourceModel, state matchingRuleResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *matchingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan matchingRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MatchingRuleApi.GetMatchingRule(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Matching Rule", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state matchingRuleResourceModel
	if readResponse.OrderingMatchingRuleResponse != nil {
		readOrderingMatchingRuleResponse(ctx, readResponse.OrderingMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ApproximateMatchingRuleResponse != nil {
		readApproximateMatchingRuleResponse(ctx, readResponse.ApproximateMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityMatchingRuleResponse != nil {
		readEqualityMatchingRuleResponse(ctx, readResponse.EqualityMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubstringMatchingRuleResponse != nil {
		readSubstringMatchingRuleResponse(ctx, readResponse.SubstringMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericMatchingRuleResponse != nil {
		readGenericMatchingRuleResponse(ctx, readResponse.GenericMatchingRuleResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.MatchingRuleApi.UpdateMatchingRule(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createMatchingRuleOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MatchingRuleApi.UpdateMatchingRuleExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Matching Rule", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.OrderingMatchingRuleResponse != nil {
			readOrderingMatchingRuleResponse(ctx, updateResponse.OrderingMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.ApproximateMatchingRuleResponse != nil {
			readApproximateMatchingRuleResponse(ctx, updateResponse.ApproximateMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.EqualityMatchingRuleResponse != nil {
			readEqualityMatchingRuleResponse(ctx, updateResponse.EqualityMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.SubstringMatchingRuleResponse != nil {
			readSubstringMatchingRuleResponse(ctx, updateResponse.SubstringMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericMatchingRuleResponse != nil {
			readGenericMatchingRuleResponse(ctx, updateResponse.GenericMatchingRuleResponse, &state, &resp.Diagnostics)
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *matchingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state matchingRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MatchingRuleApi.GetMatchingRule(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Matching Rule", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.OrderingMatchingRuleResponse != nil {
		readOrderingMatchingRuleResponse(ctx, readResponse.OrderingMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ApproximateMatchingRuleResponse != nil {
		readApproximateMatchingRuleResponse(ctx, readResponse.ApproximateMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityMatchingRuleResponse != nil {
		readEqualityMatchingRuleResponse(ctx, readResponse.EqualityMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubstringMatchingRuleResponse != nil {
		readSubstringMatchingRuleResponse(ctx, readResponse.SubstringMatchingRuleResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericMatchingRuleResponse != nil {
		readGenericMatchingRuleResponse(ctx, readResponse.GenericMatchingRuleResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *matchingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan matchingRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state matchingRuleResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.MatchingRuleApi.UpdateMatchingRule(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createMatchingRuleOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MatchingRuleApi.UpdateMatchingRuleExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Matching Rule", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.OrderingMatchingRuleResponse != nil {
			readOrderingMatchingRuleResponse(ctx, updateResponse.OrderingMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.ApproximateMatchingRuleResponse != nil {
			readApproximateMatchingRuleResponse(ctx, updateResponse.ApproximateMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.EqualityMatchingRuleResponse != nil {
			readEqualityMatchingRuleResponse(ctx, updateResponse.EqualityMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.SubstringMatchingRuleResponse != nil {
			readSubstringMatchingRuleResponse(ctx, updateResponse.SubstringMatchingRuleResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericMatchingRuleResponse != nil {
			readGenericMatchingRuleResponse(ctx, updateResponse.GenericMatchingRuleResponse, &state, &resp.Diagnostics)
		}
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *matchingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *matchingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
