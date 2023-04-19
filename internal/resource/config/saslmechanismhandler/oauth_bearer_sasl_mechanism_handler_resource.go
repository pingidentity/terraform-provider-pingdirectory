package saslmechanismhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	_ resource.Resource                = &oauthBearerSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &oauthBearerSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &oauthBearerSaslMechanismHandlerResource{}
	_ resource.Resource                = &defaultOauthBearerSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultOauthBearerSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &defaultOauthBearerSaslMechanismHandlerResource{}
)

// Create a Oauth Bearer Sasl Mechanism Handler resource
func NewOauthBearerSaslMechanismHandlerResource() resource.Resource {
	return &oauthBearerSaslMechanismHandlerResource{}
}

func NewDefaultOauthBearerSaslMechanismHandlerResource() resource.Resource {
	return &defaultOauthBearerSaslMechanismHandlerResource{}
}

// oauthBearerSaslMechanismHandlerResource is the resource implementation.
type oauthBearerSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultOauthBearerSaslMechanismHandlerResource is the resource implementation.
type defaultOauthBearerSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *oauthBearerSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_bearer_sasl_mechanism_handler"
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_oauth_bearer_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *oauthBearerSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type oauthBearerSaslMechanismHandlerResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	AccessTokenValidator                         types.Set    `tfsdk:"access_token_validator"`
	IdTokenValidator                             types.Set    `tfsdk:"id_token_validator"`
	RequireBothAccessTokenAndIDToken             types.Bool   `tfsdk:"require_both_access_token_and_id_token"`
	ValidateAccessTokenWhenIDTokenIsAlsoProvided types.String `tfsdk:"validate_access_token_when_id_token_is_also_provided"`
	AlternateAuthorizationIdentityMapper         types.String `tfsdk:"alternate_authorization_identity_mapper"`
	AllRequiredScope                             types.Set    `tfsdk:"all_required_scope"`
	AnyRequiredScope                             types.Set    `tfsdk:"any_required_scope"`
	ServerFqdn                                   types.String `tfsdk:"server_fqdn"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *oauthBearerSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthBearerSaslMechanismHandlerSchema(ctx, req, resp, false)
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthBearerSaslMechanismHandlerSchema(ctx, req, resp, true)
}

func oauthBearerSaslMechanismHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Oauth Bearer Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"access_token_validator": schema.SetAttribute{
				Description: "An access token validator that will ensure that each presented OAuth access token is authentic and trustworthy. It must be configured with an identity mapper that will be used to map the access token to a local entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "An ID token validator that will ensure that each presented OpenID Connect ID token is authentic and trustworthy, and that will map the token to a local entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"require_both_access_token_and_id_token": schema.BoolAttribute{
				Description: "Indicates whether bind requests will be required to have both an OAuth access token (in the \"auth\" element of the bind request) and an OpenID Connect ID token (in the \"pingidentityidtoken\" element of the bind request).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"validate_access_token_when_id_token_is_also_provided": schema.StringAttribute{
				Description: "Indicates whether to validate the OAuth access token in addition to the OpenID Connect ID token in OAUTHBEARER bind requests that contain both types of tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alternate_authorization_identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"all_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that will all be required for any access tokens that will be allowed for authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that a token may have to be allowed for authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"server_fqdn": schema.StringAttribute{
				Description: "The fully-qualified name that clients are expected to use when communicating with the server.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
				Required:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add config validators
func (r oauthBearerSaslMechanismHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("id_token_validator"),
			path.MatchRoot("access_token_validator"),
		),
	}
}

// Add optional fields to create request
func addOptionalOauthBearerSaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddOauthBearerSaslMechanismHandlerRequest, plan oauthBearerSaslMechanismHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.AccessTokenValidator) {
		var slice []string
		plan.AccessTokenValidator.ElementsAs(ctx, &slice, false)
		addRequest.AccessTokenValidator = slice
	}
	if internaltypes.IsDefined(plan.IdTokenValidator) {
		var slice []string
		plan.IdTokenValidator.ElementsAs(ctx, &slice, false)
		addRequest.IdTokenValidator = slice
	}
	if internaltypes.IsDefined(plan.RequireBothAccessTokenAndIDToken) {
		addRequest.RequireBothAccessTokenAndIDToken = plan.RequireBothAccessTokenAndIDToken.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidateAccessTokenWhenIDTokenIsAlsoProvided) {
		validateAccessTokenWhenIDTokenIsAlsoProvided, err := client.NewEnumsaslMechanismHandlerValidateAccessTokenWhenIDTokenIsAlsoProvidedPropFromValue(plan.ValidateAccessTokenWhenIDTokenIsAlsoProvided.ValueString())
		if err != nil {
			return err
		}
		addRequest.ValidateAccessTokenWhenIDTokenIsAlsoProvided = validateAccessTokenWhenIDTokenIsAlsoProvided
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AlternateAuthorizationIdentityMapper) {
		addRequest.AlternateAuthorizationIdentityMapper = plan.AlternateAuthorizationIdentityMapper.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllRequiredScope) {
		var slice []string
		plan.AllRequiredScope.ElementsAs(ctx, &slice, false)
		addRequest.AllRequiredScope = slice
	}
	if internaltypes.IsDefined(plan.AnyRequiredScope) {
		var slice []string
		plan.AnyRequiredScope.ElementsAs(ctx, &slice, false)
		addRequest.AnyRequiredScope = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerFqdn) {
		addRequest.ServerFqdn = plan.ServerFqdn.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Read a OauthBearerSaslMechanismHandlerResponse object into the model struct
func readOauthBearerSaslMechanismHandlerResponse(ctx context.Context, r *client.OauthBearerSaslMechanismHandlerResponse, state *oauthBearerSaslMechanismHandlerResourceModel, expectedValues *oauthBearerSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireBothAccessTokenAndIDToken = internaltypes.BoolTypeOrNil(r.RequireBothAccessTokenAndIDToken)
	state.ValidateAccessTokenWhenIDTokenIsAlsoProvided = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerValidateAccessTokenWhenIDTokenIsAlsoProvidedProp(r.ValidateAccessTokenWhenIDTokenIsAlsoProvided), internaltypes.IsEmptyString(expectedValues.ValidateAccessTokenWhenIDTokenIsAlsoProvided))
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, internaltypes.IsEmptyString(expectedValues.AlternateAuthorizationIdentityMapper))
	state.AllRequiredScope = internaltypes.GetStringSet(r.AllRequiredScope)
	state.AnyRequiredScope = internaltypes.GetStringSet(r.AnyRequiredScope)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, internaltypes.IsEmptyString(expectedValues.ServerFqdn))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createOauthBearerSaslMechanismHandlerOperations(plan oauthBearerSaslMechanismHandlerResourceModel, state oauthBearerSaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBothAccessTokenAndIDToken, state.RequireBothAccessTokenAndIDToken, "require-both-access-token-and-id-token")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidateAccessTokenWhenIDTokenIsAlsoProvided, state.ValidateAccessTokenWhenIDTokenIsAlsoProvided, "validate-access-token-when-id-token-is-also-provided")
	operations.AddStringOperationIfNecessary(&ops, plan.AlternateAuthorizationIdentityMapper, state.AlternateAuthorizationIdentityMapper, "alternate-authorization-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllRequiredScope, state.AllRequiredScope, "all-required-scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredScope, state.AnyRequiredScope, "any-required-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerFqdn, state.ServerFqdn, "server-fqdn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *oauthBearerSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan oauthBearerSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddOauthBearerSaslMechanismHandlerRequest(plan.Id.ValueString(),
		[]client.EnumoauthBearerSaslMechanismHandlerSchemaUrn{client.ENUMOAUTHBEARERSASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLEROAUTH_BEARER},
		plan.Enabled.ValueBool())
	err := addOptionalOauthBearerSaslMechanismHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Oauth Bearer Sasl Mechanism Handler", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddOauthBearerSaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthBearerSaslMechanismHandlerResourceModel
	readOauthBearerSaslMechanismHandlerResponse(ctx, addResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultOauthBearerSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan oauthBearerSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state oauthBearerSaslMechanismHandlerResourceModel
	readOauthBearerSaslMechanismHandlerResponse(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createOauthBearerSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readOauthBearerSaslMechanismHandlerResponse(ctx, updateResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *oauthBearerSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthBearerSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthBearerSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOauthBearerSaslMechanismHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state oauthBearerSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthBearerSaslMechanismHandlerResponse(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *oauthBearerSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthBearerSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthBearerSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOauthBearerSaslMechanismHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan oauthBearerSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthBearerSaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createOauthBearerSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readOauthBearerSaslMechanismHandlerResponse(ctx, updateResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultOauthBearerSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *oauthBearerSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthBearerSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandlerExecute(r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Oauth Bearer Sasl Mechanism Handler", err, httpResp)
		return
	}
}

func (r *oauthBearerSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthBearerSaslMechanismHandler(ctx, req, resp)
}

func (r *defaultOauthBearerSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthBearerSaslMechanismHandler(ctx, req, resp)
}

func importOauthBearerSaslMechanismHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
