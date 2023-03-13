package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pluggablePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithConfigure   = &pluggablePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithImportState = &pluggablePassThroughAuthenticationPluginResource{}
	_ resource.Resource                = &defaultPluggablePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultPluggablePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithImportState = &defaultPluggablePassThroughAuthenticationPluginResource{}
)

// Create a Pluggable Pass Through Authentication Plugin resource
func NewPluggablePassThroughAuthenticationPluginResource() resource.Resource {
	return &pluggablePassThroughAuthenticationPluginResource{}
}

func NewDefaultPluggablePassThroughAuthenticationPluginResource() resource.Resource {
	return &defaultPluggablePassThroughAuthenticationPluginResource{}
}

// pluggablePassThroughAuthenticationPluginResource is the resource implementation.
type pluggablePassThroughAuthenticationPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPluggablePassThroughAuthenticationPluginResource is the resource implementation.
type defaultPluggablePassThroughAuthenticationPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pluggablePassThroughAuthenticationPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pluggable_pass_through_authentication_plugin"
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_pluggable_pass_through_authentication_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *pluggablePassThroughAuthenticationPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type pluggablePassThroughAuthenticationPluginResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	LastUpdated                                types.String `tfsdk:"last_updated"`
	Notifications                              types.Set    `tfsdk:"notifications"`
	RequiredActions                            types.Set    `tfsdk:"required_actions"`
	PassThroughAuthenticationHandler           types.String `tfsdk:"pass_through_authentication_handler"`
	IncludedLocalEntryBaseDN                   types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                         types.String `tfsdk:"connection_criteria"`
	RequestCriteria                            types.String `tfsdk:"request_criteria"`
	TryLocalBind                               types.Bool   `tfsdk:"try_local_bind"`
	OverrideLocalPassword                      types.Bool   `tfsdk:"override_local_password"`
	UpdateLocalPassword                        types.Bool   `tfsdk:"update_local_password"`
	UpdateLocalPasswordDN                      types.String `tfsdk:"update_local_password_dn"`
	AllowLaxPassThroughAuthenticationPasswords types.Bool   `tfsdk:"allow_lax_pass_through_authentication_passwords"`
	IgnoredPasswordPolicyStateErrorCondition   types.Set    `tfsdk:"ignored_password_policy_state_error_condition"`
	Description                                types.String `tfsdk:"description"`
	Enabled                                    types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations                types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *pluggablePassThroughAuthenticationPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pluggablePassThroughAuthenticationPluginSchema(ctx, req, resp, false)
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pluggablePassThroughAuthenticationPluginSchema(ctx, req, resp, true)
}

func pluggablePassThroughAuthenticationPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Pluggable Pass Through Authentication Plugin.",
		Attributes: map[string]schema.Attribute{
			"pass_through_authentication_handler": schema.StringAttribute{
				Description: "The component used to manage authentication with the external authentication service.",
				Required:    true,
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the external authentication service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A reference to request criteria that will be used to indicate which bind requests should be passed through to the external authentication service.",
				Optional:    true,
			},
			"try_local_bind": schema.BoolAttribute{
				Description: "Indicates whether to attempt the bind in the local server first and only send the request to the external authentication service if the local bind attempt fails, or to only attempt the bind in the external service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"override_local_password": schema.BoolAttribute{
				Description: "Indicates whether to attempt the authentication in the external service if the local user entry includes a password. This property will be ignored if try-local-bind is false.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password if the local bind fails but the authentication attempt succeeds when attempted in the external service. This property may only be set to true if try-local-bind is also true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password_dn": schema.StringAttribute{
				Description: "The DN of the authorization identity that will be used when updating the user's local password if update-local-password is true. This is primarily intended for use if the Data Sync Server will be used to synchronize passwords between the local server and the external service, and in that case, the DN used here should also be added to the ignore-changes-by-dn property in the appropriate Sync Source object in the Data Sync Server configuration.",
				Optional:    true,
			},
			"allow_lax_pass_through_authentication_passwords": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password even if the password used to authenticate to the external service would have failed validation if the user attempted to set it directly.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ignored_password_policy_state_error_condition": schema.SetAttribute{
				Description: "A set of password policy state error conditions that should not be enforced when authentication succeeds when attempted in the external service. This option can only be used if try-local-bind is true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalPluggablePassThroughAuthenticationPluginFields(ctx context.Context, addRequest *client.AddPluggablePassThroughAuthenticationPluginRequest, plan pluggablePassThroughAuthenticationPluginResourceModel) error {
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		stringVal := plan.ConnectionCriteria.ValueString()
		addRequest.ConnectionCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		stringVal := plan.RequestCriteria.ValueString()
		addRequest.RequestCriteria = &stringVal
	}
	if internaltypes.IsDefined(plan.TryLocalBind) {
		boolVal := plan.TryLocalBind.ValueBool()
		addRequest.TryLocalBind = &boolVal
	}
	if internaltypes.IsDefined(plan.OverrideLocalPassword) {
		boolVal := plan.OverrideLocalPassword.ValueBool()
		addRequest.OverrideLocalPassword = &boolVal
	}
	if internaltypes.IsDefined(plan.UpdateLocalPassword) {
		boolVal := plan.UpdateLocalPassword.ValueBool()
		addRequest.UpdateLocalPassword = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateLocalPasswordDN) {
		stringVal := plan.UpdateLocalPasswordDN.ValueString()
		addRequest.UpdateLocalPasswordDN = &stringVal
	}
	if internaltypes.IsDefined(plan.AllowLaxPassThroughAuthenticationPasswords) {
		boolVal := plan.AllowLaxPassThroughAuthenticationPasswords.ValueBool()
		addRequest.AllowLaxPassThroughAuthenticationPasswords = &boolVal
	}
	if internaltypes.IsDefined(plan.IgnoredPasswordPolicyStateErrorCondition) {
		var slice []string
		plan.IgnoredPasswordPolicyStateErrorCondition.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIgnoredPasswordPolicyStateErrorConditionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIgnoredPasswordPolicyStateErrorConditionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IgnoredPasswordPolicyStateErrorCondition = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		boolVal := plan.InvokeForInternalOperations.ValueBool()
		addRequest.InvokeForInternalOperations = &boolVal
	}
	return nil
}

// Read a PluggablePassThroughAuthenticationPluginResponse object into the model struct
func readPluggablePassThroughAuthenticationPluginResponse(ctx context.Context, r *client.PluggablePassThroughAuthenticationPluginResponse, state *pluggablePassThroughAuthenticationPluginResourceModel, expectedValues *pluggablePassThroughAuthenticationPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PassThroughAuthenticationHandler = types.StringValue(r.PassThroughAuthenticationHandler)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, internaltypes.IsEmptyString(expectedValues.UpdateLocalPasswordDN))
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPluggablePassThroughAuthenticationPluginOperations(plan pluggablePassThroughAuthenticationPluginResourceModel, state pluggablePassThroughAuthenticationPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PassThroughAuthenticationHandler, state.PassThroughAuthenticationHandler, "pass-through-authentication-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLocalEntryBaseDN, state.IncludedLocalEntryBaseDN, "included-local-entry-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.TryLocalBind, state.TryLocalBind, "try-local-bind")
	operations.AddBoolOperationIfNecessary(&ops, plan.OverrideLocalPassword, state.OverrideLocalPassword, "override-local-password")
	operations.AddBoolOperationIfNecessary(&ops, plan.UpdateLocalPassword, state.UpdateLocalPassword, "update-local-password")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateLocalPasswordDN, state.UpdateLocalPasswordDN, "update-local-password-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLaxPassThroughAuthenticationPasswords, state.AllowLaxPassThroughAuthenticationPasswords, "allow-lax-pass-through-authentication-passwords")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IgnoredPasswordPolicyStateErrorCondition, state.IgnoredPasswordPolicyStateErrorCondition, "ignored-password-policy-state-error-condition")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *pluggablePassThroughAuthenticationPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pluggablePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPluggablePassThroughAuthenticationPluginRequest(plan.Id.ValueString(),
		[]client.EnumpluggablePassThroughAuthenticationPluginSchemaUrn{client.ENUMPLUGGABLEPASSTHROUGHAUTHENTICATIONPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPLUGGABLE_PASS_THROUGH_AUTHENTICATION},
		plan.PassThroughAuthenticationHandler.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPluggablePassThroughAuthenticationPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pluggable Pass Through Authentication Plugin", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPluggablePassThroughAuthenticationPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pluggable Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluggablePassThroughAuthenticationPluginResourceModel
	readPluggablePassThroughAuthenticationPluginResponse(ctx, addResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultPluggablePassThroughAuthenticationPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pluggablePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pluggable Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state pluggablePassThroughAuthenticationPluginResourceModel
	readPluggablePassThroughAuthenticationPluginResponse(ctx, readResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPluggablePassThroughAuthenticationPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pluggable Pass Through Authentication Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPluggablePassThroughAuthenticationPluginResponse(ctx, updateResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pluggablePassThroughAuthenticationPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPluggablePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPluggablePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPluggablePassThroughAuthenticationPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state pluggablePassThroughAuthenticationPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pluggable Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPluggablePassThroughAuthenticationPluginResponse(ctx, readResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pluggablePassThroughAuthenticationPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePluggablePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePluggablePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePluggablePassThroughAuthenticationPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan pluggablePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pluggablePassThroughAuthenticationPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPluggablePassThroughAuthenticationPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pluggable Pass Through Authentication Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPluggablePassThroughAuthenticationPluginResponse(ctx, updateResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPluggablePassThroughAuthenticationPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pluggablePassThroughAuthenticationPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pluggablePassThroughAuthenticationPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Pluggable Pass Through Authentication Plugin", err, httpResp)
		return
	}
}

func (r *pluggablePassThroughAuthenticationPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPluggablePassThroughAuthenticationPlugin(ctx, req, resp)
}

func (r *defaultPluggablePassThroughAuthenticationPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPluggablePassThroughAuthenticationPlugin(ctx, req, resp)
}

func importPluggablePassThroughAuthenticationPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
