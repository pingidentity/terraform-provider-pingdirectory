package azureauthenticationmethod

import (
	"context"
	"time"

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
	_ resource.Resource                = &azureAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &azureAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &azureAuthenticationMethodResource{}
	_ resource.Resource                = &defaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultAzureAuthenticationMethodResource{}
)

// Create a Azure Authentication Method resource
func NewAzureAuthenticationMethodResource() resource.Resource {
	return &azureAuthenticationMethodResource{}
}

func NewDefaultAzureAuthenticationMethodResource() resource.Resource {
	return &defaultAzureAuthenticationMethodResource{}
}

// azureAuthenticationMethodResource is the resource implementation.
type azureAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAzureAuthenticationMethodResource is the resource implementation.
type defaultAzureAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *azureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_authentication_method"
}

func (r *defaultAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_azure_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *azureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type azureAuthenticationMethodResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	TenantID        types.String `tfsdk:"tenant_id"`
	ClientID        types.String `tfsdk:"client_id"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *azureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	azureAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	azureAuthenticationMethodSchema(ctx, req, resp, true)
}

func azureAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Azure Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Azure Authentication Method resource. Options are ['default', 'client-secret', 'username-password']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"default", "client-secret", "username-password"}...),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`client-secret`, `username-password`]: The tenant ID to use to authenticate. When the `type` attribute is set to `default`: The tenant ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_TENANT_ID environment variable.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`client-secret`, `username-password`]: The tenant ID to use to authenticate.\n  - `default`: The tenant ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_TENANT_ID environment variable.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`client-secret`, `username-password`]: The client ID to use to authenticate. When the `type` attribute is set to `default`: The client ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_CLIENT_ID",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`client-secret`, `username-password`]: The client ID to use to authenticate.\n  - `default`: The client ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_CLIENT_ID",
				Optional:            true,
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
			"client_secret": schema.StringAttribute{
				Description: "The client secret to use to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Azure Authentication Method",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsAzureAuthenticationMethod() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password"),
			path.MatchRoot("type"),
			[]string{"username-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_secret"),
			path.MatchRoot("type"),
			[]string{"client-secret"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("username"),
			path.MatchRoot("type"),
			[]string{"username-password"},
		),
	}
}

// Add config validators
func (r azureAuthenticationMethodResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsAzureAuthenticationMethod()
}

// Add config validators
func (r defaultAzureAuthenticationMethodResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsAzureAuthenticationMethod()
}

// Add optional fields to create request for default azure-authentication-method
func addOptionalDefaultAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddDefaultAzureAuthenticationMethodRequest, plan azureAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TenantID) {
		addRequest.TenantID = plan.TenantID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientID) {
		addRequest.ClientID = plan.ClientID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for client-secret azure-authentication-method
func addOptionalClientSecretAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddClientSecretAzureAuthenticationMethodRequest, plan azureAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for username-password azure-authentication-method
func addOptionalUsernamePasswordAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddUsernamePasswordAzureAuthenticationMethodRequest, plan azureAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAzureAuthenticationMethodUnknownValues(ctx context.Context, model *azureAuthenticationMethodResourceModel) {
	if model.ClientSecret.IsUnknown() {
		model.ClientSecret = types.StringNull()
	}
	if model.Password.IsUnknown() {
		model.Password = types.StringNull()
	}
}

// Read a DefaultAzureAuthenticationMethodResponse object into the model struct
func readDefaultAzureAuthenticationMethodResponse(ctx context.Context, r *client.DefaultAzureAuthenticationMethodResponse, state *azureAuthenticationMethodResourceModel, expectedValues *azureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = internaltypes.StringTypeOrNil(r.TenantID, internaltypes.IsEmptyString(expectedValues.TenantID))
	state.ClientID = internaltypes.StringTypeOrNil(r.ClientID, internaltypes.IsEmptyString(expectedValues.ClientID))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAzureAuthenticationMethodUnknownValues(ctx, state)
}

// Read a ClientSecretAzureAuthenticationMethodResponse object into the model struct
func readClientSecretAzureAuthenticationMethodResponse(ctx context.Context, r *client.ClientSecretAzureAuthenticationMethodResponse, state *azureAuthenticationMethodResourceModel, expectedValues *azureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-secret")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = types.StringValue(r.TenantID)
	state.ClientID = types.StringValue(r.ClientID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAzureAuthenticationMethodUnknownValues(ctx, state)
}

// Read a UsernamePasswordAzureAuthenticationMethodResponse object into the model struct
func readUsernamePasswordAzureAuthenticationMethodResponse(ctx context.Context, r *client.UsernamePasswordAzureAuthenticationMethodResponse, state *azureAuthenticationMethodResourceModel, expectedValues *azureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("username-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TenantID = types.StringValue(r.TenantID)
	state.ClientID = types.StringValue(r.ClientID)
	state.Username = types.StringValue(r.Username)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAzureAuthenticationMethodUnknownValues(ctx, state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *azureAuthenticationMethodResourceModel) setStateValuesNotReturnedByAPI(expectedValues *azureAuthenticationMethodResourceModel) {
	if !expectedValues.ClientSecret.IsUnknown() {
		state.ClientSecret = expectedValues.ClientSecret
	}
	if !expectedValues.Password.IsUnknown() {
		state.Password = expectedValues.Password
	}
}

// Create any update operations necessary to make the state match the plan
func createAzureAuthenticationMethodOperations(plan azureAuthenticationMethodResourceModel, state azureAuthenticationMethodResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.TenantID, state.TenantID, "tenant-id")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientSecret, state.ClientSecret, "client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a default azure-authentication-method
func (r *azureAuthenticationMethodResource) CreateDefaultAzureAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan azureAuthenticationMethodResourceModel) (*azureAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddDefaultAzureAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumdefaultAzureAuthenticationMethodSchemaUrn{client.ENUMDEFAULTAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODDEFAULT})
	addOptionalDefaultAzureAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
		client.AddDefaultAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Azure Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state azureAuthenticationMethodResourceModel
	readDefaultAzureAuthenticationMethodResponse(ctx, addResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a client-secret azure-authentication-method
func (r *azureAuthenticationMethodResource) CreateClientSecretAzureAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan azureAuthenticationMethodResourceModel) (*azureAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddClientSecretAzureAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumclientSecretAzureAuthenticationMethodSchemaUrn{client.ENUMCLIENTSECRETAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODCLIENT_SECRET},
		plan.TenantID.ValueString(),
		plan.ClientID.ValueString(),
		plan.ClientSecret.ValueString())
	addOptionalClientSecretAzureAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
		client.AddClientSecretAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Azure Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state azureAuthenticationMethodResourceModel
	readClientSecretAzureAuthenticationMethodResponse(ctx, addResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a username-password azure-authentication-method
func (r *azureAuthenticationMethodResource) CreateUsernamePasswordAzureAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan azureAuthenticationMethodResourceModel) (*azureAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddUsernamePasswordAzureAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumusernamePasswordAzureAuthenticationMethodSchemaUrn{client.ENUMUSERNAMEPASSWORDAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODUSERNAME_PASSWORD},
		plan.TenantID.ValueString(),
		plan.ClientID.ValueString(),
		plan.Username.ValueString(),
		plan.Password.ValueString())
	addOptionalUsernamePasswordAzureAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
		client.AddUsernamePasswordAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Azure Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state azureAuthenticationMethodResourceModel
	readUsernamePasswordAzureAuthenticationMethodResponse(ctx, addResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *azureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan azureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *azureAuthenticationMethodResourceModel
	var err error
	if plan.Type.ValueString() == "default" {
		state, err = r.CreateDefaultAzureAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "client-secret" {
		state, err = r.CreateClientSecretAzureAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "username-password" {
		state, err = r.CreateUsernamePasswordAzureAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan azureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state azureAuthenticationMethodResourceModel
	if readResponse.DefaultAzureAuthenticationMethodResponse != nil {
		readDefaultAzureAuthenticationMethodResponse(ctx, readResponse.DefaultAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClientSecretAzureAuthenticationMethodResponse != nil {
		readClientSecretAzureAuthenticationMethodResponse(ctx, readResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UsernamePasswordAzureAuthenticationMethodResponse != nil {
		readUsernamePasswordAzureAuthenticationMethodResponse(ctx, readResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createAzureAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Azure Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultAzureAuthenticationMethodResponse != nil {
			readDefaultAzureAuthenticationMethodResponse(ctx, updateResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClientSecretAzureAuthenticationMethodResponse != nil {
			readClientSecretAzureAuthenticationMethodResponse(ctx, updateResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UsernamePasswordAzureAuthenticationMethodResponse != nil {
			readUsernamePasswordAzureAuthenticationMethodResponse(ctx, updateResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *azureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readAzureAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state azureAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Authentication Method", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Authentication Method", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultAzureAuthenticationMethodResponse != nil {
		readDefaultAzureAuthenticationMethodResponse(ctx, readResponse.DefaultAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClientSecretAzureAuthenticationMethodResponse != nil {
		readClientSecretAzureAuthenticationMethodResponse(ctx, readResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UsernamePasswordAzureAuthenticationMethodResponse != nil {
		readUsernamePasswordAzureAuthenticationMethodResponse(ctx, readResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *azureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAzureAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan azureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state azureAuthenticationMethodResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createAzureAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Azure Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultAzureAuthenticationMethodResponse != nil {
			readDefaultAzureAuthenticationMethodResponse(ctx, updateResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClientSecretAzureAuthenticationMethodResponse != nil {
			readClientSecretAzureAuthenticationMethodResponse(ctx, updateResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UsernamePasswordAzureAuthenticationMethodResponse != nil {
			readUsernamePasswordAzureAuthenticationMethodResponse(ctx, updateResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *defaultAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *azureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state azureAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethodExecute(r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Azure Authentication Method", err, httpResp)
		return
	}
}

func (r *azureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAzureAuthenticationMethod(ctx, req, resp)
}

func (r *defaultAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAzureAuthenticationMethod(ctx, req, resp)
}

func importAzureAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
