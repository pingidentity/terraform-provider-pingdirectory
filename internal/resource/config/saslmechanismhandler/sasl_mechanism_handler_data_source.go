package saslmechanismhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
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
	r.apiClient = providerCfg.ApiClient
}

type saslMechanismHandlerDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	Type                                         types.String `tfsdk:"type"`
	ExtensionClass                               types.String `tfsdk:"extension_class"`
	ExtensionArgument                            types.Set    `tfsdk:"extension_argument"`
	AccessTokenValidator                         types.Set    `tfsdk:"access_token_validator"`
	KdcAddress                                   types.String `tfsdk:"kdc_address"`
	Keytab                                       types.String `tfsdk:"keytab"`
	AllowNullServerFqdn                          types.Bool   `tfsdk:"allow_null_server_fqdn"`
	IdTokenValidator                             types.Set    `tfsdk:"id_token_validator"`
	AllowedQualityOfProtection                   types.Set    `tfsdk:"allowed_quality_of_protection"`
	RequireBothAccessTokenAndIDToken             types.Bool   `tfsdk:"require_both_access_token_and_id_token"`
	ValidateAccessTokenWhenIDTokenIsAlsoProvided types.String `tfsdk:"validate_access_token_when_id_token_is_also_provided"`
	KerberosServicePrincipal                     types.String `tfsdk:"kerberos_service_principal"`
	GssapiRole                                   types.String `tfsdk:"gssapi_role"`
	JaasConfigFile                               types.String `tfsdk:"jaas_config_file"`
	EnableDebug                                  types.Bool   `tfsdk:"enable_debug"`
	AlternateAuthorizationIdentityMapper         types.String `tfsdk:"alternate_authorization_identity_mapper"`
	AllRequiredScope                             types.Set    `tfsdk:"all_required_scope"`
	AnyRequiredScope                             types.Set    `tfsdk:"any_required_scope"`
	Realm                                        types.String `tfsdk:"realm"`
	OtpValidityDuration                          types.String `tfsdk:"otp_validity_duration"`
	CertificateValidationPolicy                  types.String `tfsdk:"certificate_validation_policy"`
	ServerFqdn                                   types.String `tfsdk:"server_fqdn"`
	CertificateAttribute                         types.String `tfsdk:"certificate_attribute"`
	CertificateMapper                            types.String `tfsdk:"certificate_mapper"`
	YubikeyClientID                              types.String `tfsdk:"yubikey_client_id"`
	YubikeyAPIKey                                types.String `tfsdk:"yubikey_api_key"`
	YubikeyAPIKeyPassphraseProvider              types.String `tfsdk:"yubikey_api_key_passphrase_provider"`
	YubikeyValidationServerBaseURL               types.Set    `tfsdk:"yubikey_validation_server_base_url"`
	HttpProxyExternalServer                      types.String `tfsdk:"http_proxy_external_server"`
	HttpConnectTimeout                           types.String `tfsdk:"http_connect_timeout"`
	HttpResponseTimeout                          types.String `tfsdk:"http_response_timeout"`
	IdentityMapper                               types.String `tfsdk:"identity_mapper"`
	SharedSecretAttributeType                    types.String `tfsdk:"shared_secret_attribute_type"`
	KeyManagerProvider                           types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                         types.String `tfsdk:"trust_manager_provider"`
	TimeIntervalDuration                         types.String `tfsdk:"time_interval_duration"`
	AdjacentIntervalsToCheck                     types.Int64  `tfsdk:"adjacent_intervals_to_check"`
	RequireStaticPassword                        types.Bool   `tfsdk:"require_static_password"`
	PreventTOTPReuse                             types.Bool   `tfsdk:"prevent_totp_reuse"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *saslMechanismHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
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
			"kdc_address": schema.StringAttribute{
				Description: "Specifies the address of the KDC that is to be used for Kerberos processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"keytab": schema.StringAttribute{
				Description: "Specifies the keytab file that should be used for Kerberos processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_null_server_fqdn": schema.BoolAttribute{
				Description: "Specifies whether or not to allow a null value for the server-fqdn.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "An ID token validator that will ensure that each presented OpenID Connect ID token is authentic and trustworthy, and that will map the token to a local entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_quality_of_protection": schema.SetAttribute{
				Description: "Specifies the supported quality of protection (QoP) levels that clients will be permitted to request when performing GSSAPI authentication.",
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
			"kerberos_service_principal": schema.StringAttribute{
				Description: "Specifies the Kerberos service principal that the Directory Server will use to identify itself to the KDC.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"gssapi_role": schema.StringAttribute{
				Description: "Specifies the role that should be declared for the server in the generated JAAS configuration file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jaas_config_file": schema.StringAttribute{
				Description: "Specifies the path to a JAAS (Java Authentication and Authorization Service) configuration file that provides the information that the JVM should use for Kerberos processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_debug": schema.BoolAttribute{
				Description: "Indicates whether to enable debugging for the Java GSSAPI provider. Debug information will be written to standard output, which should be captured in the server.out log file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alternate_authorization_identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `oauth-bearer`: The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry. When the `type` attribute is set to `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to map the alternate authorization identity (if provided, and if different from the Kerberos principal used as the authentication identity) to the corresponding user in the directory. If no value is specified, then the mapper specified in the identity-mapper configuration property will be used.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `oauth-bearer`: The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry.\n  - `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to map the alternate authorization identity (if provided, and if different from the Kerberos principal used as the authentication identity) to the corresponding user in the directory. If no value is specified, then the mapper specified in the identity-mapper configuration property will be used.",
				Required:            false,
				Optional:            false,
				Computed:            true,
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
			"realm": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `digest-md5`: Specifies the realm that is to be used by the server for DIGEST-MD5 authentication. When the `type` attribute is set to `gssapi`: Specifies the realm to be used for GSSAPI authentication.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `digest-md5`: Specifies the realm that is to be used by the server for DIGEST-MD5 authentication.\n  - `gssapi`: Specifies the realm to be used for GSSAPI authentication.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"otp_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a one-time password value should be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"certificate_validation_policy": schema.StringAttribute{
				Description: "Indicates whether to attempt to validate the peer certificate against a certificate held in the user's entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_fqdn": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `digest-md5`: Specifies the DNS-resolvable fully-qualified domain name for the server that is used when validating the digest-uri parameter during the authentication process. When the `type` attribute is set to `oauth-bearer`: The fully-qualified name that clients are expected to use when communicating with the server. When the `type` attribute is set to `gssapi`: Specifies the DNS-resolvable fully-qualified domain name for the system.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `digest-md5`: Specifies the DNS-resolvable fully-qualified domain name for the server that is used when validating the digest-uri parameter during the authentication process.\n  - `oauth-bearer`: The fully-qualified name that clients are expected to use when communicating with the server.\n  - `gssapi`: Specifies the DNS-resolvable fully-qualified domain name for the system.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"certificate_attribute": schema.StringAttribute{
				Description: "Specifies the name of the attribute to hold user certificates.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"certificate_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `external`: Specifies the name of the certificate mapper that should be used to match client certificates to user entries. When the `type` attribute is set to `unboundid-certificate-plus-password`: The certificate mapper that will be used to identify the target user based on the certificate that was presented to the server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `external`: Specifies the name of the certificate mapper that should be used to match client certificates to user entries.\n  - `unboundid-certificate-plus-password`: The certificate mapper that will be used to identify the target user based on the certificate that was presented to the server.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"yubikey_client_id": schema.StringAttribute{
				Description: "The client ID to include in requests to the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"yubikey_api_key": schema.StringAttribute{
				Description: "The API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"yubikey_api_key_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"yubikey_validation_server_base_url": schema.SetAttribute{
				Description: "The base URL of the validation server to use to verify one-time passwords. You should only need to change the value if you wish to use your own validation server instead of using one of the Yubico servers. The server must use the YubiKey Validation Protocol version 2.0.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the YubiKey validation service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_connect_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait to obtain an HTTP connection.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_response_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait for a response to an HTTP request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`unboundid-totp`, `unboundid-yubikey-otp`, `unboundid-delivered-otp`]: The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities. When the `type` attribute is set to  one of [`digest-md5`, `plain`]: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the authentication or authorization ID included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `unboundid-ms-chap-v2`: The identity mapper that should be used to identify the entry associated with the username provided in the bind request. When the `type` attribute is set to `unboundid-external-auth`: The identity mapper that should be used to identify the user targeted by the authentication ID contained in the bind request. This will only be used for \"u:\"-style authentication ID values. When the `type` attribute is set to `cram-md5`: Specifies the name of the identity mapper used with this SASL mechanism handler to match the authentication ID included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the Kerberos principal included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `third-party`: The identity mapper that may be used to map usernames to user entries. If the custom SASL mechanism involves a username or some other form of authentication and/or authorization identity, then this may be used to map that ID to an entry for that user.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`unboundid-totp`, `unboundid-yubikey-otp`, `unboundid-delivered-otp`]: The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities.\n  - One of [`digest-md5`, `plain`]: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the authentication or authorization ID included in the SASL bind request to the corresponding user in the directory.\n  - `unboundid-ms-chap-v2`: The identity mapper that should be used to identify the entry associated with the username provided in the bind request.\n  - `unboundid-external-auth`: The identity mapper that should be used to identify the user targeted by the authentication ID contained in the bind request. This will only be used for \"u:\"-style authentication ID values.\n  - `cram-md5`: Specifies the name of the identity mapper used with this SASL mechanism handler to match the authentication ID included in the SASL bind request to the corresponding user in the directory.\n  - `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the Kerberos principal included in the SASL bind request to the corresponding user in the directory.\n  - `third-party`: The identity mapper that may be used to map usernames to user entries. If the custom SASL mechanism involves a username or some other form of authentication and/or authorization identity, then this may be used to map that ID to an entry for that user.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"shared_secret_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies which key manager provider should be used to obtain a client certificate to present to the validation server when performing HTTPS communication. This may be left undefined if communication will not be secured with HTTPS, or if there is no need to present a client certificate to the validation service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies which trust manager provider should be used to determine whether to trust the certificate presented by the server when performing HTTPS communication. This may be left undefined if HTTPS communication is not needed, or if the validation service presents a certificate that is trusted by the default JVM configuration (which should be the case for the validation servers that Yubico provides, but may not be the case if an alternate validation server is configured).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_interval_duration": schema.StringAttribute{
				Description: "The duration of the time interval used for TOTP processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"adjacent_intervals_to_check": schema.Int64Attribute{
				Description: "The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_static_password": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `unboundid-totp`: Indicates whether to require a static password (as might be held in the userPassword attribute, or whatever password attribute is defined in the password policy governing the user) in addition to the one-time password. When the `type` attribute is set to `unboundid-yubikey-otp`: Indicates whether a user will be required to provide a static password when authenticating via the UNBOUNDID-YUBIKEY-OTP SASL mechanism.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `unboundid-totp`: Indicates whether to require a static password (as might be held in the userPassword attribute, or whatever password attribute is defined in the password policy governing the user) in addition to the one-time password.\n  - `unboundid-yubikey-otp`: Indicates whether a user will be required to provide a static password when authenticating via the UNBOUNDID-YUBIKEY-OTP SASL mechanism.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"prevent_totp_reuse": schema.BoolAttribute{
				Description: "Indicates whether to prevent clients from re-using TOTP passwords.",
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
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a UnboundidMsChapV2SaslMechanismHandlerResponse object into the model struct
func readUnboundidMsChapV2SaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidMsChapV2SaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-ms-chap-v2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a UnboundidTotpSaslMechanismHandlerResponse object into the model struct
func readUnboundidTotpSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidTotpSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-totp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, false)
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, false)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.RequireStaticPassword = internaltypes.BoolTypeOrNil(r.RequireStaticPassword)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a UnboundidYubikeyOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidYubikeyOtpSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidYubikeyOtpSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-yubikey-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.YubikeyClientID = internaltypes.StringTypeOrNil(r.YubikeyClientID, false)
	state.YubikeyAPIKeyPassphraseProvider = internaltypes.StringTypeOrNil(r.YubikeyAPIKeyPassphraseProvider, false)
	state.YubikeyValidationServerBaseURL = internaltypes.GetStringSet(r.YubikeyValidationServerBaseURL)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.HttpConnectTimeout = internaltypes.StringTypeOrNil(r.HttpConnectTimeout, false)
	state.HttpResponseTimeout = internaltypes.StringTypeOrNil(r.HttpResponseTimeout, false)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.RequireStaticPassword = internaltypes.BoolTypeOrNil(r.RequireStaticPassword)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ExternalSaslMechanismHandlerResponse object into the model struct
func readExternalSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.ExternalSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("external")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CertificateValidationPolicy = types.StringValue(r.CertificateValidationPolicy.String())
	state.CertificateAttribute = internaltypes.StringTypeOrNil(r.CertificateAttribute, false)
	state.CertificateMapper = types.StringValue(r.CertificateMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a DigestMd5SaslMechanismHandlerResponse object into the model struct
func readDigestMd5SaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.DigestMd5SaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("digest-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Realm = internaltypes.StringTypeOrNil(r.Realm, false)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PlainSaslMechanismHandlerResponse object into the model struct
func readPlainSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.PlainSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("plain")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a UnboundidDeliveredOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidDeliveredOtpSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidDeliveredOtpSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-delivered-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.OtpValidityDuration = types.StringValue(r.OtpValidityDuration)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a UnboundidExternalAuthSaslMechanismHandlerResponse object into the model struct
func readUnboundidExternalAuthSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidExternalAuthSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-external-auth")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AnonymousSaslMechanismHandlerResponse object into the model struct
func readAnonymousSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.AnonymousSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("anonymous")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CramMd5SaslMechanismHandlerResponse object into the model struct
func readCramMd5SaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.CramMd5SaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("cram-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a OauthBearerSaslMechanismHandlerResponse object into the model struct
func readOauthBearerSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.OauthBearerSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oauth-bearer")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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

// Read a UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse object into the model struct
func readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-certificate-plus-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CertificateMapper = types.StringValue(r.CertificateMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GssapiSaslMechanismHandlerResponse object into the model struct
func readGssapiSaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.GssapiSaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("gssapi")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Realm = internaltypes.StringTypeOrNil(r.Realm, false)
	state.KdcAddress = internaltypes.StringTypeOrNil(r.KdcAddress, false)
	state.Keytab = internaltypes.StringTypeOrNil(r.Keytab, false)
	state.AllowNullServerFqdn = internaltypes.BoolTypeOrNil(r.AllowNullServerFqdn)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, false)
	state.AllowedQualityOfProtection = internaltypes.GetStringSet(
		client.StringSliceEnumsaslMechanismHandlerAllowedQualityOfProtectionProp(r.AllowedQualityOfProtection))
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, false)
	state.KerberosServicePrincipal = internaltypes.StringTypeOrNil(r.KerberosServicePrincipal, false)
	state.GssapiRole = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerGssapiRoleProp(r.GssapiRole), false)
	state.JaasConfigFile = internaltypes.StringTypeOrNil(r.JaasConfigFile, false)
	state.EnableDebug = internaltypes.BoolTypeOrNil(r.EnableDebug)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartySaslMechanismHandlerResponse object into the model struct
func readThirdPartySaslMechanismHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartySaslMechanismHandlerResponse, state *saslMechanismHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
	if readResponse.UnboundidTotpSaslMechanismHandlerResponse != nil {
		readUnboundidTotpSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidTotpSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
		readUnboundidYubikeyOtpSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ExternalSaslMechanismHandlerResponse != nil {
		readExternalSaslMechanismHandlerResponseDataSource(ctx, readResponse.ExternalSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DigestMd5SaslMechanismHandlerResponse != nil {
		readDigestMd5SaslMechanismHandlerResponseDataSource(ctx, readResponse.DigestMd5SaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PlainSaslMechanismHandlerResponse != nil {
		readPlainSaslMechanismHandlerResponseDataSource(ctx, readResponse.PlainSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
		readUnboundidDeliveredOtpSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
		readUnboundidExternalAuthSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AnonymousSaslMechanismHandlerResponse != nil {
		readAnonymousSaslMechanismHandlerResponseDataSource(ctx, readResponse.AnonymousSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CramMd5SaslMechanismHandlerResponse != nil {
		readCramMd5SaslMechanismHandlerResponseDataSource(ctx, readResponse.CramMd5SaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OauthBearerSaslMechanismHandlerResponse != nil {
		readOauthBearerSaslMechanismHandlerResponseDataSource(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
		readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDataSource(ctx, readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GssapiSaslMechanismHandlerResponse != nil {
		readGssapiSaslMechanismHandlerResponseDataSource(ctx, readResponse.GssapiSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySaslMechanismHandlerResponse != nil {
		readThirdPartySaslMechanismHandlerResponseDataSource(ctx, readResponse.ThirdPartySaslMechanismHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
