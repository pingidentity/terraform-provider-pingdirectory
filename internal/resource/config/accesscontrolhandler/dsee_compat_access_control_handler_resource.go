package accesscontrolhandler

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dseeCompatAccessControlHandlerResource{}
	_ resource.ResourceWithConfigure   = &dseeCompatAccessControlHandlerResource{}
	_ resource.ResourceWithImportState = &dseeCompatAccessControlHandlerResource{}
)

// Create a Dsee Compat Access Control Handler resource
func NewDseeCompatAccessControlHandlerResource() resource.Resource {
	return &dseeCompatAccessControlHandlerResource{}
}

// dseeCompatAccessControlHandlerResource is the resource implementation.
type dseeCompatAccessControlHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *dseeCompatAccessControlHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_dsee_compat_access_control_handler"
}

// Configure adds the provider configured client to the resource.
func (r *dseeCompatAccessControlHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type dseeCompatAccessControlHandlerResourceModel struct {
	// Id field required for acceptance testing framework
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	GlobalACI             types.Set    `tfsdk:"global_aci"`
	AllowedBindControl    types.Set    `tfsdk:"allowed_bind_control"`
	AllowedBindControlOID types.Set    `tfsdk:"allowed_bind_control_oid"`
	Enabled               types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *dseeCompatAccessControlHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Dsee Compat Access Control Handler.",
		Attributes: map[string]schema.Attribute{
			"global_aci": schema.SetAttribute{
				Description: "Defines global access control rules.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"allowed_bind_control": schema.SetAttribute{
				Description: "Specifies a set of controls that clients should be allowed to include in bind requests. As bind requests are evaluated as the unauthenticated user, any controls included in this set will be permitted for any bind attempt. If you wish to grant permission for any bind controls not listed here, then the allowed-bind-control-oid property may be used to accomplish that.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"allowed_bind_control_oid": schema.SetAttribute{
				Description: "Specifies the OIDs of any additional controls (not covered by the allowed-bind-control property) that should be permitted in bind requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Control Handler is enabled. If set to FALSE, then no access control is enforced, and any client (including unauthenticated or anonymous clients) could be allowed to perform any operation if not subject to other restrictions, such as those enforced by the privilege subsystem.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Read a DseeCompatAccessControlHandlerResponse object into the model struct
func readDseeCompatAccessControlHandlerResponse(ctx context.Context, r *client.DseeCompatAccessControlHandlerResponse, state *dseeCompatAccessControlHandlerResourceModel, diagnostics *diag.Diagnostics) {
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.GlobalACI = internaltypes.GetStringSet(r.GlobalACI)
	state.AllowedBindControl = internaltypes.GetStringSet(
		client.StringSliceEnumaccessControlHandlerAllowedBindControlProp(r.AllowedBindControl))
	state.AllowedBindControlOID = internaltypes.GetStringSet(r.AllowedBindControlOID)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDseeCompatAccessControlHandlerOperations(plan dseeCompatAccessControlHandlerResourceModel, state dseeCompatAccessControlHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GlobalACI, state.GlobalACI, "global-aci")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedBindControl, state.AllowedBindControl, "allowed-bind-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedBindControlOID, state.AllowedBindControlOID, "allowed-bind-control-oid")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *dseeCompatAccessControlHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dseeCompatAccessControlHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessControlHandlerApi.GetAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dsee Compat Access Control Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state dseeCompatAccessControlHandlerResourceModel
	readDseeCompatAccessControlHandlerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccessControlHandlerApi.UpdateAccessControlHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createDseeCompatAccessControlHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessControlHandlerApi.UpdateAccessControlHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dsee Compat Access Control Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDseeCompatAccessControlHandlerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *dseeCompatAccessControlHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dseeCompatAccessControlHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessControlHandlerApi.GetAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dsee Compat Access Control Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDseeCompatAccessControlHandlerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *dseeCompatAccessControlHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dseeCompatAccessControlHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state dseeCompatAccessControlHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AccessControlHandlerApi.UpdateAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createDseeCompatAccessControlHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessControlHandlerApi.UpdateAccessControlHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dsee Compat Access Control Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDseeCompatAccessControlHandlerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *dseeCompatAccessControlHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *dseeCompatAccessControlHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
