package saslmechanismhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &saslMechanismHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &saslMechanismHandlerDataSource{}
)

// Create a Sasl Mechanism Handler data source
func NewSaslMechanismHandlerDataSource() datasource.DataSource {
	return &saslMechanismHandlerDataSource{}
}

// saslMechanismHandlerDataSource is the datasource implementation.
type saslMechanismHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *saslMechanismHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the data source.
func (r *saslMechanismHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type saslMechanismHandlerDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Type                                         types.String `tfsdk:"type"`
	ExtensionClass                               types.String `tfsdk:"extension_class"`
	ExtensionArgument                            types.Set    `tfsdk:"extension_argument"`
	AccessTokenValidator                         types.Set    `tfsdk:"access_token_validator"`
	IdTokenValidator                             types.Set    `tfsdk:"id_token_validator"`
	RequireBothAccessTokenAndIDToken             types.Bool   `tfsdk:"require_both_access_token_and_id_token"`
	ValidateAccessTokenWhenIDTokenIsAlsoProvided types.String `tfsdk:"validate_access_token_when_id_token_is_also_provided"`
	AlternateAuthorizationIdentityMapper         types.String `tfsdk:"alternate_authorization_identity_mapper"`
	AllRequiredScope                             types.Set    `tfsdk:"all_required_scope"`
	AnyRequiredScope                             types.Set    `tfsdk:"any_required_scope"`
	OtpValidityDuration                          types.String `tfsdk:"otp_validity_duration"`
	ServerFqdn                                   types.String `tfsdk:"server_fqdn"`
	IdentityMapper                               types.String `tfsdk:"identity_mapper"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *saslMechanismHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of SASL Mechanism Handler resource. Options are ['unboundid-ms-chap-v2', 'unboundid-totp', 'unboundid-yubikey-otp', 'external', 'digest-md5', 'plain', 'unboundid-delivered-otp', 'unboundid-external-auth', 'anonymous', 'cram-md5', 'oauth-bearer', 'unboundid-certificate-plus-password', 'gssapi', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party SASL Mechanism Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party SASL Mechanism Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"access_token_validator": schema.SetAttribute{
				Description: "An access token validator that will ensure that each presented OAuth access token is authentic and trustworthy. It must be configured with an identity mapper that will be used to map the access token to a local entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "An ID token validator that will ensure that each presented OpenID Connect ID token is authentic and trustworthy, and that will map the token to a local entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_both_access_token_and_id_token": schema.BoolAttribute{
				Description: "Indicates whether bind requests will be required to have both an OAuth access token (in the \"auth\" element of the bind request) and an OpenID Connect ID token (in the \"pingidentityidtoken\" element of the bind request).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"validate_access_token_when_id_token_is_also_provided": schema.StringAttribute{
				Description: "Indicates whether to validate the OAuth access token in addition to the OpenID Connect ID token in OAUTHBEARER bind requests that contain both types of tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alternate_authorization_identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"all_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that will all be required for any access tokens that will be allowed for authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that a token may have to be allowed for authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"otp_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a one-time password value should be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_fqdn": schema.StringAttribute{
				Description: "Specifies the DNS-resolvable fully-qualified domain name for the server that is used when validating the digest-uri parameter during the authentication process.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that should be used to identify the entry associated with the username provided in the bind request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a UnboundidMsChapV2SaslMechanismHandlerResponse object into the model struct
func readUnboundidMsChapV2SaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidMsChapV2SaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-ms-chap-v2")
	state.Id = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a UnboundidDeliveredOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidDeliveredOtpSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidDeliveredOtpSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-delivered-otp")
	state.Id = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.OtpValidityDuration = types.StringValue(r.OtpValidityDuration)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a OauthBearerSaslMechanismHandlerResponse object into the model struct
func readOauthBearerSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.OauthBearerSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oauth-bearer")
	state.Id = types.StringValue(r.Id)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireBothAccessTokenAndIDToken = internaltypes.BoolTypeOrNil(r.RequireBothAccessTokenAndIDToken)
	state.ValidateAccessTokenWhenIDTokenIsAlsoProvided = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerValidateAccessTokenWhenIDTokenIsAlsoProvidedProp(r.ValidateAccessTokenWhenIDTokenIsAlsoProvided), false)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, false)
	state.AllRequiredScope = internaltypes.GetStringSet(r.AllRequiredScope)
	state.AnyRequiredScope = internaltypes.GetStringSet(r.AnyRequiredScope)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartySaslMechanismHandlerResponse object into the model struct
func readThirdPartySaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartySaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *saslMechanismHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state saslMechanismHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
		readUnboundidMsChapV2SaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
		readUnboundidDeliveredOtpSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OauthBearerSaslMechanismHandlerResponse != nil {
		readOauthBearerSaslMechanismHandlerResponseDataSource(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySaslMechanismHandlerResponse != nil {
		readThirdPartySaslMechanismHandlerResponseDataSource(ctx, readResponse.ThirdPartySaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
