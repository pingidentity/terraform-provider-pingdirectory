package idtokenvalidator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &openidConnectIdTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &openidConnectIdTokenValidatorResource{}
	_ resource.ResourceWithImportState = &openidConnectIdTokenValidatorResource{}
	_ resource.Resource                = &defaultOpenidConnectIdTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultOpenidConnectIdTokenValidatorResource{}
	_ resource.ResourceWithImportState = &defaultOpenidConnectIdTokenValidatorResource{}
)

// Create a Openid Connect Id Token Validator resource
func NewOpenidConnectIdTokenValidatorResource() resource.Resource {
	return &openidConnectIdTokenValidatorResource{}
}

func NewDefaultOpenidConnectIdTokenValidatorResource() resource.Resource {
	return &defaultOpenidConnectIdTokenValidatorResource{}
}

// openidConnectIdTokenValidatorResource is the resource implementation.
type openidConnectIdTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultOpenidConnectIdTokenValidatorResource is the resource implementation.
type defaultOpenidConnectIdTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *openidConnectIdTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openid_connect_id_token_validator"
}

func (r *defaultOpenidConnectIdTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_openid_connect_id_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *openidConnectIdTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultOpenidConnectIdTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type openidConnectIdTokenValidatorResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	AllowedSigningAlgorithm types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate      types.Set    `tfsdk:"signing_certificate"`
	OpenIDConnectProvider   types.String `tfsdk:"openid_connect_provider"`
	JwksEndpointPath        types.String `tfsdk:"jwks_endpoint_path"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	IdentityMapper          types.String `tfsdk:"identity_mapper"`
	SubjectClaimName        types.String `tfsdk:"subject_claim_name"`
	IssuerURL               types.String `tfsdk:"issuer_url"`
	ClockSkewGracePeriod    types.String `tfsdk:"clock_skew_grace_period"`
	JwksCacheDuration       types.String `tfsdk:"jwks_cache_duration"`
	EvaluationOrderIndex    types.Int64  `tfsdk:"evaluation_order_index"`
}

// GetSchema defines the schema for the resource.
func (r *openidConnectIdTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	openidConnectIdTokenValidatorSchema(ctx, req, resp, false)
}

func (r *defaultOpenidConnectIdTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	openidConnectIdTokenValidatorSchema(ctx, req, resp, true)
}

func openidConnectIdTokenValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Openid Connect Id Token Validator.",
		Attributes: map[string]schema.Attribute{
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the OpenID Connect ID Token Validator.",
				Required:    true,
				ElementType: types.StringType,
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming ID token. This property may be specified if a JWKS endpoint should not be used to retrieve public signing keys.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"openid_connect_provider": schema.StringAttribute{
				Description: "Specifies the OpenID Connect provider that issues ID tokens handled by this OpenID Connect ID Token Validator. This property is used in conjunction with the jwks-endpoint-path property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to the JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming ID token. This path is relative to the base_url property defined for the validator's OpenID Connect provider. If jwks-endpoint-path is specified, the OpenID Connect ID Token Validator will not consult locally stored certificates for validating token signatures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this ID Token Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this ID Token Validator is enabled for use in the Directory Server.",
				Required:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used to correlate an ID token subject value to a user entry. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Required:    true,
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject; i.e., the authenticated user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"issuer_url": schema.StringAttribute{
				Description: "Specifies the OpenID Connect provider's issuer URL.",
				Required:    true,
			},
			"clock_skew_grace_period": schema.StringAttribute{
				Description: "Specifies the amount of clock skew that is tolerated by the ID Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_cache_duration": schema.StringAttribute{
				Description: "How often the ID Token Validator should refresh its cache of JWKS token signing keys.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple ID Token Validators are defined for a single Directory Server, this property determines the order in which the ID Token Validators are consulted. Values of this property must be unique among all ID Token Validators defined within Directory Server but not necessarily contiguous. ID Token Validators with lower values will be evaluated first to determine if they are able to validate the ID token.",
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

// Add optional fields to create request
func addOptionalOpenidConnectIdTokenValidatorFields(ctx context.Context, addRequest *client.AddOpenidConnectIdTokenValidatorRequest, plan openidConnectIdTokenValidatorResourceModel) {
	if internaltypes.IsDefined(plan.SigningCertificate) {
		var slice []string
		plan.SigningCertificate.ElementsAs(ctx, &slice, false)
		addRequest.SigningCertificate = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OpenIDConnectProvider) {
		addRequest.OpenIDConnectProvider = plan.OpenIDConnectProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JwksEndpointPath) {
		addRequest.JwksEndpointPath = plan.JwksEndpointPath.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		addRequest.SubjectClaimName = plan.SubjectClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClockSkewGracePeriod) {
		addRequest.ClockSkewGracePeriod = plan.ClockSkewGracePeriod.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JwksCacheDuration) {
		addRequest.JwksCacheDuration = plan.JwksCacheDuration.ValueStringPointer()
	}
}

// Read a OpenidConnectIdTokenValidatorResponse object into the model struct
func readOpenidConnectIdTokenValidatorResponse(ctx context.Context, r *client.OpenidConnectIdTokenValidatorResponse, state *openidConnectIdTokenValidatorResourceModel, expectedValues *openidConnectIdTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllowedSigningAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumidTokenValidatorAllowedSigningAlgorithmProp(r.AllowedSigningAlgorithm))
	state.SigningCertificate = internaltypes.GetStringSet(r.SigningCertificate)
	state.OpenIDConnectProvider = internaltypes.StringTypeOrNil(r.OpenIDConnectProvider, internaltypes.IsEmptyString(expectedValues.OpenIDConnectProvider))
	state.JwksEndpointPath = internaltypes.StringTypeOrNil(r.JwksEndpointPath, internaltypes.IsEmptyString(expectedValues.JwksEndpointPath))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.IssuerURL = types.StringValue(r.IssuerURL)
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, internaltypes.IsEmptyString(expectedValues.ClockSkewGracePeriod))
	config.CheckMismatchedPDFormattedAttributes("clock_skew_grace_period",
		expectedValues.ClockSkewGracePeriod, state.ClockSkewGracePeriod, diagnostics)
	state.JwksCacheDuration = internaltypes.StringTypeOrNil(r.JwksCacheDuration, internaltypes.IsEmptyString(expectedValues.JwksCacheDuration))
	config.CheckMismatchedPDFormattedAttributes("jwks_cache_duration",
		expectedValues.JwksCacheDuration, state.JwksCacheDuration, diagnostics)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createOpenidConnectIdTokenValidatorOperations(plan openidConnectIdTokenValidatorResourceModel, state openidConnectIdTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedSigningAlgorithm, state.AllowedSigningAlgorithm, "allowed-signing-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SigningCertificate, state.SigningCertificate, "signing-certificate")
	operations.AddStringOperationIfNecessary(&ops, plan.OpenIDConnectProvider, state.OpenIDConnectProvider, "openid-connect-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksEndpointPath, state.JwksEndpointPath, "jwks-endpoint-path")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.IssuerURL, state.IssuerURL, "issuer-url")
	operations.AddStringOperationIfNecessary(&ops, plan.ClockSkewGracePeriod, state.ClockSkewGracePeriod, "clock-skew-grace-period")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksCacheDuration, state.JwksCacheDuration, "jwks-cache-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	return ops
}

// Create a new resource
func (r *openidConnectIdTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan openidConnectIdTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var AllowedSigningAlgorithmSlice []client.EnumidTokenValidatorAllowedSigningAlgorithmProp
	plan.AllowedSigningAlgorithm.ElementsAs(ctx, &AllowedSigningAlgorithmSlice, false)
	addRequest := client.NewAddOpenidConnectIdTokenValidatorRequest(plan.Id.ValueString(),
		[]client.EnumopenidConnectIdTokenValidatorSchemaUrn{client.ENUMOPENIDCONNECTIDTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ID_TOKEN_VALIDATOROPENID_CONNECT},
		AllowedSigningAlgorithmSlice,
		plan.Enabled.ValueBool(),
		plan.IdentityMapper.ValueString(),
		plan.IssuerURL.ValueString(),
		plan.EvaluationOrderIndex.ValueInt64())
	addOptionalOpenidConnectIdTokenValidatorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdTokenValidatorApi.AddIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdTokenValidatorRequest(
		client.AddOpenidConnectIdTokenValidatorRequestAsAddIdTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.AddIdTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Openid Connect Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state openidConnectIdTokenValidatorResourceModel
	readOpenidConnectIdTokenValidatorResponse(ctx, addResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultOpenidConnectIdTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan openidConnectIdTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.GetIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Openid Connect Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state openidConnectIdTokenValidatorResourceModel
	readOpenidConnectIdTokenValidatorResponse(ctx, readResponse.OpenidConnectIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.IdTokenValidatorApi.UpdateIdTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createOpenidConnectIdTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.UpdateIdTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Openid Connect Id Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readOpenidConnectIdTokenValidatorResponse(ctx, updateResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *openidConnectIdTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOpenidConnectIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOpenidConnectIdTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOpenidConnectIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOpenidConnectIdTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state openidConnectIdTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.IdTokenValidatorApi.GetIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Openid Connect Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOpenidConnectIdTokenValidatorResponse(ctx, readResponse.OpenidConnectIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *openidConnectIdTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOpenidConnectIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOpenidConnectIdTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOpenidConnectIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOpenidConnectIdTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan openidConnectIdTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state openidConnectIdTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.IdTokenValidatorApi.UpdateIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createOpenidConnectIdTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.IdTokenValidatorApi.UpdateIdTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Openid Connect Id Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readOpenidConnectIdTokenValidatorResponse(ctx, updateResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultOpenidConnectIdTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *openidConnectIdTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state openidConnectIdTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdTokenValidatorApi.DeleteIdTokenValidatorExecute(r.apiClient.IdTokenValidatorApi.DeleteIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Openid Connect Id Token Validator", err, httpResp)
		return
	}
}

func (r *openidConnectIdTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOpenidConnectIdTokenValidator(ctx, req, resp)
}

func (r *defaultOpenidConnectIdTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOpenidConnectIdTokenValidator(ctx, req, resp)
}

func importOpenidConnectIdTokenValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
