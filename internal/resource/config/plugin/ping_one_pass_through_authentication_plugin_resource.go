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
	_ resource.Resource                = &pingOnePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithConfigure   = &pingOnePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithImportState = &pingOnePassThroughAuthenticationPluginResource{}
	_ resource.Resource                = &defaultPingOnePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultPingOnePassThroughAuthenticationPluginResource{}
	_ resource.ResourceWithImportState = &defaultPingOnePassThroughAuthenticationPluginResource{}
)

// Create a Ping One Pass Through Authentication Plugin resource
func NewPingOnePassThroughAuthenticationPluginResource() resource.Resource {
	return &pingOnePassThroughAuthenticationPluginResource{}
}

func NewDefaultPingOnePassThroughAuthenticationPluginResource() resource.Resource {
	return &defaultPingOnePassThroughAuthenticationPluginResource{}
}

// pingOnePassThroughAuthenticationPluginResource is the resource implementation.
type pingOnePassThroughAuthenticationPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPingOnePassThroughAuthenticationPluginResource is the resource implementation.
type defaultPingOnePassThroughAuthenticationPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pingOnePassThroughAuthenticationPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ping_one_pass_through_authentication_plugin"
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ping_one_pass_through_authentication_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *pingOnePassThroughAuthenticationPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type pingOnePassThroughAuthenticationPluginResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	LastUpdated                                types.String `tfsdk:"last_updated"`
	Notifications                              types.Set    `tfsdk:"notifications"`
	RequiredActions                            types.Set    `tfsdk:"required_actions"`
	ApiURL                                     types.String `tfsdk:"api_url"`
	AuthURL                                    types.String `tfsdk:"auth_url"`
	OAuthClientID                              types.String `tfsdk:"oauth_client_id"`
	OAuthClientSecret                          types.String `tfsdk:"oauth_client_secret"`
	OAuthClientSecretPassphraseProvider        types.String `tfsdk:"oauth_client_secret_passphrase_provider"`
	EnvironmentID                              types.String `tfsdk:"environment_id"`
	IncludedLocalEntryBaseDN                   types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                         types.String `tfsdk:"connection_criteria"`
	RequestCriteria                            types.String `tfsdk:"request_criteria"`
	TryLocalBind                               types.Bool   `tfsdk:"try_local_bind"`
	OverrideLocalPassword                      types.Bool   `tfsdk:"override_local_password"`
	UpdateLocalPassword                        types.Bool   `tfsdk:"update_local_password"`
	UpdateLocalPasswordDN                      types.String `tfsdk:"update_local_password_dn"`
	AllowLaxPassThroughAuthenticationPasswords types.Bool   `tfsdk:"allow_lax_pass_through_authentication_passwords"`
	IgnoredPasswordPolicyStateErrorCondition   types.Set    `tfsdk:"ignored_password_policy_state_error_condition"`
	UserMappingLocalAttribute                  types.Set    `tfsdk:"user_mapping_local_attribute"`
	UserMappingRemoteJSONField                 types.Set    `tfsdk:"user_mapping_remote_json_field"`
	AdditionalUserMappingSCIMFilter            types.String `tfsdk:"additional_user_mapping_scim_filter"`
	Description                                types.String `tfsdk:"description"`
	Enabled                                    types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations                types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *pingOnePassThroughAuthenticationPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pingOnePassThroughAuthenticationPluginSchema(ctx, req, resp, false)
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pingOnePassThroughAuthenticationPluginSchema(ctx, req, resp, true)
}

func pingOnePassThroughAuthenticationPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ping One Pass Through Authentication Plugin.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne web service.",
				Required:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne authentication service.",
				Required:    true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Specifies the OAuth Client ID used to authenticate connections to the PingOne API.",
				Required:    true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Specifies the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
				Sensitive:   true,
			},
			"oauth_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "Specifies a passphrase provider that can be used to obtain the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Specifies the PingOne Environment that will be associated with this PingOne Pass Through Authentication Plugin.",
				Required:    true,
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the PingOne service.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A reference to request criteria that will be used to indicate which bind requests should be passed through to the PingOne service.",
				Optional:    true,
			},
			"try_local_bind": schema.BoolAttribute{
				Description: "Indicates whether to attempt the bind in the local server first, or to only send it to the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"override_local_password": schema.BoolAttribute{
				Description: "Indicates whether to attempt the authentication in the PingOne service if the local user entry includes a password. This property will only be used if try-local-bind is true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password if the local bind fails but the authentication attempt succeeds when attempted in the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password_dn": schema.StringAttribute{
				Description: "This is the DN of the user that will be used to overwrite the user's local password if update-local-password is set. The DN put here should be added to 'ignore-changes-by-dn' in the appropriate Sync Source.",
				Optional:    true,
			},
			"allow_lax_pass_through_authentication_passwords": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password even if the password used to authenticate to the PingOne service would have failed validation if the user attempted to set it directly.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ignored_password_policy_state_error_condition": schema.SetAttribute{
				Description: "A set of password policy state error conditions that should not be enforced when authentication succeeds when attempted in the PingOne service. This option can only be used if try-local-bind is true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"user_mapping_local_attribute": schema.SetAttribute{
				Description: "The names of the attributes in the local user entry whose values must match the values of the corresponding fields in the PingOne service.",
				Required:    true,
				ElementType: types.StringType,
			},
			"user_mapping_remote_json_field": schema.SetAttribute{
				Description: "The names of the fields in the PingOne service whose values must match the values of the corresponding attributes in the local user entry, as specified in the user-mapping-local-attribute property.",
				Required:    true,
				ElementType: types.StringType,
			},
			"additional_user_mapping_scim_filter": schema.StringAttribute{
				Description: "An optional SCIM filter that will be ANDed with the filter created to identify the account in the PingOne service that corresponds to the local entry. Only the \"eq\", \"sw\", \"and\", and \"or\" filter types may be used.",
				Optional:    true,
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
func addOptionalPingOnePassThroughAuthenticationPluginFields(ctx context.Context, addRequest *client.AddPingOnePassThroughAuthenticationPluginRequest, plan pingOnePassThroughAuthenticationPluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecret) {
		stringVal := plan.OAuthClientSecret.ValueString()
		addRequest.OAuthClientSecret = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecretPassphraseProvider) {
		stringVal := plan.OAuthClientSecretPassphraseProvider.ValueString()
		addRequest.OAuthClientSecretPassphraseProvider = &stringVal
	}
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
	if internaltypes.IsNonEmptyString(plan.AdditionalUserMappingSCIMFilter) {
		stringVal := plan.AdditionalUserMappingSCIMFilter.ValueString()
		addRequest.AdditionalUserMappingSCIMFilter = &stringVal
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

// Read a PingOnePassThroughAuthenticationPluginResponse object into the model struct
func readPingOnePassThroughAuthenticationPluginResponse(ctx context.Context, r *client.PingOnePassThroughAuthenticationPluginResponse, state *pingOnePassThroughAuthenticationPluginResourceModel, expectedValues *pingOnePassThroughAuthenticationPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.OAuthClientSecret = expectedValues.OAuthClientSecret
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OAuthClientSecretPassphraseProvider))
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
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
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, internaltypes.IsEmptyString(expectedValues.AdditionalUserMappingSCIMFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPingOnePassThroughAuthenticationPluginOperations(plan pingOnePassThroughAuthenticationPluginResourceModel, state pingOnePassThroughAuthenticationPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ApiURL, state.ApiURL, "api-url")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthURL, state.AuthURL, "auth-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientID, state.OAuthClientID, "oauth-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecret, state.OAuthClientSecret, "oauth-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecretPassphraseProvider, state.OAuthClientSecretPassphraseProvider, "oauth-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.EnvironmentID, state.EnvironmentID, "environment-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLocalEntryBaseDN, state.IncludedLocalEntryBaseDN, "included-local-entry-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.TryLocalBind, state.TryLocalBind, "try-local-bind")
	operations.AddBoolOperationIfNecessary(&ops, plan.OverrideLocalPassword, state.OverrideLocalPassword, "override-local-password")
	operations.AddBoolOperationIfNecessary(&ops, plan.UpdateLocalPassword, state.UpdateLocalPassword, "update-local-password")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateLocalPasswordDN, state.UpdateLocalPasswordDN, "update-local-password-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLaxPassThroughAuthenticationPasswords, state.AllowLaxPassThroughAuthenticationPasswords, "allow-lax-pass-through-authentication-passwords")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IgnoredPasswordPolicyStateErrorCondition, state.IgnoredPasswordPolicyStateErrorCondition, "ignored-password-policy-state-error-condition")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingLocalAttribute, state.UserMappingLocalAttribute, "user-mapping-local-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingRemoteJSONField, state.UserMappingRemoteJSONField, "user-mapping-remote-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalUserMappingSCIMFilter, state.AdditionalUserMappingSCIMFilter, "additional-user-mapping-scim-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *pingOnePassThroughAuthenticationPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pingOnePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var UserMappingLocalAttributeSlice []string
	plan.UserMappingLocalAttribute.ElementsAs(ctx, &UserMappingLocalAttributeSlice, false)
	var UserMappingRemoteJSONFieldSlice []string
	plan.UserMappingRemoteJSONField.ElementsAs(ctx, &UserMappingRemoteJSONFieldSlice, false)
	addRequest := client.NewAddPingOnePassThroughAuthenticationPluginRequest(plan.Id.ValueString(),
		[]client.EnumpingOnePassThroughAuthenticationPluginSchemaUrn{client.ENUMPINGONEPASSTHROUGHAUTHENTICATIONPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPING_ONE_PASS_THROUGH_AUTHENTICATION},
		plan.ApiURL.ValueString(),
		plan.AuthURL.ValueString(),
		plan.OAuthClientID.ValueString(),
		plan.EnvironmentID.ValueString(),
		UserMappingLocalAttributeSlice,
		UserMappingRemoteJSONFieldSlice,
		plan.Enabled.ValueBool())
	err := addOptionalPingOnePassThroughAuthenticationPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ping One Pass Through Authentication Plugin", err.Error())
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
		client.AddPingOnePassThroughAuthenticationPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ping One Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pingOnePassThroughAuthenticationPluginResourceModel
	readPingOnePassThroughAuthenticationPluginResponse(ctx, addResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultPingOnePassThroughAuthenticationPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pingOnePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ping One Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state pingOnePassThroughAuthenticationPluginResourceModel
	readPingOnePassThroughAuthenticationPluginResponse(ctx, readResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPingOnePassThroughAuthenticationPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ping One Pass Through Authentication Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPingOnePassThroughAuthenticationPluginResponse(ctx, updateResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pingOnePassThroughAuthenticationPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPingOnePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPingOnePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPingOnePassThroughAuthenticationPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state pingOnePassThroughAuthenticationPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ping One Pass Through Authentication Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPingOnePassThroughAuthenticationPluginResponse(ctx, readResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pingOnePassThroughAuthenticationPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePingOnePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePingOnePassThroughAuthenticationPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePingOnePassThroughAuthenticationPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan pingOnePassThroughAuthenticationPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pingOnePassThroughAuthenticationPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPingOnePassThroughAuthenticationPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ping One Pass Through Authentication Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPingOnePassThroughAuthenticationPluginResponse(ctx, updateResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPingOnePassThroughAuthenticationPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pingOnePassThroughAuthenticationPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pingOnePassThroughAuthenticationPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ping One Pass Through Authentication Plugin", err, httpResp)
		return
	}
}

func (r *pingOnePassThroughAuthenticationPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPingOnePassThroughAuthenticationPlugin(ctx, req, resp)
}

func (r *defaultPingOnePassThroughAuthenticationPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPingOnePassThroughAuthenticationPlugin(ctx, req, resp)
}

func importPingOnePassThroughAuthenticationPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
