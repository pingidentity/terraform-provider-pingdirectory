// Copyright © 2025 Ping Identity Corporation

package accesstokenvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &accessTokenValidatorDataSource{}
	_ datasource.DataSourceWithConfigure = &accessTokenValidatorDataSource{}
)

// Create a Access Token Validator data source
func NewAccessTokenValidatorDataSource() datasource.DataSource {
	return &accessTokenValidatorDataSource{}
}

// accessTokenValidatorDataSource is the datasource implementation.
type accessTokenValidatorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *accessTokenValidatorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_token_validator"
}

// Configure adds the provider configured client to the data source.
func (r *accessTokenValidatorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type accessTokenValidatorDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	AllowedSigningAlgorithm           types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate                types.Set    `tfsdk:"signing_certificate"`
	JwksEndpointPath                  types.String `tfsdk:"jwks_endpoint_path"`
	EncryptionKeyPair                 types.String `tfsdk:"encryption_key_pair"`
	AllowedKeyEncryptionAlgorithm     types.Set    `tfsdk:"allowed_key_encryption_algorithm"`
	AllowedContentEncryptionAlgorithm types.Set    `tfsdk:"allowed_content_encryption_algorithm"`
	ClockSkewGracePeriod              types.String `tfsdk:"clock_skew_grace_period"`
	ClientIDClaimName                 types.String `tfsdk:"client_id_claim_name"`
	ScopeClaimName                    types.String `tfsdk:"scope_claim_name"`
	ClientID                          types.String `tfsdk:"client_id"`
	ClientSecret                      types.String `tfsdk:"client_secret"`
	ClientSecretPassphraseProvider    types.String `tfsdk:"client_secret_passphrase_provider"`
	IncludeAudParameter               types.Bool   `tfsdk:"include_aud_parameter"`
	AccessTokenManagerID              types.String `tfsdk:"access_token_manager_id"`
	EndpointCacheRefresh              types.String `tfsdk:"endpoint_cache_refresh"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	AuthorizationServer               types.String `tfsdk:"authorization_server"`
	PersistAccessTokens               types.Bool   `tfsdk:"persist_access_tokens"`
	MaximumTokenLifetime              types.String `tfsdk:"maximum_token_lifetime"`
	AllowedAuthenticationType         types.Set    `tfsdk:"allowed_authentication_type"`
	AllowedSASLMechanism              types.Set    `tfsdk:"allowed_sasl_mechanism"`
	GenerateTokenResultCriteria       types.String `tfsdk:"generate_token_result_criteria"`
	IncludedScope                     types.Set    `tfsdk:"included_scope"`
	IdentityMapper                    types.String `tfsdk:"identity_mapper"`
	SubjectClaimName                  types.String `tfsdk:"subject_claim_name"`
	Description                       types.String `tfsdk:"description"`
	EvaluationOrderIndex              types.Int64  `tfsdk:"evaluation_order_index"`
}

// GetSchema defines the schema for the datasource.
func (r *accessTokenValidatorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Access Token Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Access Token Validator resource. Options are ['bind', 'ping-federate', 'jwt', 'mock', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Access Token Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Access Token Validator. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the JWT Access Token Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming JWT access token. If this property is specified, the JWT Access Token Validator will not use a JWKS endpoint to retrieve public keys.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming JWT access token. This path is relative to the base_url property defined for the validator's external authorization server. If jwks-endpoint-path is specified, the JWT Access Token Validator will not consult locally stored certificates for validating token signatures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_key_pair": schema.StringAttribute{
				Description: "The public-private key pair that is used to encrypt the JWT payload. If specified, the JWT Access Token Validator will use the private key to decrypt the JWT payload, and the public key must be exported to the Authorization Server that is issuing access tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_key_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT key encryption algorithms that will be accepted by the JWT Access Token Validator. This setting is only used if encryption-key-pair is set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_content_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT content encryption algorithms that will be accepted by the JWT Access Token Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"clock_skew_grace_period": schema.StringAttribute{
				Description: "Specifies the amount of clock skew that is tolerated by the JWT Access Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"client_id_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the OAuth2 client Id.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scope_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the scopes granted by the token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client identifier to use when authenticating to the PingFederate authorization server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret to use when authenticating to the PingFederate authorization server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"client_secret_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider for obtaining the client secret to use when authenticating to the PingFederate authorization server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_aud_parameter": schema.BoolAttribute{
				Description: "Whether to include the incoming request URL as the \"aud\" parameter when calling the PingFederate introspection endpoint. This property is ignored if the access-token-manager-id property is set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"access_token_manager_id": schema.StringAttribute{
				Description: "The Access Token Manager instance ID to specify when calling the PingFederate introspection endpoint. If this property is set the include-aud-parameter property is ignored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"endpoint_cache_refresh": schema.StringAttribute{
				Description: "How often the Access Token Validator should refresh its stored value of the PingFederate server's token introspection endpoint.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`ping-federate`, `jwt`, `mock`, `third-party`]: Indicates whether this Access Token Validator is enabled for use in Directory Server. When the `type` attribute is set to `bind`: Indicates whether this Bind Access Token Validator is enabled for use in Directory Server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`ping-federate`, `jwt`, `mock`, `third-party`]: Indicates whether this Access Token Validator is enabled for use in Directory Server.\n  - `bind`: Indicates whether this Bind Access Token Validator is enabled for use in Directory Server.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"authorization_server": schema.StringAttribute{
				Description: "Specifies the external server that will be used to aid in validating access tokens. In most cases this will be the Authorization Server that minted the token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persist_access_tokens": schema.BoolAttribute{
				Description: "Indicates whether access tokens should be persisted in user entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_token_lifetime": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a generated token should be considered valid. If this is not specified, then generated access tokens will not expire.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "Specifies the authentication types for bind operations that may be used to generate access tokens, and for which generated access tokens will be accepted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_sasl_mechanism": schema.SetAttribute{
				Description: "Specifies the names of the SASL mechanisms for which access tokens may be generated, and for which generated access tokens will be accepted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"generate_token_result_criteria": schema.StringAttribute{
				Description: "A reference to a request criteria object that may be used to identify the types of bind operations for which access tokens may be generated. If no criteria is specified, then access tokens may be generated for any bind operations that satisfy the other requirements configured in this validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_scope": schema.SetAttribute{
				Description: "Specifies the names of any scopes that should be granted to a client that authenticates with a bind access token. By default, no scopes will be granted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Access Token Validator",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`bind`, `third-party`]: When multiple Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Access Token Validators defined within Directory Server but not necessarily contiguous. Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token. When the `type` attribute is set to `ping-federate`: When multiple Ping Federate Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Ping Federate Access Token Validators defined within Directory Server but not necessarily contiguous. Ping Federate Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token. When the `type` attribute is set to `jwt`: When multiple JWT Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all JWT Access Token Validators defined within Directory Server but not necessarily contiguous. JWT Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token. When the `type` attribute is set to `mock`: When multiple Mock Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Mock Access Token Validators defined within Directory Server but not necessarily contiguous. Mock Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`bind`, `third-party`]: When multiple Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Access Token Validators defined within Directory Server but not necessarily contiguous. Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.\n  - `ping-federate`: When multiple Ping Federate Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Ping Federate Access Token Validators defined within Directory Server but not necessarily contiguous. Ping Federate Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.\n  - `jwt`: When multiple JWT Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all JWT Access Token Validators defined within Directory Server but not necessarily contiguous. JWT Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.\n  - `mock`: When multiple Mock Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Mock Access Token Validators defined within Directory Server but not necessarily contiguous. Mock Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a BindAccessTokenValidatorResponse object into the model struct
func readBindAccessTokenValidatorResponseDataSource(ctx context.Context, r *client.BindAccessTokenValidatorResponse, state *accessTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bind")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.PersistAccessTokens = internaltypes.BoolTypeOrNil(r.PersistAccessTokens)
	state.MaximumTokenLifetime = internaltypes.StringTypeOrNil(r.MaximumTokenLifetime, false)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedAuthenticationTypeProp(r.AllowedAuthenticationType))
	state.AllowedSASLMechanism = internaltypes.GetStringSet(r.AllowedSASLMechanism)
	state.GenerateTokenResultCriteria = internaltypes.StringTypeOrNil(r.GenerateTokenResultCriteria, false)
	state.IncludedScope = internaltypes.GetStringSet(r.IncludedScope)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
}

// Read a PingFederateAccessTokenValidatorResponse object into the model struct
func readPingFederateAccessTokenValidatorResponseDataSource(ctx context.Context, r *client.PingFederateAccessTokenValidatorResponse, state *accessTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-federate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ClientID = types.StringValue(r.ClientID)
	state.ClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.ClientSecretPassphraseProvider, false)
	state.IncludeAudParameter = internaltypes.BoolTypeOrNil(r.IncludeAudParameter)
	state.AccessTokenManagerID = internaltypes.StringTypeOrNil(r.AccessTokenManagerID, false)
	state.EndpointCacheRefresh = internaltypes.StringTypeOrNil(r.EndpointCacheRefresh, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, false)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a JwtAccessTokenValidatorResponse object into the model struct
func readJwtAccessTokenValidatorResponseDataSource(ctx context.Context, r *client.JwtAccessTokenValidatorResponse, state *accessTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jwt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowedSigningAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedSigningAlgorithmProp(r.AllowedSigningAlgorithm))
	state.SigningCertificate = internaltypes.GetStringSet(r.SigningCertificate)
	state.JwksEndpointPath = internaltypes.StringTypeOrNil(r.JwksEndpointPath, false)
	state.EncryptionKeyPair = internaltypes.StringTypeOrNil(r.EncryptionKeyPair, false)
	state.AllowedKeyEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmProp(r.AllowedKeyEncryptionAlgorithm))
	state.AllowedContentEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedContentEncryptionAlgorithmProp(r.AllowedContentEncryptionAlgorithm))
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, false)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, false)
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, false)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a MockAccessTokenValidatorResponse object into the model struct
func readMockAccessTokenValidatorResponseDataSource(ctx context.Context, r *client.MockAccessTokenValidatorResponse, state *accessTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("mock")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, false)
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyAccessTokenValidatorResponse object into the model struct
func readThirdPartyAccessTokenValidatorResponseDataSource(ctx context.Context, r *client.ThirdPartyAccessTokenValidatorResponse, state *accessTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
}

// Read resource information
func (r *accessTokenValidatorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state accessTokenValidatorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessTokenValidatorAPI.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.BindAccessTokenValidatorResponse != nil {
		readBindAccessTokenValidatorResponseDataSource(ctx, readResponse.BindAccessTokenValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PingFederateAccessTokenValidatorResponse != nil {
		readPingFederateAccessTokenValidatorResponseDataSource(ctx, readResponse.PingFederateAccessTokenValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JwtAccessTokenValidatorResponse != nil {
		readJwtAccessTokenValidatorResponseDataSource(ctx, readResponse.JwtAccessTokenValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MockAccessTokenValidatorResponse != nil {
		readMockAccessTokenValidatorResponseDataSource(ctx, readResponse.MockAccessTokenValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessTokenValidatorResponse != nil {
		readThirdPartyAccessTokenValidatorResponseDataSource(ctx, readResponse.ThirdPartyAccessTokenValidatorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
