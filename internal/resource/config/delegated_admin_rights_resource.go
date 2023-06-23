package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delegatedAdminRightsResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminRightsResource{}
	_ resource.ResourceWithImportState = &delegatedAdminRightsResource{}
	_ resource.Resource                = &defaultDelegatedAdminRightsResource{}
	_ resource.ResourceWithConfigure   = &defaultDelegatedAdminRightsResource{}
	_ resource.ResourceWithImportState = &defaultDelegatedAdminRightsResource{}
)

// Create a Delegated Admin Rights resource
func NewDelegatedAdminRightsResource() resource.Resource {
	return &delegatedAdminRightsResource{}
}

func NewDefaultDelegatedAdminRightsResource() resource.Resource {
	return &defaultDelegatedAdminRightsResource{}
}

// delegatedAdminRightsResource is the resource implementation.
type delegatedAdminRightsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelegatedAdminRightsResource is the resource implementation.
type defaultDelegatedAdminRightsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminRightsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_rights"
}

func (r *defaultDelegatedAdminRightsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delegated_admin_rights"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminRightsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDelegatedAdminRightsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type delegatedAdminRightsResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	AdminUserDN     types.String `tfsdk:"admin_user_dn"`
	AdminGroupDN    types.String `tfsdk:"admin_group_dn"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminRightsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminRightsSchema(ctx, req, resp, false)
}

func (r *defaultDelegatedAdminRightsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminRightsSchema(ctx, req, resp, true)
}

func delegatedAdminRightsSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Delegated Admin Rights.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Rights",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Delegated Admin Rights is enabled.",
				Required:    true,
			},
			"admin_user_dn": schema.StringAttribute{
				Description: "Specifies the DN of an administrative user who has authority to manage resources. Either admin-user-dn or admin-group-dn must be specified, but not both.",
				Optional:    true,
			},
			"admin_group_dn": schema.StringAttribute{
				Description: "Specifies the DN of a group of administrative users who have authority to manage resources. Either admin-user-dn or admin-group-dn must be specified, but not both.",
				Optional:    true,
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators
func (r delegatedAdminRightsResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("admin_group_dn"),
			path.MatchRoot("admin_user_dn"),
		),
	}
}

// Add optional fields to create request for delegated-admin-rights delegated-admin-rights
func addOptionalDelegatedAdminRightsFields(ctx context.Context, addRequest *client.AddDelegatedAdminRightsRequest, plan delegatedAdminRightsResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdminUserDN) {
		addRequest.AdminUserDN = plan.AdminUserDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdminGroupDN) {
		addRequest.AdminGroupDN = plan.AdminGroupDN.ValueStringPointer()
	}
}

// Read a DelegatedAdminRightsResponse object into the model struct
func readDelegatedAdminRightsResponse(ctx context.Context, r *client.DelegatedAdminRightsResponse, state *delegatedAdminRightsResourceModel, expectedValues *delegatedAdminRightsResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AdminUserDN = internaltypes.StringTypeOrNil(r.AdminUserDN, internaltypes.IsEmptyString(expectedValues.AdminUserDN))
	state.AdminGroupDN = internaltypes.StringTypeOrNil(r.AdminGroupDN, internaltypes.IsEmptyString(expectedValues.AdminGroupDN))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminRightsOperations(plan delegatedAdminRightsResourceModel, state delegatedAdminRightsResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AdminUserDN, state.AdminUserDN, "admin-user-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.AdminGroupDN, state.AdminGroupDN, "admin-group-dn")
	return ops
}

// Create a delegated-admin-rights delegated-admin-rights
func (r *delegatedAdminRightsResource) CreateDelegatedAdminRights(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminRightsResourceModel) (*delegatedAdminRightsResourceModel, error) {
	addRequest := client.NewAddDelegatedAdminRightsRequest(plan.Id.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalDelegatedAdminRightsFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminRightsApi.AddDelegatedAdminRights(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDelegatedAdminRightsRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.DelegatedAdminRightsApi.AddDelegatedAdminRightsExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Rights", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminRightsResourceModel
	readDelegatedAdminRightsResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *delegatedAdminRightsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateDelegatedAdminRights(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultDelegatedAdminRightsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminRightsApi.GetDelegatedAdminRights(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Rights", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminRightsResourceModel
	readDelegatedAdminRightsResponse(ctx, readResponse, &state, &plan, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminRightsApi.UpdateDelegatedAdminRights(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDelegatedAdminRightsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminRightsApi.UpdateDelegatedAdminRightsExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Rights", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminRightsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *delegatedAdminRightsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminRightsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDelegatedAdminRights(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state delegatedAdminRightsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminRightsApi.GetDelegatedAdminRights(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Rights", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminRightsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *delegatedAdminRightsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminRightsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelegatedAdminRights(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delegatedAdminRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminRightsResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminRightsApi.UpdateDelegatedAdminRights(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminRightsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminRightsApi.UpdateDelegatedAdminRightsExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Rights", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminRightsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDelegatedAdminRightsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminRightsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delegatedAdminRightsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminRightsApi.DeleteDelegatedAdminRightsExecute(r.apiClient.DelegatedAdminRightsApi.DeleteDelegatedAdminRights(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delegated Admin Rights", err, httpResp)
		return
	}
}

func (r *delegatedAdminRightsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminRights(ctx, req, resp)
}

func (r *defaultDelegatedAdminRightsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminRights(ctx, req, resp)
}

func importDelegatedAdminRights(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
