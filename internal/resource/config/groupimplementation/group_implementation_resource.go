package groupimplementation

import (
	"context"
	"time"

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
	_ resource.Resource                = &groupImplementationResource{}
	_ resource.ResourceWithConfigure   = &groupImplementationResource{}
	_ resource.ResourceWithImportState = &groupImplementationResource{}
)

// Create a Group Implementation resource
func NewGroupImplementationResource() resource.Resource {
	return &groupImplementationResource{}
}

// groupImplementationResource is the resource implementation.
type groupImplementationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groupImplementationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_group_implementation"
}

// Configure adds the provider configured client to the resource.
func (r *groupImplementationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type groupImplementationResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *groupImplementationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Group Implementation.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Group Implementation resource. Options are ['static', 'virtual-static', 'dynamic']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"static", "virtual-static", "dynamic"}...),
				},
			},
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
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a StaticGroupImplementationResponse object into the model struct
func readStaticGroupImplementationResponse(ctx context.Context, r *client.StaticGroupImplementationResponse, state *groupImplementationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("static")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a VirtualStaticGroupImplementationResponse object into the model struct
func readVirtualStaticGroupImplementationResponse(ctx context.Context, r *client.VirtualStaticGroupImplementationResponse, state *groupImplementationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("virtual-static")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a DynamicGroupImplementationResponse object into the model struct
func readDynamicGroupImplementationResponse(ctx context.Context, r *client.DynamicGroupImplementationResponse, state *groupImplementationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dynamic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroupImplementationOperations(plan groupImplementationResourceModel, state groupImplementationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *groupImplementationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupImplementationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationApi.GetGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Group Implementation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groupImplementationResourceModel
	if readResponse.StaticGroupImplementationResponse != nil {
		readStaticGroupImplementationResponse(ctx, readResponse.StaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VirtualStaticGroupImplementationResponse != nil {
		readVirtualStaticGroupImplementationResponse(ctx, readResponse.VirtualStaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DynamicGroupImplementationResponse != nil {
		readDynamicGroupImplementationResponse(ctx, readResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GroupImplementationApi.UpdateGroupImplementation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createGroupImplementationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GroupImplementationApi.UpdateGroupImplementationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Group Implementation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StaticGroupImplementationResponse != nil {
			readStaticGroupImplementationResponse(ctx, updateResponse.StaticGroupImplementationResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.VirtualStaticGroupImplementationResponse != nil {
			readVirtualStaticGroupImplementationResponse(ctx, updateResponse.VirtualStaticGroupImplementationResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DynamicGroupImplementationResponse != nil {
			readDynamicGroupImplementationResponse(ctx, updateResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
		}
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
func (r *groupImplementationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state groupImplementationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationApi.GetGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Group Implementation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StaticGroupImplementationResponse != nil {
		readStaticGroupImplementationResponse(ctx, readResponse.StaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VirtualStaticGroupImplementationResponse != nil {
		readVirtualStaticGroupImplementationResponse(ctx, readResponse.VirtualStaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DynamicGroupImplementationResponse != nil {
		readDynamicGroupImplementationResponse(ctx, readResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *groupImplementationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan groupImplementationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groupImplementationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.GroupImplementationApi.UpdateGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createGroupImplementationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GroupImplementationApi.UpdateGroupImplementationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Group Implementation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StaticGroupImplementationResponse != nil {
			readStaticGroupImplementationResponse(ctx, updateResponse.StaticGroupImplementationResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.VirtualStaticGroupImplementationResponse != nil {
			readVirtualStaticGroupImplementationResponse(ctx, updateResponse.VirtualStaticGroupImplementationResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DynamicGroupImplementationResponse != nil {
			readDynamicGroupImplementationResponse(ctx, updateResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
		}
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
func (r *groupImplementationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groupImplementationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
