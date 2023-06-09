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
	_ resource.Resource                = &idTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &idTokenValidatorResource{}
	_ resource.ResourceWithImportState = &idTokenValidatorResource{}
	_ resource.Resource                = &defaultIdTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultIdTokenValidatorResource{}
	_ resource.ResourceWithImportState = &defaultIdTokenValidatorResource{}
)

// Create a Id Token Validator resource
func NewIdTokenValidatorResource() resource.Resource {
	return &idTokenValidatorResource{}
}

func NewDefaultIdTokenValidatorResource() resource.Resource {
	return &defaultIdTokenValidatorResource{}
}

// idTokenValidatorResource is the resource implementation.
type idTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultIdTokenValidatorResource is the resource implementation.
type defaultIdTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *idTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_id_token_validator"
}

func (r *defaultIdTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_id_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *idTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultIdTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type idTokenValidatorResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	Type                               types.String `tfsdk:"type"`
	AllowedSigningAlgorithm            types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate                 types.Set    `tfsdk:"signing_certificate"`
	IssuerURL                          types.String `tfsdk:"issuer_url"`
	JwksEndpointPath                   types.String `tfsdk:"jwks_endpoint_path"`
	OpenIDConnectProvider              types.String `tfsdk:"openid_connect_provider"`
	OpenIDConnectMetadataCacheDuration types.String `tfsdk:"openid_connect_metadata_cache_duration"`
	Description                        types.String `tfsdk:"description"`
	Enabled                            types.Bool   `tfsdk:"enabled"`
	IdentityMapper                     types.String `tfsdk:"identity_mapper"`
	SubjectClaimName                   types.String `tfsdk:"subject_claim_name"`
	ClockSkewGracePeriod               types.String `tfsdk:"clock_skew_grace_period"`
	JwksCacheDuration                  types.String `tfsdk:"jwks_cache_duration"`
	EvaluationOrderIndex               types.Int64  `tfsdk:"evaluation_order_index"`
}

// GetSchema defines the schema for the resource.
func (r *idTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idTokenValidatorSchema(ctx, req, resp, false)
}

func (r *defaultIdTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idTokenValidatorSchema(ctx, req, resp, true)
}

func idTokenValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Id Token Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of ID Token Validator resource. Options are ['ping-one', 'openid-connect']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the OpenID Connect ID Token Validator.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming ID token. This property may be specified if a JWKS endpoint should not be used to retrieve public signing keys.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"issuer_url": schema.StringAttribute{
				Description: "Specifies a PingOne base issuer URL.",
				Required:    true,
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to the JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming ID token. This path is relative to the base_url property defined for the validator's OpenID Connect provider. If jwks-endpoint-path is specified, the OpenID Connect ID Token Validator will not consult locally stored certificates for validating token signatures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"openid_connect_provider": schema.StringAttribute{
				Description: "Specifies HTTPS connection settings for the PingOne OpenID Connect provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"openid_connect_metadata_cache_duration": schema.StringAttribute{
				Description: "How often the PingOne ID Token Validator should refresh its stored cache of OpenID Connect-related metadata.",
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
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *idTokenValidatorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdTokenValidatorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanIdTokenValidator(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model idTokenValidatorResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.AllowedSigningAlgorithm) && model.Type.ValueString() != "openid-connect" {
		resp.Diagnostics.AddError("Attribute 'allowed_signing_algorithm' not supported by pingdirectory_id_token_validator resources with type '"+model.Type.ValueString()+"'", "")
	}
	if internaltypes.IsDefined(model.SigningCertificate) && model.Type.ValueString() != "openid-connect" {
		resp.Diagnostics.AddError("Attribute 'signing_certificate' not supported by pingdirectory_id_token_validator resources with type '"+model.Type.ValueString()+"'", "")
	}
	if internaltypes.IsDefined(model.JwksEndpointPath) && model.Type.ValueString() != "openid-connect" {
		resp.Diagnostics.AddError("Attribute 'jwks_endpoint_path' not supported by pingdirectory_id_token_validator resources with type '"+model.Type.ValueString()+"'", "")
	}
	if internaltypes.IsDefined(model.OpenIDConnectMetadataCacheDuration) && model.Type.ValueString() != "ping-one" {
		resp.Diagnostics.AddError("Attribute 'openid_connect_metadata_cache_duration' not supported by pingdirectory_id_token_validator resources with type '"+model.Type.ValueString()+"'", "")
	}
}

// Add optional fields to create request for ping-one id-token-validator
func addOptionalPingOneIdTokenValidatorFields(ctx context.Context, addRequest *client.AddPingOneIdTokenValidatorRequest, plan idTokenValidatorResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OpenIDConnectProvider) {
		addRequest.OpenIDConnectProvider = plan.OpenIDConnectProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OpenIDConnectMetadataCacheDuration) {
		addRequest.OpenIDConnectMetadataCacheDuration = plan.OpenIDConnectMetadataCacheDuration.ValueStringPointer()
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

// Add optional fields to create request for openid-connect id-token-validator
func addOptionalOpenidConnectIdTokenValidatorFields(ctx context.Context, addRequest *client.AddOpenidConnectIdTokenValidatorRequest, plan idTokenValidatorResourceModel) {
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

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateIdTokenValidatorNilSets(ctx context.Context, model *idTokenValidatorResourceModel) {
	if model.SigningCertificate.ElementType(ctx) == nil {
		model.SigningCertificate = types.SetNull(types.StringType)
	}
	if model.AllowedSigningAlgorithm.ElementType(ctx) == nil {
		model.AllowedSigningAlgorithm = types.SetNull(types.StringType)
	}
}

// Read a PingOneIdTokenValidatorResponse object into the model struct
func readPingOneIdTokenValidatorResponse(ctx context.Context, r *client.PingOneIdTokenValidatorResponse, state *idTokenValidatorResourceModel, expectedValues *idTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one")
	state.Id = types.StringValue(r.Id)
	state.IssuerURL = types.StringValue(r.IssuerURL)
	state.OpenIDConnectProvider = types.StringValue(r.OpenIDConnectProvider)
	state.OpenIDConnectMetadataCacheDuration = internaltypes.StringTypeOrNil(r.OpenIDConnectMetadataCacheDuration, internaltypes.IsEmptyString(expectedValues.OpenIDConnectMetadataCacheDuration))
	config.CheckMismatchedPDFormattedAttributes("openid_connect_metadata_cache_duration",
		expectedValues.OpenIDConnectMetadataCacheDuration, state.OpenIDConnectMetadataCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, internaltypes.IsEmptyString(expectedValues.ClockSkewGracePeriod))
	config.CheckMismatchedPDFormattedAttributes("clock_skew_grace_period",
		expectedValues.ClockSkewGracePeriod, state.ClockSkewGracePeriod, diagnostics)
	state.JwksCacheDuration = internaltypes.StringTypeOrNil(r.JwksCacheDuration, internaltypes.IsEmptyString(expectedValues.JwksCacheDuration))
	config.CheckMismatchedPDFormattedAttributes("jwks_cache_duration",
		expectedValues.JwksCacheDuration, state.JwksCacheDuration, diagnostics)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdTokenValidatorNilSets(ctx, state)
}

// Read a OpenidConnectIdTokenValidatorResponse object into the model struct
func readOpenidConnectIdTokenValidatorResponse(ctx context.Context, r *client.OpenidConnectIdTokenValidatorResponse, state *idTokenValidatorResourceModel, expectedValues *idTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("openid-connect")
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
	populateIdTokenValidatorNilSets(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createIdTokenValidatorOperations(plan idTokenValidatorResourceModel, state idTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedSigningAlgorithm, state.AllowedSigningAlgorithm, "allowed-signing-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SigningCertificate, state.SigningCertificate, "signing-certificate")
	operations.AddStringOperationIfNecessary(&ops, plan.IssuerURL, state.IssuerURL, "issuer-url")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksEndpointPath, state.JwksEndpointPath, "jwks-endpoint-path")
	operations.AddStringOperationIfNecessary(&ops, plan.OpenIDConnectProvider, state.OpenIDConnectProvider, "openid-connect-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.OpenIDConnectMetadataCacheDuration, state.OpenIDConnectMetadataCacheDuration, "openid-connect-metadata-cache-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ClockSkewGracePeriod, state.ClockSkewGracePeriod, "clock-skew-grace-period")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksCacheDuration, state.JwksCacheDuration, "jwks-cache-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	return ops
}

// Create a ping-one id-token-validator
func (r *idTokenValidatorResource) CreatePingOneIdTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan idTokenValidatorResourceModel) (*idTokenValidatorResourceModel, error) {
	addRequest := client.NewAddPingOneIdTokenValidatorRequest(plan.Id.ValueString(),
		[]client.EnumpingOneIdTokenValidatorSchemaUrn{client.ENUMPINGONEIDTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ID_TOKEN_VALIDATORPING_ONE},
		plan.IssuerURL.ValueString(),
		plan.Enabled.ValueBool(),
		plan.IdentityMapper.ValueString(),
		plan.EvaluationOrderIndex.ValueInt64())
	addOptionalPingOneIdTokenValidatorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdTokenValidatorApi.AddIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdTokenValidatorRequest(
		client.AddPingOneIdTokenValidatorRequestAsAddIdTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.AddIdTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Id Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idTokenValidatorResourceModel
	readPingOneIdTokenValidatorResponse(ctx, addResponse.PingOneIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a openid-connect id-token-validator
func (r *idTokenValidatorResource) CreateOpenidConnectIdTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan idTokenValidatorResourceModel) (*idTokenValidatorResourceModel, error) {
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Id Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idTokenValidatorResourceModel
	readOpenidConnectIdTokenValidatorResponse(ctx, addResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *idTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan idTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *idTokenValidatorResourceModel
	var err error
	if plan.Type.ValueString() == "ping-one" {
		state, err = r.CreatePingOneIdTokenValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "openid-connect" {
		state, err = r.CreateOpenidConnectIdTokenValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultIdTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan idTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.GetIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state idTokenValidatorResourceModel
	if plan.Type.ValueString() == "ping-one" {
		readPingOneIdTokenValidatorResponse(ctx, readResponse.PingOneIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "openid-connect" {
		readOpenidConnectIdTokenValidatorResponse(ctx, readResponse.OpenidConnectIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.IdTokenValidatorApi.UpdateIdTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createIdTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.UpdateIdTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Id Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ping-one" {
			readPingOneIdTokenValidatorResponse(ctx, updateResponse.PingOneIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "openid-connect" {
			readOpenidConnectIdTokenValidatorResponse(ctx, updateResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *idTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state idTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.IdTokenValidatorApi.GetIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PingOneIdTokenValidatorResponse != nil {
		readPingOneIdTokenValidatorResponse(ctx, readResponse.PingOneIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OpenidConnectIdTokenValidatorResponse != nil {
		readOpenidConnectIdTokenValidatorResponse(ctx, readResponse.OpenidConnectIdTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *idTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan idTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state idTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.IdTokenValidatorApi.UpdateIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createIdTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.IdTokenValidatorApi.UpdateIdTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Id Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ping-one" {
			readPingOneIdTokenValidatorResponse(ctx, updateResponse.PingOneIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "openid-connect" {
			readOpenidConnectIdTokenValidatorResponse(ctx, updateResponse.OpenidConnectIdTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultIdTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *idTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdTokenValidatorApi.DeleteIdTokenValidatorExecute(r.apiClient.IdTokenValidatorApi.DeleteIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Id Token Validator", err, httpResp)
		return
	}
}

func (r *idTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdTokenValidator(ctx, req, resp)
}

func (r *defaultIdTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdTokenValidator(ctx, req, resp)
}

func importIdTokenValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
