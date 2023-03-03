package accesstokenvalidator

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jwtAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &jwtAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &jwtAccessTokenValidatorResource{}
	_ resource.Resource                = &defaultJwtAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultJwtAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &defaultJwtAccessTokenValidatorResource{}
)

// Create a Jwt Access Token Validator resource
func NewJwtAccessTokenValidatorResource() resource.Resource {
	return &jwtAccessTokenValidatorResource{}
}

func NewDefaultJwtAccessTokenValidatorResource() resource.Resource {
	return &defaultJwtAccessTokenValidatorResource{}
}

// jwtAccessTokenValidatorResource is the resource implementation.
type jwtAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJwtAccessTokenValidatorResource is the resource implementation.
type defaultJwtAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jwtAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwt_access_token_validator"
}

func (r *defaultJwtAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_jwt_access_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *jwtAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultJwtAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type jwtAccessTokenValidatorResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	AllowedSigningAlgorithm           types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate                types.Set    `tfsdk:"signing_certificate"`
	JwksEndpointPath                  types.String `tfsdk:"jwks_endpoint_path"`
	EncryptionKeyPair                 types.String `tfsdk:"encryption_key_pair"`
	AllowedKeyEncryptionAlgorithm     types.Set    `tfsdk:"allowed_key_encryption_algorithm"`
	AllowedContentEncryptionAlgorithm types.Set    `tfsdk:"allowed_content_encryption_algorithm"`
	ClockSkewGracePeriod              types.String `tfsdk:"clock_skew_grace_period"`
	ClientIDClaimName                 types.String `tfsdk:"client_id_claim_name"`
	ScopeClaimName                    types.String `tfsdk:"scope_claim_name"`
	EvaluationOrderIndex              types.Int64  `tfsdk:"evaluation_order_index"`
	AuthorizationServer               types.String `tfsdk:"authorization_server"`
	IdentityMapper                    types.String `tfsdk:"identity_mapper"`
	SubjectClaimName                  types.String `tfsdk:"subject_claim_name"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *jwtAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jwtAccessTokenValidatorSchema(ctx, req, resp, false)
}

func (r *defaultJwtAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jwtAccessTokenValidatorSchema(ctx, req, resp, true)
}

func jwtAccessTokenValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Jwt Access Token Validator.",
		Attributes: map[string]schema.Attribute{
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the JWT Access Token Validator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming JWT access token. If this property is specified, the JWT Access Token Validator will not use a JWKS endpoint to retrieve public keys.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming JWT access token. This path is relative to the base_url property defined for the validator's external authorization server. If jwks-endpoint-path is specified, the JWT Access Token Validator will not consult locally stored certificates for validating token signatures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_key_pair": schema.StringAttribute{
				Description: "The public-private key pair that is used to encrypt the JWT payload. If specified, the JWT Access Token Validator will use the private key to decrypt the JWT payload, and the public key must be exported to the Authorization Server that is issuing access tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_key_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT key encryption algorithms that will be accepted by the JWT Access Token Validator. This setting is only used if encryption-key-pair is set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"allowed_content_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT content encryption algorithms that will be accepted by the JWT Access Token Validator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"clock_skew_grace_period": schema.StringAttribute{
				Description: "Specifies the amount of clock skew that is tolerated by the JWT Access Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the OAuth2 client Id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the scopes granted by the token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple JWT Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all JWT Access Token Validators defined within Directory Server but not necessarily contiguous. JWT Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authorization_server": schema.StringAttribute{
				Description: "Specifies the external server that will be used to aid in validating access tokens. In most cases this will be the Authorization Server that minted the token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Access Token Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Token Validator is enabled for use in Directory Server.",
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
func addOptionalJwtAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddJwtAccessTokenValidatorRequest, plan jwtAccessTokenValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedSigningAlgorithm) {
		var slice []string
		plan.AllowedSigningAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedSigningAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedSigningAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedSigningAlgorithm = enumSlice
	}
	if internaltypes.IsDefined(plan.SigningCertificate) {
		var slice []string
		plan.SigningCertificate.ElementsAs(ctx, &slice, false)
		addRequest.SigningCertificate = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JwksEndpointPath) {
		stringVal := plan.JwksEndpointPath.ValueString()
		addRequest.JwksEndpointPath = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionKeyPair) {
		stringVal := plan.EncryptionKeyPair.ValueString()
		addRequest.EncryptionKeyPair = &stringVal
	}
	if internaltypes.IsDefined(plan.AllowedKeyEncryptionAlgorithm) {
		var slice []string
		plan.AllowedKeyEncryptionAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedKeyEncryptionAlgorithm = enumSlice
	}
	if internaltypes.IsDefined(plan.AllowedContentEncryptionAlgorithm) {
		var slice []string
		plan.AllowedContentEncryptionAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedContentEncryptionAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedContentEncryptionAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedContentEncryptionAlgorithm = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClockSkewGracePeriod) {
		stringVal := plan.ClockSkewGracePeriod.ValueString()
		addRequest.ClockSkewGracePeriod = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientIDClaimName) {
		stringVal := plan.ClientIDClaimName.ValueString()
		addRequest.ClientIDClaimName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ScopeClaimName) {
		stringVal := plan.ScopeClaimName.ValueString()
		addRequest.ScopeClaimName = &stringVal
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		intVal := int32(plan.EvaluationOrderIndex.ValueInt64())
		addRequest.EvaluationOrderIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthorizationServer) {
		stringVal := plan.AuthorizationServer.ValueString()
		addRequest.AuthorizationServer = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		stringVal := plan.IdentityMapper.ValueString()
		addRequest.IdentityMapper = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		stringVal := plan.SubjectClaimName.ValueString()
		addRequest.SubjectClaimName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a JwtAccessTokenValidatorResponse object into the model struct
func readJwtAccessTokenValidatorResponse(ctx context.Context, r *client.JwtAccessTokenValidatorResponse, state *jwtAccessTokenValidatorResourceModel, expectedValues *jwtAccessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllowedSigningAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedSigningAlgorithmProp(r.AllowedSigningAlgorithm))
	state.SigningCertificate = internaltypes.GetStringSet(r.SigningCertificate)
	state.JwksEndpointPath = internaltypes.StringTypeOrNil(r.JwksEndpointPath, internaltypes.IsEmptyString(expectedValues.JwksEndpointPath))
	state.EncryptionKeyPair = internaltypes.StringTypeOrNil(r.EncryptionKeyPair, internaltypes.IsEmptyString(expectedValues.EncryptionKeyPair))
	state.AllowedKeyEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmProp(r.AllowedKeyEncryptionAlgorithm))
	state.AllowedContentEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedContentEncryptionAlgorithmProp(r.AllowedContentEncryptionAlgorithm))
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, internaltypes.IsEmptyString(expectedValues.ClockSkewGracePeriod))
	config.CheckMismatchedPDFormattedAttributes("clock_skew_grace_period",
		expectedValues.ClockSkewGracePeriod, state.ClockSkewGracePeriod, diagnostics)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, internaltypes.IsEmptyString(expectedValues.ClientIDClaimName))
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, internaltypes.IsEmptyString(expectedValues.ScopeClaimName))
	state.EvaluationOrderIndex = types.Int64Value(int64(r.EvaluationOrderIndex))
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, internaltypes.IsEmptyString(expectedValues.AuthorizationServer))
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createJwtAccessTokenValidatorOperations(plan jwtAccessTokenValidatorResourceModel, state jwtAccessTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedSigningAlgorithm, state.AllowedSigningAlgorithm, "allowed-signing-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SigningCertificate, state.SigningCertificate, "signing-certificate")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksEndpointPath, state.JwksEndpointPath, "jwks-endpoint-path")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionKeyPair, state.EncryptionKeyPair, "encryption-key-pair")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedKeyEncryptionAlgorithm, state.AllowedKeyEncryptionAlgorithm, "allowed-key-encryption-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedContentEncryptionAlgorithm, state.AllowedContentEncryptionAlgorithm, "allowed-content-encryption-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.ClockSkewGracePeriod, state.ClockSkewGracePeriod, "clock-skew-grace-period")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientIDClaimName, state.ClientIDClaimName, "client-id-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ScopeClaimName, state.ScopeClaimName, "scope-claim-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthorizationServer, state.AuthorizationServer, "authorization-server")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *jwtAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jwtAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddJwtAccessTokenValidatorRequest(plan.Id.ValueString(),
		[]client.EnumjwtAccessTokenValidatorSchemaUrn{client.ENUMJWTACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORJWT},
		plan.Enabled.ValueBool())
	err := addOptionalJwtAccessTokenValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Jwt Access Token Validator", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddJwtAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Jwt Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jwtAccessTokenValidatorResourceModel
	readJwtAccessTokenValidatorResponse(ctx, addResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultJwtAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jwtAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jwt Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jwtAccessTokenValidatorResourceModel
	readJwtAccessTokenValidatorResponse(ctx, readResponse.JwtAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createJwtAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jwt Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJwtAccessTokenValidatorResponse(ctx, updateResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jwtAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJwtAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJwtAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJwtAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readJwtAccessTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state jwtAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jwt Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJwtAccessTokenValidatorResponse(ctx, readResponse.JwtAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *jwtAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJwtAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJwtAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJwtAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJwtAccessTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jwtAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jwtAccessTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createJwtAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jwt Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJwtAccessTokenValidatorResponse(ctx, updateResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJwtAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jwtAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jwtAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidatorExecute(r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Jwt Access Token Validator", err, httpResp)
		return
	}
}

func (r *jwtAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJwtAccessTokenValidator(ctx, req, resp)
}

func (r *defaultJwtAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJwtAccessTokenValidator(ctx, req, resp)
}

func importJwtAccessTokenValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
