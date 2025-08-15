// Copyright © 2025 Ping Identity Corporation

package accesscontrolhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accessControlHandlerResource{}
	_ resource.ResourceWithConfigure   = &accessControlHandlerResource{}
	_ resource.ResourceWithImportState = &accessControlHandlerResource{}
)

// Create a Access Control Handler resource
func NewAccessControlHandlerResource() resource.Resource {
	return &accessControlHandlerResource{}
}

// accessControlHandlerResource is the resource implementation.
type accessControlHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *accessControlHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_access_control_handler"
}

// Configure adds the provider configured client to the resource.
func (r *accessControlHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type accessControlHandlerResourceModel struct {
	Id                                            types.String `tfsdk:"id"`
	Notifications                                 types.Set    `tfsdk:"notifications"`
	RequiredActions                               types.Set    `tfsdk:"required_actions"`
	Type                                          types.String `tfsdk:"type"`
	GlobalACI                                     types.Set    `tfsdk:"global_aci"`
	AllowedBindControl                            types.Set    `tfsdk:"allowed_bind_control"`
	AllowedBindControlOID                         types.Set    `tfsdk:"allowed_bind_control_oid"`
	EvaluateTargetAttributeRightsForAddOperations types.Bool   `tfsdk:"evaluate_target_attribute_rights_for_add_operations"`
	Enabled                                       types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *accessControlHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Access Control Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Access Control Handler resource. Options are ['dsee-compat']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"dsee-compat"}...),
				},
			},
			"global_aci": schema.SetAttribute{
				Description: "Defines global access control rules.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_bind_control": schema.SetAttribute{
				Description: "Specifies a set of controls that clients should be allowed to include in bind requests. As bind requests are evaluated as the unauthenticated user, any controls included in this set will be permitted for any bind attempt. If you wish to grant permission for any bind controls not listed here, then the allowed-bind-control-oid property may be used to accomplish that.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_bind_control_oid": schema.SetAttribute{
				Description: "Specifies the OIDs of any additional controls (not covered by the allowed-bind-control property) that should be permitted in bind requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluate_target_attribute_rights_for_add_operations": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.1.0.0+. Indicates whether the server should ensure that the requester has the \"add\" right for each attribute included in an add request, and is not denied \"add\" rights for any attributes in the request. Historically, any user who has been granted the \"add\" right has been allowed to create an entry of any type, even for add requests that include attributes for which they do not have the \"add\" right (that is, the \"targetattr\" portion of an access control rule was not considered when evaluating access control rights for add operations). This is still the default behavior in order to preserve backward compatibility, but setting the value of this property to true will cause the server to only permit add operations in which the requester has the \"add\" right for each of the attributes included in the add request, and deny add operations if the requester is denied \"add\" rights for any attributes included in the add request. It is strongly recommended that you thoroughly test your existing access control configuration before enabling this setting in a production environment to identify any cases in which you may need to add or augment access control rules to ensure that authorized users are allowed to add the entries they need to be able to create.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *accessControlHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory10100)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model accessControlHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.EvaluateTargetAttributeRightsForAddOperations) {
		resp.Diagnostics.AddError("Attribute 'evaluate_target_attribute_rights_for_add_operations' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a DseeCompatAccessControlHandlerResponse object into the model struct
func readDseeCompatAccessControlHandlerResponse(ctx context.Context, r *client.DseeCompatAccessControlHandlerResponse, state *accessControlHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dsee-compat")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.GlobalACI = internaltypes.GetStringSet(r.GlobalACI)
	state.AllowedBindControl = internaltypes.GetStringSet(
		client.StringSliceEnumaccessControlHandlerAllowedBindControlProp(r.AllowedBindControl))
	state.AllowedBindControlOID = internaltypes.GetStringSet(r.AllowedBindControlOID)
	state.EvaluateTargetAttributeRightsForAddOperations = internaltypes.BoolTypeOrNil(r.EvaluateTargetAttributeRightsForAddOperations)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAccessControlHandlerOperations(plan accessControlHandlerResourceModel, state accessControlHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GlobalACI, state.GlobalACI, "global-aci")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedBindControl, state.AllowedBindControl, "allowed-bind-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedBindControlOID, state.AllowedBindControlOID, "allowed-bind-control-oid")
	operations.AddBoolOperationIfNecessary(&ops, plan.EvaluateTargetAttributeRightsForAddOperations, state.EvaluateTargetAttributeRightsForAddOperations, "evaluate-target-attribute-rights-for-add-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *accessControlHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accessControlHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessControlHandlerAPI.GetAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Control Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state accessControlHandlerResourceModel
	readDseeCompatAccessControlHandlerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccessControlHandlerAPI.UpdateAccessControlHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createAccessControlHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessControlHandlerAPI.UpdateAccessControlHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Access Control Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDseeCompatAccessControlHandlerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *accessControlHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state accessControlHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessControlHandlerAPI.GetAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Control Handler", err, httpResp)
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
}

// Update a resource
func (r *accessControlHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan accessControlHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state accessControlHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AccessControlHandlerAPI.UpdateAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createAccessControlHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessControlHandlerAPI.UpdateAccessControlHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Access Control Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDseeCompatAccessControlHandlerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *accessControlHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *accessControlHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
