package saslmechanismhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &saslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &saslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &saslMechanismHandlerResource{}
	_ resource.Resource                = &defaultSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &defaultSaslMechanismHandlerResource{}
)

// Create a Sasl Mechanism Handler resource
func NewSaslMechanismHandlerResource() resource.Resource {
	return &saslMechanismHandlerResource{}
}

func NewDefaultSaslMechanismHandlerResource() resource.Resource {
	return &defaultSaslMechanismHandlerResource{}
}

// saslMechanismHandlerResource is the resource implementation.
type saslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSaslMechanismHandlerResource is the resource implementation.
type defaultSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *saslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sasl_mechanism_handler"
}

func (r *defaultSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *saslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type saslMechanismHandlerResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
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

type defaultSaslMechanismHandlerResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *saslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	saslMechanismHandlerSchema(ctx, req, resp, false)
}

func (r *defaultSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	saslMechanismHandlerSchema(ctx, req, resp, true)
}

func saslMechanismHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of SASL Mechanism Handler resource. Options are ['unboundid-ms-chap-v2', 'unboundid-totp', 'unboundid-yubikey-otp', 'external', 'digest-md5', 'plain', 'unboundid-delivered-otp', 'unboundid-external-auth', 'anonymous', 'cram-md5', 'oauth-bearer', 'unboundid-certificate-plus-password', 'gssapi', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"unboundid-ms-chap-v2", "unboundid-delivered-otp", "oauth-bearer", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party SASL Mechanism Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party SASL Mechanism Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"access_token_validator": schema.SetAttribute{
				Description: "An access token validator that will ensure that each presented OAuth access token is authentic and trustworthy. It must be configured with an identity mapper that will be used to map the access token to a local entry.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "An ID token validator that will ensure that each presented OpenID Connect ID token is authentic and trustworthy, and that will map the token to a local entry.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"require_both_access_token_and_id_token": schema.BoolAttribute{
				Description: "Indicates whether bind requests will be required to have both an OAuth access token (in the \"auth\" element of the bind request) and an OpenID Connect ID token (in the \"pingidentityidtoken\" element of the bind request).",
				Optional:    true,
				Computed:    true,
			},
			"validate_access_token_when_id_token_is_also_provided": schema.StringAttribute{
				Description: "Indicates whether to validate the OAuth access token in addition to the OpenID Connect ID token in OAUTHBEARER bind requests that contain both types of tokens.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"validate-only-the-id-token", "validate-both-tokens-but-only-map-the-id-token", "validate-and-map-both-tokens"}...),
				},
			},
			"alternate_authorization_identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `oauth-bearer`: The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry. When the `type` attribute is set to `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to map the alternate authorization identity (if provided, and if different from the Kerberos principal used as the authentication identity) to the corresponding user in the directory. If no value is specified, then the mapper specified in the identity-mapper configuration property will be used.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `oauth-bearer`: The identity mapper that will be used to map an alternate authorization identity (provided in the GS2 header of the encoded OAUTHBEARER bind request credentials) to the corresponding local entry.\n  - `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to map the alternate authorization identity (if provided, and if different from the Kerberos principal used as the authentication identity) to the corresponding user in the directory. If no value is specified, then the mapper specified in the identity-mapper configuration property will be used.",
				Optional:            true,
			},
			"all_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that will all be required for any access tokens that will be allowed for authentication.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_required_scope": schema.SetAttribute{
				Description: "The set of OAuth scopes that a token may have to be allowed for authentication.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"otp_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a one-time password value should be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_fqdn": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `digest-md5`: Specifies the DNS-resolvable fully-qualified domain name for the server that is used when validating the digest-uri parameter during the authentication process. When the `type` attribute is set to `oauth-bearer`: The fully-qualified name that clients are expected to use when communicating with the server. When the `type` attribute is set to `gssapi`: Specifies the DNS-resolvable fully-qualified domain name for the system.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `digest-md5`: Specifies the DNS-resolvable fully-qualified domain name for the server that is used when validating the digest-uri parameter during the authentication process.\n  - `oauth-bearer`: The fully-qualified name that clients are expected to use when communicating with the server.\n  - `gssapi`: Specifies the DNS-resolvable fully-qualified domain name for the system.",
				Optional:            true,
			},
			"identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`unboundid-totp`, `unboundid-yubikey-otp`, `unboundid-delivered-otp`]: The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities. When the `type` attribute is set to  one of [`digest-md5`, `plain`]: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the authentication or authorization ID included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `unboundid-ms-chap-v2`: The identity mapper that should be used to identify the entry associated with the username provided in the bind request. When the `type` attribute is set to `unboundid-external-auth`: The identity mapper that should be used to identify the user targeted by the authentication ID contained in the bind request. This will only be used for \"u:\"-style authentication ID values. When the `type` attribute is set to `cram-md5`: Specifies the name of the identity mapper used with this SASL mechanism handler to match the authentication ID included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the Kerberos principal included in the SASL bind request to the corresponding user in the directory. When the `type` attribute is set to `third-party`: The identity mapper that may be used to map usernames to user entries. If the custom SASL mechanism involves a username or some other form of authentication and/or authorization identity, then this may be used to map that ID to an entry for that user.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`unboundid-totp`, `unboundid-yubikey-otp`, `unboundid-delivered-otp`]: The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities.\n  - One of [`digest-md5`, `plain`]: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the authentication or authorization ID included in the SASL bind request to the corresponding user in the directory.\n  - `unboundid-ms-chap-v2`: The identity mapper that should be used to identify the entry associated with the username provided in the bind request.\n  - `unboundid-external-auth`: The identity mapper that should be used to identify the user targeted by the authentication ID contained in the bind request. This will only be used for \"u:\"-style authentication ID values.\n  - `cram-md5`: Specifies the name of the identity mapper used with this SASL mechanism handler to match the authentication ID included in the SASL bind request to the corresponding user in the directory.\n  - `gssapi`: Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the Kerberos principal included in the SASL bind request to the corresponding user in the directory.\n  - `third-party`: The identity mapper that may be used to map usernames to user entries. If the custom SASL mechanism involves a username or some other form of authentication and/or authorization identity, then this may be used to map that ID to an entry for that user.",
				Optional:            true,
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
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"unboundid-ms-chap-v2", "unboundid-totp", "unboundid-yubikey-otp", "external", "digest-md5", "plain", "unboundid-delivered-otp", "unboundid-external-auth", "anonymous", "cram-md5", "oauth-bearer", "unboundid-certificate-plus-password", "gssapi", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["kdc_address"] = schema.StringAttribute{
			Description: "Specifies the address of the KDC that is to be used for Kerberos processing.",
		}
		schemaDef.Attributes["keytab"] = schema.StringAttribute{
			Description: "Specifies the keytab file that should be used for Kerberos processing.",
		}
		schemaDef.Attributes["allow_null_server_fqdn"] = schema.BoolAttribute{
			Description: "Specifies whether or not to allow a null value for the server-fqdn.",
		}
		schemaDef.Attributes["allowed_quality_of_protection"] = schema.SetAttribute{
			Description: "Specifies the supported quality of protection (QoP) levels that clients will be permitted to request when performing GSSAPI authentication.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["kerberos_service_principal"] = schema.StringAttribute{
			Description: "Specifies the Kerberos service principal that the Directory Server will use to identify itself to the KDC.",
		}
		schemaDef.Attributes["gssapi_role"] = schema.StringAttribute{
			Description: "Specifies the role that should be declared for the server in the generated JAAS configuration file.",
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"acceptor", "initiator"}...),
			},
		}
		schemaDef.Attributes["jaas_config_file"] = schema.StringAttribute{
			Description: "Specifies the path to a JAAS (Java Authentication and Authorization Service) configuration file that provides the information that the JVM should use for Kerberos processing.",
		}
		schemaDef.Attributes["enable_debug"] = schema.BoolAttribute{
			Description: "Indicates whether to enable debugging for the Java GSSAPI provider. Debug information will be written to standard output, which should be captured in the server.out log file.",
		}
		schemaDef.Attributes["realm"] = schema.StringAttribute{
			Description:         "When the `type` attribute is set to `digest-md5`: Specifies the realm that is to be used by the server for DIGEST-MD5 authentication. When the `type` attribute is set to `gssapi`: Specifies the realm to be used for GSSAPI authentication.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - `digest-md5`: Specifies the realm that is to be used by the server for DIGEST-MD5 authentication.\n  - `gssapi`: Specifies the realm to be used for GSSAPI authentication.",
		}
		schemaDef.Attributes["certificate_validation_policy"] = schema.StringAttribute{
			Description: "Indicates whether to attempt to validate the peer certificate against a certificate held in the user's entry.",
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"always", "ifpresent", "never"}...),
			},
		}
		schemaDef.Attributes["certificate_attribute"] = schema.StringAttribute{
			Description: "Specifies the name of the attribute to hold user certificates.",
		}
		schemaDef.Attributes["certificate_mapper"] = schema.StringAttribute{
			Description:         "When the `type` attribute is set to `external`: Specifies the name of the certificate mapper that should be used to match client certificates to user entries. When the `type` attribute is set to `unboundid-certificate-plus-password`: The certificate mapper that will be used to identify the target user based on the certificate that was presented to the server.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - `external`: Specifies the name of the certificate mapper that should be used to match client certificates to user entries.\n  - `unboundid-certificate-plus-password`: The certificate mapper that will be used to identify the target user based on the certificate that was presented to the server.",
		}
		schemaDef.Attributes["yubikey_client_id"] = schema.StringAttribute{
			Description: "The client ID to include in requests to the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
		}
		schemaDef.Attributes["yubikey_api_key"] = schema.StringAttribute{
			Description: "The API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
			Sensitive:   true,
		}
		schemaDef.Attributes["yubikey_api_key_passphrase_provider"] = schema.StringAttribute{
			Description: "The passphrase provider to use to obtain the API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
		}
		schemaDef.Attributes["yubikey_validation_server_base_url"] = schema.SetAttribute{
			Description: "The base URL of the validation server to use to verify one-time passwords. You should only need to change the value if you wish to use your own validation server instead of using one of the Yubico servers. The server must use the YubiKey Validation Protocol version 2.0.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["http_proxy_external_server"] = schema.StringAttribute{
			Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the YubiKey validation service.",
		}
		schemaDef.Attributes["http_connect_timeout"] = schema.StringAttribute{
			Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait to obtain an HTTP connection.",
		}
		schemaDef.Attributes["http_response_timeout"] = schema.StringAttribute{
			Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait for a response to an HTTP request.",
		}
		schemaDef.Attributes["shared_secret_attribute_type"] = schema.StringAttribute{
			Description: "The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.",
		}
		schemaDef.Attributes["key_manager_provider"] = schema.StringAttribute{
			Description: "Specifies which key manager provider should be used to obtain a client certificate to present to the validation server when performing HTTPS communication. This may be left undefined if communication will not be secured with HTTPS, or if there is no need to present a client certificate to the validation service.",
		}
		schemaDef.Attributes["trust_manager_provider"] = schema.StringAttribute{
			Description: "Specifies which trust manager provider should be used to determine whether to trust the certificate presented by the server when performing HTTPS communication. This may be left undefined if HTTPS communication is not needed, or if the validation service presents a certificate that is trusted by the default JVM configuration (which should be the case for the validation servers that Yubico provides, but may not be the case if an alternate validation server is configured).",
		}
		schemaDef.Attributes["time_interval_duration"] = schema.StringAttribute{
			Description: "The duration of the time interval used for TOTP processing.",
		}
		schemaDef.Attributes["adjacent_intervals_to_check"] = schema.Int64Attribute{
			Description: "The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.",
		}
		schemaDef.Attributes["require_static_password"] = schema.BoolAttribute{
			Description:         "When the `type` attribute is set to `unboundid-totp`: Indicates whether to require a static password (as might be held in the userPassword attribute, or whatever password attribute is defined in the password policy governing the user) in addition to the one-time password. When the `type` attribute is set to `unboundid-yubikey-otp`: Indicates whether a user will be required to provide a static password when authenticating via the UNBOUNDID-YUBIKEY-OTP SASL mechanism.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - `unboundid-totp`: Indicates whether to require a static password (as might be held in the userPassword attribute, or whatever password attribute is defined in the password policy governing the user) in addition to the one-time password.\n  - `unboundid-yubikey-otp`: Indicates whether a user will be required to provide a static password when authenticating via the UNBOUNDID-YUBIKEY-OTP SASL mechanism.",
		}
		schemaDef.Attributes["prevent_totp_reuse"] = schema.BoolAttribute{
			Description: "Indicates whether to prevent clients from re-using TOTP passwords.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *saslMechanismHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
	var planModel, configModel saslMechanismHandlerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for oauth-bearer type
	if resourceType == "oauth-bearer" {
		if !internaltypes.IsDefined(configModel.RequireBothAccessTokenAndIDToken) {
			defaultVal := types.BoolValue(false)
			if !planModel.RequireBothAccessTokenAndIDToken.Equal(defaultVal) {
				planModel.RequireBothAccessTokenAndIDToken = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ValidateAccessTokenWhenIDTokenIsAlsoProvided) {
			defaultVal := types.StringValue("validate-only-the-id-token")
			if !planModel.ValidateAccessTokenWhenIDTokenIsAlsoProvided.Equal(defaultVal) {
				planModel.ValidateAccessTokenWhenIDTokenIsAlsoProvided = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultSaslMechanismHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanSaslMechanismHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model defaultSaslMechanismHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpConnectTimeout) {
		resp.Diagnostics.AddError("Attribute 'http_connect_timeout' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.HttpResponseTimeout) {
		resp.Diagnostics.AddError("Attribute 'http_response_timeout' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	compare, err = version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

func (model *saslMechanismHandlerResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "unboundid-ms-chap-v2" {
		model.RequireBothAccessTokenAndIDToken = types.BoolNull()
		model.ValidateAccessTokenWhenIDTokenIsAlsoProvided = types.StringNull()
		model.OtpValidityDuration = types.StringNull()
	}
	if resourceType == "unboundid-delivered-otp" {
		model.RequireBothAccessTokenAndIDToken = types.BoolNull()
		model.ValidateAccessTokenWhenIDTokenIsAlsoProvided = types.StringNull()
	}
	if resourceType == "oauth-bearer" {
		model.OtpValidityDuration = types.StringNull()
	}
	if resourceType == "third-party" {
		model.RequireBothAccessTokenAndIDToken = types.BoolNull()
		model.ValidateAccessTokenWhenIDTokenIsAlsoProvided = types.StringNull()
		model.OtpValidityDuration = types.StringNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsSaslMechanismHandler() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("access_token_validator"),
				path.MatchRoot("id_token_validator"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("identity_mapper"),
			path.MatchRoot("type"),
			[]string{"unboundid-ms-chap-v2", "unboundid-totp", "unboundid-yubikey-otp", "digest-md5", "plain", "unboundid-delivered-otp", "unboundid-external-auth", "cram-md5", "gssapi", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_fqdn"),
			path.MatchRoot("type"),
			[]string{"digest-md5", "oauth-bearer", "gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("otp_validity_duration"),
			path.MatchRoot("type"),
			[]string{"unboundid-delivered-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("access_token_validator"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id_token_validator"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_both_access_token_and_id_token"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("validate_access_token_when_id_token_is_also_provided"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("alternate_authorization_identity_mapper"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer", "gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_required_scope"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_required_scope"),
			path.MatchRoot("type"),
			[]string{"oauth-bearer"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"unboundid-ms-chap-v2",
			[]path.Expression{path.MatchRoot("identity_mapper"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"unboundid-delivered-otp",
			[]path.Expression{path.MatchRoot("identity_mapper"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"oauth-bearer",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r saslMechanismHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsSaslMechanismHandler()
}

// Add config validators
func (r defaultSaslMechanismHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("shared_secret_attribute_type"),
			path.MatchRoot("type"),
			[]string{"unboundid-totp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("time_interval_duration"),
			path.MatchRoot("type"),
			[]string{"unboundid-totp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("adjacent_intervals_to_check"),
			path.MatchRoot("type"),
			[]string{"unboundid-totp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_static_password"),
			path.MatchRoot("type"),
			[]string{"unboundid-totp", "unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prevent_totp_reuse"),
			path.MatchRoot("type"),
			[]string{"unboundid-totp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("yubikey_client_id"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("yubikey_api_key"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("yubikey_api_key_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("yubikey_validation_server_base_url"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_connect_timeout"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_response_timeout"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_manager_provider"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_manager_provider"),
			path.MatchRoot("type"),
			[]string{"unboundid-yubikey-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("certificate_validation_policy"),
			path.MatchRoot("type"),
			[]string{"external"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("certificate_attribute"),
			path.MatchRoot("type"),
			[]string{"external"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("certificate_mapper"),
			path.MatchRoot("type"),
			[]string{"external", "unboundid-certificate-plus-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("realm"),
			path.MatchRoot("type"),
			[]string{"digest-md5", "gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("kdc_address"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("keytab"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_null_server_fqdn"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_quality_of_protection"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("kerberos_service_principal"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("gssapi_role"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jaas_config_file"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("enable_debug"),
			path.MatchRoot("type"),
			[]string{"gssapi"},
		),
	}
	return append(configValidatorsSaslMechanismHandler(), validators...)
}

// Add optional fields to create request for unboundid-ms-chap-v2 sasl-mechanism-handler
func addOptionalUnboundidMsChapV2SaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddUnboundidMsChapV2SaslMechanismHandlerRequest, plan saslMechanismHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for unboundid-delivered-otp sasl-mechanism-handler
func addOptionalUnboundidDeliveredOtpSaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddUnboundidDeliveredOtpSaslMechanismHandlerRequest, plan saslMechanismHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OtpValidityDuration) {
		addRequest.OtpValidityDuration = plan.OtpValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for oauth-bearer sasl-mechanism-handler
func addOptionalOauthBearerSaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddOauthBearerSaslMechanismHandlerRequest, plan saslMechanismHandlerResourceModel) error {
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

// Add optional fields to create request for third-party sasl-mechanism-handler
func addOptionalThirdPartySaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddThirdPartySaslMechanismHandlerRequest, plan saslMechanismHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateSaslMechanismHandlerUnknownValues(model *saslMechanismHandlerResourceModel) {
	if model.AccessTokenValidator.IsUnknown() || model.AccessTokenValidator.IsNull() {
		model.AccessTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyRequiredScope.IsUnknown() || model.AnyRequiredScope.IsNull() {
		model.AnyRequiredScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllRequiredScope.IsUnknown() || model.AllRequiredScope.IsNull() {
		model.AllRequiredScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IdTokenValidator.IsUnknown() || model.IdTokenValidator.IsNull() {
		model.IdTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateSaslMechanismHandlerUnknownValuesDefault(model *defaultSaslMechanismHandlerResourceModel) {
	if model.AccessTokenValidator.IsUnknown() || model.AccessTokenValidator.IsNull() {
		model.AccessTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyRequiredScope.IsUnknown() || model.AnyRequiredScope.IsNull() {
		model.AnyRequiredScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedQualityOfProtection.IsUnknown() || model.AllowedQualityOfProtection.IsNull() {
		model.AllowedQualityOfProtection, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllRequiredScope.IsUnknown() || model.AllRequiredScope.IsNull() {
		model.AllRequiredScope, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.YubikeyValidationServerBaseURL.IsUnknown() || model.YubikeyValidationServerBaseURL.IsNull() {
		model.YubikeyValidationServerBaseURL, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IdTokenValidator.IsUnknown() || model.IdTokenValidator.IsNull() {
		model.IdTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Read a UnboundidMsChapV2SaslMechanismHandlerResponse object into the model struct
func readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx context.Context, r *client.UnboundidMsChapV2SaslMechanismHandlerResponse, state *saslMechanismHandlerResourceModel, expectedValues *saslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-ms-chap-v2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValues(state)
}

// Read a UnboundidMsChapV2SaslMechanismHandlerResponse object into the model struct
func readUnboundidMsChapV2SaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidMsChapV2SaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-ms-chap-v2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a UnboundidTotpSaslMechanismHandlerResponse object into the model struct
func readUnboundidTotpSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidTotpSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-totp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, true)
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval_duration",
		expectedValues.TimeIntervalDuration, state.TimeIntervalDuration, diagnostics)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.RequireStaticPassword = internaltypes.BoolTypeOrNil(r.RequireStaticPassword)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a UnboundidYubikeyOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidYubikeyOtpSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidYubikeyOtpSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-yubikey-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.YubikeyClientID = internaltypes.StringTypeOrNil(r.YubikeyClientID, true)
	state.YubikeyAPIKeyPassphraseProvider = internaltypes.StringTypeOrNil(r.YubikeyAPIKeyPassphraseProvider, true)
	state.YubikeyValidationServerBaseURL = internaltypes.GetStringSet(r.YubikeyValidationServerBaseURL)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, true)
	state.HttpConnectTimeout = internaltypes.StringTypeOrNil(r.HttpConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_connect_timeout",
		expectedValues.HttpConnectTimeout, state.HttpConnectTimeout, diagnostics)
	state.HttpResponseTimeout = internaltypes.StringTypeOrNil(r.HttpResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_response_timeout",
		expectedValues.HttpResponseTimeout, state.HttpResponseTimeout, diagnostics)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.RequireStaticPassword = internaltypes.BoolTypeOrNil(r.RequireStaticPassword)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, true)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a ExternalSaslMechanismHandlerResponse object into the model struct
func readExternalSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.ExternalSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("external")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CertificateValidationPolicy = types.StringValue(r.CertificateValidationPolicy.String())
	state.CertificateAttribute = internaltypes.StringTypeOrNil(r.CertificateAttribute, true)
	state.CertificateMapper = types.StringValue(r.CertificateMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a DigestMd5SaslMechanismHandlerResponse object into the model struct
func readDigestMd5SaslMechanismHandlerResponseDefault(ctx context.Context, r *client.DigestMd5SaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("digest-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Realm = internaltypes.StringTypeOrNil(r.Realm, true)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a PlainSaslMechanismHandlerResponse object into the model struct
func readPlainSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.PlainSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("plain")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a UnboundidDeliveredOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx context.Context, r *client.UnboundidDeliveredOtpSaslMechanismHandlerResponse, state *saslMechanismHandlerResourceModel, expectedValues *saslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-delivered-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.OtpValidityDuration = types.StringValue(r.OtpValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("otp_validity_duration",
		expectedValues.OtpValidityDuration, state.OtpValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValues(state)
}

// Read a UnboundidDeliveredOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidDeliveredOtpSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidDeliveredOtpSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-delivered-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.OtpValidityDuration = types.StringValue(r.OtpValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("otp_validity_duration",
		expectedValues.OtpValidityDuration, state.OtpValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a UnboundidExternalAuthSaslMechanismHandlerResponse object into the model struct
func readUnboundidExternalAuthSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidExternalAuthSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-external-auth")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a AnonymousSaslMechanismHandlerResponse object into the model struct
func readAnonymousSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.AnonymousSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("anonymous")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a CramMd5SaslMechanismHandlerResponse object into the model struct
func readCramMd5SaslMechanismHandlerResponseDefault(ctx context.Context, r *client.CramMd5SaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("cram-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a OauthBearerSaslMechanismHandlerResponse object into the model struct
func readOauthBearerSaslMechanismHandlerResponse(ctx context.Context, r *client.OauthBearerSaslMechanismHandlerResponse, state *saslMechanismHandlerResourceModel, expectedValues *saslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oauth-bearer")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireBothAccessTokenAndIDToken = internaltypes.BoolTypeOrNil(r.RequireBothAccessTokenAndIDToken)
	state.ValidateAccessTokenWhenIDTokenIsAlsoProvided = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerValidateAccessTokenWhenIDTokenIsAlsoProvidedProp(r.ValidateAccessTokenWhenIDTokenIsAlsoProvided), true)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, internaltypes.IsEmptyString(expectedValues.AlternateAuthorizationIdentityMapper))
	state.AllRequiredScope = internaltypes.GetStringSet(r.AllRequiredScope)
	state.AnyRequiredScope = internaltypes.GetStringSet(r.AnyRequiredScope)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, internaltypes.IsEmptyString(expectedValues.ServerFqdn))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValues(state)
}

// Read a OauthBearerSaslMechanismHandlerResponse object into the model struct
func readOauthBearerSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.OauthBearerSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oauth-bearer")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireBothAccessTokenAndIDToken = internaltypes.BoolTypeOrNil(r.RequireBothAccessTokenAndIDToken)
	state.ValidateAccessTokenWhenIDTokenIsAlsoProvided = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerValidateAccessTokenWhenIDTokenIsAlsoProvidedProp(r.ValidateAccessTokenWhenIDTokenIsAlsoProvided), true)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, true)
	state.AllRequiredScope = internaltypes.GetStringSet(r.AllRequiredScope)
	state.AnyRequiredScope = internaltypes.GetStringSet(r.AnyRequiredScope)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse object into the model struct
func readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unboundid-certificate-plus-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CertificateMapper = types.StringValue(r.CertificateMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a GssapiSaslMechanismHandlerResponse object into the model struct
func readGssapiSaslMechanismHandlerResponseDefault(ctx context.Context, r *client.GssapiSaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("gssapi")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Realm = internaltypes.StringTypeOrNil(r.Realm, true)
	state.KdcAddress = internaltypes.StringTypeOrNil(r.KdcAddress, true)
	state.Keytab = internaltypes.StringTypeOrNil(r.Keytab, true)
	state.AllowNullServerFqdn = internaltypes.BoolTypeOrNil(r.AllowNullServerFqdn)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, true)
	state.AllowedQualityOfProtection = internaltypes.GetStringSet(
		client.StringSliceEnumsaslMechanismHandlerAllowedQualityOfProtectionProp(r.AllowedQualityOfProtection))
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, true)
	state.KerberosServicePrincipal = internaltypes.StringTypeOrNil(r.KerberosServicePrincipal, true)
	state.GssapiRole = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerGssapiRoleProp(r.GssapiRole), true)
	state.JaasConfigFile = internaltypes.StringTypeOrNil(r.JaasConfigFile, true)
	state.EnableDebug = internaltypes.BoolTypeOrNil(r.EnableDebug)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Read a ThirdPartySaslMechanismHandlerResponse object into the model struct
func readThirdPartySaslMechanismHandlerResponse(ctx context.Context, r *client.ThirdPartySaslMechanismHandlerResponse, state *saslMechanismHandlerResourceModel, expectedValues *saslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValues(state)
}

// Read a ThirdPartySaslMechanismHandlerResponse object into the model struct
func readThirdPartySaslMechanismHandlerResponseDefault(ctx context.Context, r *client.ThirdPartySaslMechanismHandlerResponse, state *defaultSaslMechanismHandlerResourceModel, expectedValues *defaultSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSaslMechanismHandlerUnknownValuesDefault(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *defaultSaslMechanismHandlerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *defaultSaslMechanismHandlerResourceModel) {
	if !expectedValues.YubikeyAPIKey.IsUnknown() {
		state.YubikeyAPIKey = expectedValues.YubikeyAPIKey
	}
}

// Create any update operations necessary to make the state match the plan
func createSaslMechanismHandlerOperations(plan saslMechanismHandlerResourceModel, state saslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBothAccessTokenAndIDToken, state.RequireBothAccessTokenAndIDToken, "require-both-access-token-and-id-token")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidateAccessTokenWhenIDTokenIsAlsoProvided, state.ValidateAccessTokenWhenIDTokenIsAlsoProvided, "validate-access-token-when-id-token-is-also-provided")
	operations.AddStringOperationIfNecessary(&ops, plan.AlternateAuthorizationIdentityMapper, state.AlternateAuthorizationIdentityMapper, "alternate-authorization-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllRequiredScope, state.AllRequiredScope, "all-required-scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredScope, state.AnyRequiredScope, "any-required-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.OtpValidityDuration, state.OtpValidityDuration, "otp-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerFqdn, state.ServerFqdn, "server-fqdn")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createSaslMechanismHandlerOperationsDefault(plan defaultSaslMechanismHandlerResourceModel, state defaultSaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.KdcAddress, state.KdcAddress, "kdc-address")
	operations.AddStringOperationIfNecessary(&ops, plan.Keytab, state.Keytab, "keytab")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowNullServerFqdn, state.AllowNullServerFqdn, "allow-null-server-fqdn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedQualityOfProtection, state.AllowedQualityOfProtection, "allowed-quality-of-protection")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBothAccessTokenAndIDToken, state.RequireBothAccessTokenAndIDToken, "require-both-access-token-and-id-token")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidateAccessTokenWhenIDTokenIsAlsoProvided, state.ValidateAccessTokenWhenIDTokenIsAlsoProvided, "validate-access-token-when-id-token-is-also-provided")
	operations.AddStringOperationIfNecessary(&ops, plan.KerberosServicePrincipal, state.KerberosServicePrincipal, "kerberos-service-principal")
	operations.AddStringOperationIfNecessary(&ops, plan.GssapiRole, state.GssapiRole, "gssapi-role")
	operations.AddStringOperationIfNecessary(&ops, plan.JaasConfigFile, state.JaasConfigFile, "jaas-config-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDebug, state.EnableDebug, "enable-debug")
	operations.AddStringOperationIfNecessary(&ops, plan.AlternateAuthorizationIdentityMapper, state.AlternateAuthorizationIdentityMapper, "alternate-authorization-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllRequiredScope, state.AllRequiredScope, "all-required-scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredScope, state.AnyRequiredScope, "any-required-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Realm, state.Realm, "realm")
	operations.AddStringOperationIfNecessary(&ops, plan.OtpValidityDuration, state.OtpValidityDuration, "otp-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.CertificateValidationPolicy, state.CertificateValidationPolicy, "certificate-validation-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerFqdn, state.ServerFqdn, "server-fqdn")
	operations.AddStringOperationIfNecessary(&ops, plan.CertificateAttribute, state.CertificateAttribute, "certificate-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.CertificateMapper, state.CertificateMapper, "certificate-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyClientID, state.YubikeyClientID, "yubikey-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyAPIKey, state.YubikeyAPIKey, "yubikey-api-key")
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyAPIKeyPassphraseProvider, state.YubikeyAPIKeyPassphraseProvider, "yubikey-api-key-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.YubikeyValidationServerBaseURL, state.YubikeyValidationServerBaseURL, "yubikey-validation-server-base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpConnectTimeout, state.HttpConnectTimeout, "http-connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpResponseTimeout, state.HttpResponseTimeout, "http-response-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SharedSecretAttributeType, state.SharedSecretAttributeType, "shared-secret-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeIntervalDuration, state.TimeIntervalDuration, "time-interval-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.AdjacentIntervalsToCheck, state.AdjacentIntervalsToCheck, "adjacent-intervals-to-check")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireStaticPassword, state.RequireStaticPassword, "require-static-password")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventTOTPReuse, state.PreventTOTPReuse, "prevent-totp-reuse")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a unboundid-ms-chap-v2 sasl-mechanism-handler
func (r *saslMechanismHandlerResource) CreateUnboundidMsChapV2SaslMechanismHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan saslMechanismHandlerResourceModel) (*saslMechanismHandlerResourceModel, error) {
	addRequest := client.NewAddUnboundidMsChapV2SaslMechanismHandlerRequest([]client.EnumunboundidMsChapV2SaslMechanismHandlerSchemaUrn{client.ENUMUNBOUNDIDMSCHAPV2SASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLERUNBOUNDID_MS_CHAP_V2},
		plan.IdentityMapper.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalUnboundidMsChapV2SaslMechanismHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sasl Mechanism Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddUnboundidMsChapV2SaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sasl Mechanism Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state saslMechanismHandlerResourceModel
	readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, addResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a unboundid-delivered-otp sasl-mechanism-handler
func (r *saslMechanismHandlerResource) CreateUnboundidDeliveredOtpSaslMechanismHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan saslMechanismHandlerResourceModel) (*saslMechanismHandlerResourceModel, error) {
	addRequest := client.NewAddUnboundidDeliveredOtpSaslMechanismHandlerRequest([]client.EnumunboundidDeliveredOtpSaslMechanismHandlerSchemaUrn{client.ENUMUNBOUNDIDDELIVEREDOTPSASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLERUNBOUNDID_DELIVERED_OTP},
		plan.IdentityMapper.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalUnboundidDeliveredOtpSaslMechanismHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sasl Mechanism Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddUnboundidDeliveredOtpSaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sasl Mechanism Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state saslMechanismHandlerResourceModel
	readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, addResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a oauth-bearer sasl-mechanism-handler
func (r *saslMechanismHandlerResource) CreateOauthBearerSaslMechanismHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan saslMechanismHandlerResourceModel) (*saslMechanismHandlerResourceModel, error) {
	addRequest := client.NewAddOauthBearerSaslMechanismHandlerRequest([]client.EnumoauthBearerSaslMechanismHandlerSchemaUrn{client.ENUMOAUTHBEARERSASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLEROAUTH_BEARER},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalOauthBearerSaslMechanismHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sasl Mechanism Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddOauthBearerSaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sasl Mechanism Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state saslMechanismHandlerResourceModel
	readOauthBearerSaslMechanismHandlerResponse(ctx, addResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party sasl-mechanism-handler
func (r *saslMechanismHandlerResource) CreateThirdPartySaslMechanismHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan saslMechanismHandlerResourceModel) (*saslMechanismHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartySaslMechanismHandlerRequest([]client.EnumthirdPartySaslMechanismHandlerSchemaUrn{client.ENUMTHIRDPARTYSASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalThirdPartySaslMechanismHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sasl Mechanism Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddThirdPartySaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sasl Mechanism Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state saslMechanismHandlerResourceModel
	readThirdPartySaslMechanismHandlerResponse(ctx, addResponse.ThirdPartySaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *saslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan saslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *saslMechanismHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "unboundid-ms-chap-v2" {
		state, err = r.CreateUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "unboundid-delivered-otp" {
		state, err = r.CreateUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "oauth-bearer" {
		state, err = r.CreateOauthBearerSaslMechanismHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartySaslMechanismHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

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
func (r *defaultSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultSaslMechanismHandlerResourceModel
	if readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
		readUnboundidMsChapV2SaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidTotpSaslMechanismHandlerResponse != nil {
		readUnboundidTotpSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidTotpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
		readUnboundidYubikeyOtpSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExternalSaslMechanismHandlerResponse != nil {
		readExternalSaslMechanismHandlerResponseDefault(ctx, readResponse.ExternalSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DigestMd5SaslMechanismHandlerResponse != nil {
		readDigestMd5SaslMechanismHandlerResponseDefault(ctx, readResponse.DigestMd5SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PlainSaslMechanismHandlerResponse != nil {
		readPlainSaslMechanismHandlerResponseDefault(ctx, readResponse.PlainSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
		readUnboundidDeliveredOtpSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
		readUnboundidExternalAuthSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AnonymousSaslMechanismHandlerResponse != nil {
		readAnonymousSaslMechanismHandlerResponseDefault(ctx, readResponse.AnonymousSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CramMd5SaslMechanismHandlerResponse != nil {
		readCramMd5SaslMechanismHandlerResponseDefault(ctx, readResponse.CramMd5SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OauthBearerSaslMechanismHandlerResponse != nil {
		readOauthBearerSaslMechanismHandlerResponseDefault(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
		readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GssapiSaslMechanismHandlerResponse != nil {
		readGssapiSaslMechanismHandlerResponseDefault(ctx, readResponse.GssapiSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySaslMechanismHandlerResponse != nil {
		readThirdPartySaslMechanismHandlerResponseDefault(ctx, readResponse.ThirdPartySaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createSaslMechanismHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
			readUnboundidMsChapV2SaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidTotpSaslMechanismHandlerResponse != nil {
			readUnboundidTotpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidTotpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
			readUnboundidYubikeyOtpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExternalSaslMechanismHandlerResponse != nil {
			readExternalSaslMechanismHandlerResponseDefault(ctx, updateResponse.ExternalSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DigestMd5SaslMechanismHandlerResponse != nil {
			readDigestMd5SaslMechanismHandlerResponseDefault(ctx, updateResponse.DigestMd5SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PlainSaslMechanismHandlerResponse != nil {
			readPlainSaslMechanismHandlerResponseDefault(ctx, updateResponse.PlainSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
			readUnboundidDeliveredOtpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
			readUnboundidExternalAuthSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidExternalAuthSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AnonymousSaslMechanismHandlerResponse != nil {
			readAnonymousSaslMechanismHandlerResponseDefault(ctx, updateResponse.AnonymousSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CramMd5SaslMechanismHandlerResponse != nil {
			readCramMd5SaslMechanismHandlerResponseDefault(ctx, updateResponse.CramMd5SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OauthBearerSaslMechanismHandlerResponse != nil {
			readOauthBearerSaslMechanismHandlerResponseDefault(ctx, updateResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
			readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GssapiSaslMechanismHandlerResponse != nil {
			readGssapiSaslMechanismHandlerResponseDefault(ctx, updateResponse.GssapiSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySaslMechanismHandlerResponse != nil {
			readThirdPartySaslMechanismHandlerResponseDefault(ctx, updateResponse.ThirdPartySaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *saslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state saslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Sasl Mechanism Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sasl Mechanism Handler", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
		readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
		readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OauthBearerSaslMechanismHandlerResponse != nil {
		readOauthBearerSaslMechanismHandlerResponse(ctx, readResponse.OauthBearerSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySaslMechanismHandlerResponse != nil {
		readThirdPartySaslMechanismHandlerResponse(ctx, readResponse.ThirdPartySaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
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
	if readResponse.UnboundidTotpSaslMechanismHandlerResponse != nil {
		readUnboundidTotpSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidTotpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
		readUnboundidYubikeyOtpSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExternalSaslMechanismHandlerResponse != nil {
		readExternalSaslMechanismHandlerResponseDefault(ctx, readResponse.ExternalSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DigestMd5SaslMechanismHandlerResponse != nil {
		readDigestMd5SaslMechanismHandlerResponseDefault(ctx, readResponse.DigestMd5SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PlainSaslMechanismHandlerResponse != nil {
		readPlainSaslMechanismHandlerResponseDefault(ctx, readResponse.PlainSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
		readUnboundidExternalAuthSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidExternalAuthSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AnonymousSaslMechanismHandlerResponse != nil {
		readAnonymousSaslMechanismHandlerResponseDefault(ctx, readResponse.AnonymousSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CramMd5SaslMechanismHandlerResponse != nil {
		readCramMd5SaslMechanismHandlerResponseDefault(ctx, readResponse.CramMd5SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
		readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDefault(ctx, readResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GssapiSaslMechanismHandlerResponse != nil {
		readGssapiSaslMechanismHandlerResponseDefault(ctx, readResponse.GssapiSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *saslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan saslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state saslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
			readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
			readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OauthBearerSaslMechanismHandlerResponse != nil {
			readOauthBearerSaslMechanismHandlerResponse(ctx, updateResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySaslMechanismHandlerResponse != nil {
			readThirdPartySaslMechanismHandlerResponse(ctx, updateResponse.ThirdPartySaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultSaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSaslMechanismHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
			readUnboundidMsChapV2SaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidTotpSaslMechanismHandlerResponse != nil {
			readUnboundidTotpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidTotpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
			readUnboundidYubikeyOtpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExternalSaslMechanismHandlerResponse != nil {
			readExternalSaslMechanismHandlerResponseDefault(ctx, updateResponse.ExternalSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DigestMd5SaslMechanismHandlerResponse != nil {
			readDigestMd5SaslMechanismHandlerResponseDefault(ctx, updateResponse.DigestMd5SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PlainSaslMechanismHandlerResponse != nil {
			readPlainSaslMechanismHandlerResponseDefault(ctx, updateResponse.PlainSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
			readUnboundidDeliveredOtpSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
			readUnboundidExternalAuthSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidExternalAuthSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AnonymousSaslMechanismHandlerResponse != nil {
			readAnonymousSaslMechanismHandlerResponseDefault(ctx, updateResponse.AnonymousSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CramMd5SaslMechanismHandlerResponse != nil {
			readCramMd5SaslMechanismHandlerResponseDefault(ctx, updateResponse.CramMd5SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OauthBearerSaslMechanismHandlerResponse != nil {
			readOauthBearerSaslMechanismHandlerResponseDefault(ctx, updateResponse.OauthBearerSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
			readUnboundidCertificatePlusPasswordSaslMechanismHandlerResponseDefault(ctx, updateResponse.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GssapiSaslMechanismHandlerResponse != nil {
			readGssapiSaslMechanismHandlerResponseDefault(ctx, updateResponse.GssapiSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySaslMechanismHandlerResponse != nil {
			readThirdPartySaslMechanismHandlerResponseDefault(ctx, updateResponse.ThirdPartySaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *saslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state saslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SaslMechanismHandlerAPI.DeleteSaslMechanismHandlerExecute(r.apiClient.SaslMechanismHandlerAPI.DeleteSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Sasl Mechanism Handler", err, httpResp)
		return
	}
}

func (r *saslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSaslMechanismHandler(ctx, req, resp)
}

func (r *defaultSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSaslMechanismHandler(ctx, req, resp)
}

func importSaslMechanismHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
