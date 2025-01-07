package delegatedadminresourcerights

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delegatedAdminResourceRightsResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminResourceRightsResource{}
	_ resource.ResourceWithImportState = &delegatedAdminResourceRightsResource{}
	_ resource.Resource                = &defaultDelegatedAdminResourceRightsResource{}
	_ resource.ResourceWithConfigure   = &defaultDelegatedAdminResourceRightsResource{}
	_ resource.ResourceWithImportState = &defaultDelegatedAdminResourceRightsResource{}
)

// Create a Delegated Admin Resource Rights resource
func NewDelegatedAdminResourceRightsResource() resource.Resource {
	return &delegatedAdminResourceRightsResource{}
}

func NewDefaultDelegatedAdminResourceRightsResource() resource.Resource {
	return &defaultDelegatedAdminResourceRightsResource{}
}

// delegatedAdminResourceRightsResource is the resource implementation.
type delegatedAdminResourceRightsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelegatedAdminResourceRightsResource is the resource implementation.
type defaultDelegatedAdminResourceRightsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminResourceRightsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_resource_rights"
}

func (r *defaultDelegatedAdminResourceRightsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delegated_admin_resource_rights"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminResourceRightsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDelegatedAdminResourceRightsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminResourceRightsResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	DelegatedAdminRightsName types.String `tfsdk:"delegated_admin_rights_name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	RestResourceType         types.String `tfsdk:"rest_resource_type"`
	AdminPermission          types.Set    `tfsdk:"admin_permission"`
	AdminScope               types.String `tfsdk:"admin_scope"`
	ResourceSubtree          types.Set    `tfsdk:"resource_subtree"`
	ResourcesInGroup         types.Set    `tfsdk:"resources_in_group"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminResourceRightsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminResourceRightsSchema(ctx, req, resp, false)
}

func (r *defaultDelegatedAdminResourceRightsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminResourceRightsSchema(ctx, req, resp, true)
}

func delegatedAdminResourceRightsSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Delegated Admin Resource Rights.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Resource Rights resource. Options are ['delegated-admin-resource-rights']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("delegated-admin-resource-rights"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"delegated-admin-resource-rights"}...),
				},
			},
			"delegated_admin_rights_name": schema.StringAttribute{
				Description: "Name of the parent Delegated Admin Rights",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Resource Rights",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether these Delegated Admin Resource Rights are enabled.",
				Required:    true,
			},
			"rest_resource_type": schema.StringAttribute{
				Description: "Specifies the resource type applicable to these Delegated Admin Resource Rights.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"admin_permission": schema.SetAttribute{
				Description: "Specifies administrator(s) permissions.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"admin_scope": schema.StringAttribute{
				Description: "Specifies the scope of these Delegated Admin Resource Rights.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("resources-in-specific-subtrees"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"resources-in-specific-groups", "resources-in-specific-subtrees", "all-resources-in-base"}...),
				},
			},
			"resource_subtree": schema.SetAttribute{
				Description: "Specifies subtrees within the search base whose entries can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-subtrees.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"resources_in_group": schema.SetAttribute{
				Description: "Specifies groups whose members can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-groups.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "rest_resource_type", "delegated_admin_rights_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		restResourceTypeAttr := schemaDef.Attributes["rest_resource_type"].(schema.StringAttribute)
		restResourceTypeAttr.PlanModifiers = append(restResourceTypeAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["rest_resource_type"] = restResourceTypeAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for delegated-admin-resource-rights delegated-admin-resource-rights
func addOptionalDelegatedAdminResourceRightsFields(ctx context.Context, addRequest *client.AddDelegatedAdminResourceRightsRequest, plan delegatedAdminResourceRightsResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AdminPermission) {
		var slice []string
		plan.AdminPermission.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdelegatedAdminResourceRightsAdminPermissionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdelegatedAdminResourceRightsAdminPermissionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AdminPermission = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdminScope) {
		adminScope, err := client.NewEnumdelegatedAdminResourceRightsAdminScopePropFromValue(plan.AdminScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.AdminScope = adminScope
	}
	if internaltypes.IsDefined(plan.ResourceSubtree) {
		var slice []string
		plan.ResourceSubtree.ElementsAs(ctx, &slice, false)
		addRequest.ResourceSubtree = slice
	}
	if internaltypes.IsDefined(plan.ResourcesInGroup) {
		var slice []string
		plan.ResourcesInGroup.ElementsAs(ctx, &slice, false)
		addRequest.ResourcesInGroup = slice
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *delegatedAdminResourceRightsResourceModel) populateAllComputedStringAttributes() {
	if model.RestResourceType.IsUnknown() || model.RestResourceType.IsNull() {
		model.RestResourceType = types.StringValue("")
	}
	if model.AdminScope.IsUnknown() || model.AdminScope.IsNull() {
		model.AdminScope = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
}

// Read a DelegatedAdminResourceRightsResponse object into the model struct
func readDelegatedAdminResourceRightsResponse(ctx context.Context, r *client.DelegatedAdminResourceRightsResponse, state *delegatedAdminResourceRightsResourceModel, expectedValues *delegatedAdminResourceRightsResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-resource-rights")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.RestResourceType = types.StringValue(r.RestResourceType)
	state.AdminPermission = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminResourceRightsAdminPermissionProp(r.AdminPermission))
	state.AdminScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminResourceRightsAdminScopeProp(r.AdminScope), true)
	state.ResourceSubtree = internaltypes.GetStringSet(r.ResourceSubtree)
	state.ResourcesInGroup = internaltypes.GetStringSet(r.ResourcesInGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *delegatedAdminResourceRightsResourceModel) setStateValuesNotReturnedByAPI(expectedValues *delegatedAdminResourceRightsResourceModel) {
	if !expectedValues.DelegatedAdminRightsName.IsUnknown() {
		state.DelegatedAdminRightsName = expectedValues.DelegatedAdminRightsName
	}
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminResourceRightsOperations(plan delegatedAdminResourceRightsResourceModel, state delegatedAdminResourceRightsResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.RestResourceType, state.RestResourceType, "rest-resource-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AdminPermission, state.AdminPermission, "admin-permission")
	operations.AddStringOperationIfNecessary(&ops, plan.AdminScope, state.AdminScope, "admin-scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResourceSubtree, state.ResourceSubtree, "resource-subtree")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResourcesInGroup, state.ResourcesInGroup, "resources-in-group")
	return ops
}

// Create a delegated-admin-resource-rights delegated-admin-resource-rights
func (r *delegatedAdminResourceRightsResource) CreateDelegatedAdminResourceRights(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminResourceRightsResourceModel) (*delegatedAdminResourceRightsResourceModel, error) {
	addRequest := client.NewAddDelegatedAdminResourceRightsRequest(plan.Enabled.ValueBool(),
		plan.RestResourceType.ValueString())
	err := addOptionalDelegatedAdminResourceRightsFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Delegated Admin Resource Rights", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminResourceRightsAPI.AddDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.DelegatedAdminRightsName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminResourceRightsRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.DelegatedAdminResourceRightsAPI.AddDelegatedAdminResourceRightsExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Resource Rights", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminResourceRightsResourceModel
	readDelegatedAdminResourceRightsResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *delegatedAdminResourceRightsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminResourceRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateDelegatedAdminResourceRights(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultDelegatedAdminResourceRightsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminResourceRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminResourceRightsAPI.GetDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceType.ValueString(), plan.DelegatedAdminRightsName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Resource Rights", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminResourceRightsResourceModel
	readDelegatedAdminResourceRightsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminResourceRightsAPI.UpdateDelegatedAdminResourceRights(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceType.ValueString(), plan.DelegatedAdminRightsName.ValueString())
	ops := createDelegatedAdminResourceRightsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminResourceRightsAPI.UpdateDelegatedAdminResourceRightsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Resource Rights", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminResourceRightsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *delegatedAdminResourceRightsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminResourceRights(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDelegatedAdminResourceRightsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminResourceRights(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDelegatedAdminResourceRights(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state delegatedAdminResourceRightsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminResourceRightsAPI.GetDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.RestResourceType.ValueString(), state.DelegatedAdminRightsName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Resource Rights", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Resource Rights", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminResourceRightsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *delegatedAdminResourceRightsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminResourceRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminResourceRightsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminResourceRights(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelegatedAdminResourceRights(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delegatedAdminResourceRightsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminResourceRightsResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminResourceRightsAPI.UpdateDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.RestResourceType.ValueString(), plan.DelegatedAdminRightsName.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminResourceRightsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminResourceRightsAPI.UpdateDelegatedAdminResourceRightsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Resource Rights", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminResourceRightsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultDelegatedAdminResourceRightsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminResourceRightsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delegatedAdminResourceRightsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminResourceRightsAPI.DeleteDelegatedAdminResourceRightsExecute(r.apiClient.DelegatedAdminResourceRightsAPI.DeleteDelegatedAdminResourceRights(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.RestResourceType.ValueString(), state.DelegatedAdminRightsName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delegated Admin Resource Rights", err, httpResp)
		return
	}
}

func (r *delegatedAdminResourceRightsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminResourceRights(ctx, req, resp)
}

func (r *defaultDelegatedAdminResourceRightsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminResourceRights(ctx, req, resp)
}

func importDelegatedAdminResourceRights(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [delegated-admin-rights-name]/[delegated-admin-resource-rights-rest-resource-type]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("delegated_admin_rights_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_resource_type"), split[1])...)
}
