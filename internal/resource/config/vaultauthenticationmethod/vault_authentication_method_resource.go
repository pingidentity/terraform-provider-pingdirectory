package vaultauthenticationmethod

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vaultAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &vaultAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &vaultAuthenticationMethodResource{}
	_ resource.Resource                = &defaultVaultAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultVaultAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultVaultAuthenticationMethodResource{}
)

// Create a Vault Authentication Method resource
func NewVaultAuthenticationMethodResource() resource.Resource {
	return &vaultAuthenticationMethodResource{}
}

func NewDefaultVaultAuthenticationMethodResource() resource.Resource {
	return &defaultVaultAuthenticationMethodResource{}
}

// vaultAuthenticationMethodResource is the resource implementation.
type vaultAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVaultAuthenticationMethodResource is the resource implementation.
type defaultVaultAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *vaultAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_authentication_method"
}

func (r *defaultVaultAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_vault_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *vaultAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultVaultAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type vaultAuthenticationMethodResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	Type               types.String `tfsdk:"type"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	VaultRoleID        types.String `tfsdk:"vault_role_id"`
	VaultSecretID      types.String `tfsdk:"vault_secret_id"`
	LoginMechanismName types.String `tfsdk:"login_mechanism_name"`
	VaultAccessToken   types.String `tfsdk:"vault_access_token"`
	Description        types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *vaultAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultVaultAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultAuthenticationMethodSchema(ctx, req, resp, true)
}

func vaultAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Vault Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Vault Authentication Method resource. Options are ['static-token', 'app-role', 'user-pass']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"static-token", "app-role", "user-pass"}...),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username for the user to authenticate.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"vault_role_id": schema.StringAttribute{
				Description: "The role ID for the AppRole to authenticate.",
				Optional:    true,
			},
			"vault_secret_id": schema.StringAttribute{
				Description: "The secret ID for the AppRole to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"login_mechanism_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `app-role`: The name used when enabling the desired AppRole authentication mechanism in the Vault server. When the `type` attribute is set to `user-pass`: The name used when enabling the desired UserPass authentication mechanism in the Vault server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `app-role`: The name used when enabling the desired AppRole authentication mechanism in the Vault server.\n  - `user-pass`: The name used when enabling the desired UserPass authentication mechanism in the Vault server.",
				Optional:            true,
				Computed:            true,
			},
			"vault_access_token": schema.StringAttribute{
				Description: "The static token used to authenticate to the Vault server.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Vault Authentication Method",
				Optional:    true,
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *vaultAuthenticationMethodResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel vaultAuthenticationMethodResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for user-pass type
	if resourceType == "user-pass" {
		if !internaltypes.IsDefined(configModel.LoginMechanismName) {
			defaultVal := types.StringValue("userpass")
			if !planModel.LoginMechanismName.Equal(defaultVal) {
				planModel.LoginMechanismName = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	resp.Plan.Set(ctx, &planModel)
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsVaultAuthenticationMethod() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_access_token"),
			path.MatchRoot("type"),
			[]string{"static-token"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_role_id"),
			path.MatchRoot("type"),
			[]string{"app-role"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_secret_id"),
			path.MatchRoot("type"),
			[]string{"app-role"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("login_mechanism_name"),
			path.MatchRoot("type"),
			[]string{"app-role", "user-pass"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("username"),
			path.MatchRoot("type"),
			[]string{"user-pass"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password"),
			path.MatchRoot("type"),
			[]string{"user-pass"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"static-token",
			[]path.Expression{path.MatchRoot("vault_access_token")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"app-role",
			[]path.Expression{path.MatchRoot("vault_role_id"), path.MatchRoot("vault_secret_id")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"user-pass",
			[]path.Expression{path.MatchRoot("username"), path.MatchRoot("password")},
		),
	}
}

// Add config validators
func (r vaultAuthenticationMethodResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsVaultAuthenticationMethod()
}

// Add config validators
func (r defaultVaultAuthenticationMethodResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsVaultAuthenticationMethod()
}

// Add optional fields to create request for static-token vault-authentication-method
func addOptionalStaticTokenVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddStaticTokenVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for app-role vault-authentication-method
func addOptionalAppRoleVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddAppRoleVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoginMechanismName) {
		addRequest.LoginMechanismName = plan.LoginMechanismName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for user-pass vault-authentication-method
func addOptionalUserPassVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddUserPassVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoginMechanismName) {
		addRequest.LoginMechanismName = plan.LoginMechanismName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateVaultAuthenticationMethodUnknownValues(model *vaultAuthenticationMethodResourceModel) {
	if model.LoginMechanismName.IsUnknown() || model.LoginMechanismName.IsNull() {
		model.LoginMechanismName = types.StringValue("")
	}
	if model.VaultAccessToken.IsUnknown() {
		model.VaultAccessToken = types.StringNull()
	}
	if model.VaultSecretID.IsUnknown() {
		model.VaultSecretID = types.StringNull()
	}
	if model.Password.IsUnknown() {
		model.Password = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *vaultAuthenticationMethodResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.Username.IsUnknown() || model.Username.IsNull() {
		model.Username = types.StringValue("")
	}
	if model.VaultRoleID.IsUnknown() || model.VaultRoleID.IsNull() {
		model.VaultRoleID = types.StringValue("")
	}
}

// Read a StaticTokenVaultAuthenticationMethodResponse object into the model struct
func readStaticTokenVaultAuthenticationMethodResponse(ctx context.Context, r *client.StaticTokenVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("static-token")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(state)
}

// Read a AppRoleVaultAuthenticationMethodResponse object into the model struct
func readAppRoleVaultAuthenticationMethodResponse(ctx context.Context, r *client.AppRoleVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("app-role")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultRoleID = types.StringValue(r.VaultRoleID)
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(state)
}

// Read a UserPassVaultAuthenticationMethodResponse object into the model struct
func readUserPassVaultAuthenticationMethodResponse(ctx context.Context, r *client.UserPassVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-pass")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Username = types.StringValue(r.Username)
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *vaultAuthenticationMethodResourceModel) setStateValuesNotReturnedByAPI(expectedValues *vaultAuthenticationMethodResourceModel) {
	if !expectedValues.VaultAccessToken.IsUnknown() {
		state.VaultAccessToken = expectedValues.VaultAccessToken
	}
	if !expectedValues.VaultSecretID.IsUnknown() {
		state.VaultSecretID = expectedValues.VaultSecretID
	}
	if !expectedValues.Password.IsUnknown() {
		state.Password = expectedValues.Password
	}
}

// Create any update operations necessary to make the state match the plan
func createVaultAuthenticationMethodOperations(plan vaultAuthenticationMethodResourceModel, state vaultAuthenticationMethodResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultRoleID, state.VaultRoleID, "vault-role-id")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretID, state.VaultSecretID, "vault-secret-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LoginMechanismName, state.LoginMechanismName, "login-mechanism-name")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultAccessToken, state.VaultAccessToken, "vault-access-token")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a static-token vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateStaticTokenVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddStaticTokenVaultAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumstaticTokenVaultAuthenticationMethodSchemaUrn{client.ENUMSTATICTOKENVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODSTATIC_TOKEN},
		plan.VaultAccessToken.ValueString())
	addOptionalStaticTokenVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddStaticTokenVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readStaticTokenVaultAuthenticationMethodResponse(ctx, addResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a app-role vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateAppRoleVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddAppRoleVaultAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumappRoleVaultAuthenticationMethodSchemaUrn{client.ENUMAPPROLEVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODAPP_ROLE},
		plan.VaultRoleID.ValueString(),
		plan.VaultSecretID.ValueString())
	addOptionalAppRoleVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddAppRoleVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readAppRoleVaultAuthenticationMethodResponse(ctx, addResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a user-pass vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateUserPassVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddUserPassVaultAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumuserPassVaultAuthenticationMethodSchemaUrn{client.ENUMUSERPASSVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODUSER_PASS},
		plan.Username.ValueString(),
		plan.Password.ValueString())
	addOptionalUserPassVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddUserPassVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readUserPassVaultAuthenticationMethodResponse(ctx, addResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *vaultAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *vaultAuthenticationMethodResourceModel
	var err error
	if plan.Type.ValueString() == "static-token" {
		state, err = r.CreateStaticTokenVaultAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "app-role" {
		state, err = r.CreateAppRoleVaultAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "user-pass" {
		state, err = r.CreateUserPassVaultAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultVaultAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state vaultAuthenticationMethodResourceModel
	if readResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
		readStaticTokenVaultAuthenticationMethodResponse(ctx, readResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AppRoleVaultAuthenticationMethodResponse != nil {
		readAppRoleVaultAuthenticationMethodResponse(ctx, readResponse.AppRoleVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UserPassVaultAuthenticationMethodResponse != nil {
		readUserPassVaultAuthenticationMethodResponse(ctx, readResponse.UserPassVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createVaultAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
			readStaticTokenVaultAuthenticationMethodResponse(ctx, updateResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AppRoleVaultAuthenticationMethodResponse != nil {
			readAppRoleVaultAuthenticationMethodResponse(ctx, updateResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UserPassVaultAuthenticationMethodResponse != nil {
			readUserPassVaultAuthenticationMethodResponse(ctx, updateResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *vaultAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultVaultAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readVaultAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state vaultAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
		readStaticTokenVaultAuthenticationMethodResponse(ctx, readResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AppRoleVaultAuthenticationMethodResponse != nil {
		readAppRoleVaultAuthenticationMethodResponse(ctx, readResponse.AppRoleVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UserPassVaultAuthenticationMethodResponse != nil {
		readUserPassVaultAuthenticationMethodResponse(ctx, readResponse.UserPassVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *vaultAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVaultAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state vaultAuthenticationMethodResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createVaultAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
			readStaticTokenVaultAuthenticationMethodResponse(ctx, updateResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AppRoleVaultAuthenticationMethodResponse != nil {
			readAppRoleVaultAuthenticationMethodResponse(ctx, updateResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UserPassVaultAuthenticationMethodResponse != nil {
			readUserPassVaultAuthenticationMethodResponse(ctx, updateResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultVaultAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *vaultAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vaultAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VaultAuthenticationMethodApi.DeleteVaultAuthenticationMethodExecute(r.apiClient.VaultAuthenticationMethodApi.DeleteVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Vault Authentication Method", err, httpResp)
		return
	}
}

func (r *vaultAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultAuthenticationMethod(ctx, req, resp)
}

func (r *defaultVaultAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultAuthenticationMethod(ctx, req, resp)
}

func importVaultAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
