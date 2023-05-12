package groupimplementation

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dynamicGroupImplementationResource{}
	_ resource.ResourceWithConfigure   = &dynamicGroupImplementationResource{}
	_ resource.ResourceWithImportState = &dynamicGroupImplementationResource{}
)

// Create a Dynamic Group Implementation resource
func NewDynamicGroupImplementationResource() resource.Resource {
	return &dynamicGroupImplementationResource{}
}

// dynamicGroupImplementationResource is the resource implementation.
type dynamicGroupImplementationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *dynamicGroupImplementationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_dynamic_group_implementation"
}

// Configure adds the provider configured client to the resource.
func (r *dynamicGroupImplementationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type dynamicGroupImplementationResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *dynamicGroupImplementationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Dynamic Group Implementation.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Group Implementation",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Group Implementation is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a DynamicGroupImplementationResponse object into the model struct
func readDynamicGroupImplementationResponse(ctx context.Context, r *client.DynamicGroupImplementationResponse, state *dynamicGroupImplementationResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDynamicGroupImplementationOperations(plan dynamicGroupImplementationResourceModel, state dynamicGroupImplementationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *dynamicGroupImplementationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dynamicGroupImplementationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationApi.GetGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dynamic Group Implementation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state dynamicGroupImplementationResourceModel
	readDynamicGroupImplementationResponse(ctx, readResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GroupImplementationApi.UpdateGroupImplementation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDynamicGroupImplementationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GroupImplementationApi.UpdateGroupImplementationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dynamic Group Implementation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDynamicGroupImplementationResponse(ctx, updateResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *dynamicGroupImplementationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dynamicGroupImplementationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationApi.GetGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dynamic Group Implementation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDynamicGroupImplementationResponse(ctx, readResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *dynamicGroupImplementationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dynamicGroupImplementationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state dynamicGroupImplementationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.GroupImplementationApi.UpdateGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDynamicGroupImplementationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GroupImplementationApi.UpdateGroupImplementationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dynamic Group Implementation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDynamicGroupImplementationResponse(ctx, updateResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *dynamicGroupImplementationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *dynamicGroupImplementationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
