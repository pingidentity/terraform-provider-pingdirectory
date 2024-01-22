package externalserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	_ resource.Resource                = &externalServerResource{}
	_ resource.ResourceWithConfigure   = &externalServerResource{}
	_ resource.ResourceWithImportState = &externalServerResource{}
	_ resource.Resource                = &defaultExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultExternalServerResource{}
)

// Create a External Server resource
func NewExternalServerResource() resource.Resource {
	return &externalServerResource{}
}

func NewDefaultExternalServerResource() resource.Resource {
	return &defaultExternalServerResource{}
}

// externalServerResource is the resource implementation.
type externalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultExternalServerResource is the resource implementation.
type defaultExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *externalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_server"
}

func (r *defaultExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *externalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type externalServerResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Name                                   types.String `tfsdk:"name"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	Type                                   types.String `tfsdk:"type"`
	VaultServerBaseURI                     types.Set    `tfsdk:"vault_server_base_uri"`
	VaultAuthenticationMethod              types.String `tfsdk:"vault_authentication_method"`
	HttpProxyExternalServer                types.String `tfsdk:"http_proxy_external_server"`
	ConjurServerBaseURI                    types.Set    `tfsdk:"conjur_server_base_uri"`
	AwsAccessKeyID                         types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKey                     types.String `tfsdk:"aws_secret_access_key"`
	AwsRegionName                          types.String `tfsdk:"aws_region_name"`
	ConjurAuthenticationMethod             types.String `tfsdk:"conjur_authentication_method"`
	ConjurAccountName                      types.String `tfsdk:"conjur_account_name"`
	HttpConnectTimeout                     types.String `tfsdk:"http_connect_timeout"`
	HttpResponseTimeout                    types.String `tfsdk:"http_response_timeout"`
	TrustStoreFile                         types.String `tfsdk:"trust_store_file"`
	TrustStorePin                          types.String `tfsdk:"trust_store_pin"`
	TrustStoreType                         types.String `tfsdk:"trust_store_type"`
	BaseURL                                types.String `tfsdk:"base_url"`
	HostnameVerificationMethod             types.String `tfsdk:"hostname_verification_method"`
	JdbcDriverType                         types.String `tfsdk:"jdbc_driver_type"`
	JdbcDriverURL                          types.String `tfsdk:"jdbc_driver_url"`
	SslCertNickname                        types.String `tfsdk:"ssl_cert_nickname"`
	ResponseTimeout                        types.String `tfsdk:"response_timeout"`
	BasicAuthenticationUsername            types.String `tfsdk:"basic_authentication_username"`
	BasicAuthenticationPassphraseProvider  types.String `tfsdk:"basic_authentication_passphrase_provider"`
	TransportMechanism                     types.String `tfsdk:"transport_mechanism"`
	DatabaseName                           types.String `tfsdk:"database_name"`
	VerifyCredentialsMethod                types.String `tfsdk:"verify_credentials_method"`
	UseAdministrativeOperationControl      types.Bool   `tfsdk:"use_administrative_operation_control"`
	ServerHostName                         types.String `tfsdk:"server_host_name"`
	ServerPort                             types.Int64  `tfsdk:"server_port"`
	Location                               types.String `tfsdk:"location"`
	ValidationQuery                        types.String `tfsdk:"validation_query"`
	ValidationQueryTimeout                 types.String `tfsdk:"validation_query_timeout"`
	JdbcConnectionProperties               types.Set    `tfsdk:"jdbc_connection_properties"`
	TransactionIsolationLevel              types.String `tfsdk:"transaction_isolation_level"`
	BindDN                                 types.String `tfsdk:"bind_dn"`
	SmtpSecurity                           types.String `tfsdk:"smtp_security"`
	UserName                               types.String `tfsdk:"user_name"`
	ConnectionSecurity                     types.String `tfsdk:"connection_security"`
	AuthenticationMethod                   types.String `tfsdk:"authentication_method"`
	HealthCheckConnectTimeout              types.String `tfsdk:"health_check_connect_timeout"`
	MaxConnectionAge                       types.String `tfsdk:"max_connection_age"`
	MinExpiredConnectionDisconnectInterval types.String `tfsdk:"min_expired_connection_disconnect_interval"`
	ConnectTimeout                         types.String `tfsdk:"connect_timeout"`
	MaxResponseSize                        types.String `tfsdk:"max_response_size"`
	KeyManagerProvider                     types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                   types.String `tfsdk:"trust_manager_provider"`
	InitialConnections                     types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                         types.Int64  `tfsdk:"max_connections"`
	DefunctConnectionResultCode            types.Set    `tfsdk:"defunct_connection_result_code"`
	AbandonOnTimeout                       types.Bool   `tfsdk:"abandon_on_timeout"`
	Password                               types.String `tfsdk:"password"`
	PassphraseProvider                     types.String `tfsdk:"passphrase_provider"`
	SmtpTimeout                            types.String `tfsdk:"smtp_timeout"`
	SmtpConnectionProperties               types.Set    `tfsdk:"smtp_connection_properties"`
	Description                            types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *externalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	externalServerSchema(ctx, req, resp, false)
}

func (r *defaultExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	externalServerSchema(ctx, req, resp, true)
}

func externalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a External Server.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of External Server resource. Options are ['smtp', 'nokia-ds', 'ping-identity-ds', 'active-directory', 'jdbc', 'syslog', 'ping-identity-proxy-server', 'http-proxy', 'nokia-proxy-server', 'opendj', 'ldap', 'ping-one-http', 'http', 'oracle-unified-directory', 'conjur', 'amazon-aws', 'vault']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "syslog", "ping-identity-proxy-server", "http-proxy", "nokia-proxy-server", "opendj", "ldap", "ping-one-http", "http", "oracle-unified-directory", "conjur", "amazon-aws", "vault"}...),
				},
			},
			"vault_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the Vault server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://vault.example.com:8200/\".",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Vault server.",
				Optional:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the AWS service.",
				Optional:    true,
			},
			"conjur_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the CyberArk Conjur server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://conjur.example.com:8443/\".",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_access_key_id": schema.StringAttribute{
				Description: "The access key ID that will be used if authentication should use an access key. If this is provided, then an aws-secret-access-key must also be provided.",
				Optional:    true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Description: "The secret access key that will be used if authentication should use an access key. If this is provided, then an aws-access-key-id must also be provided.",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_region_name": schema.StringAttribute{
				Description: "The name of the AWS region containing the resources that will be accessed.",
				Optional:    true,
			},
			"conjur_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Conjur server.",
				Optional:    true,
			},
			"conjur_account_name": schema.StringAttribute{
				Description: "The name of the account with which the desired secrets are associated.",
				Optional:    true,
			},
			"http_connect_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait to obtain an HTTP connection.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_response_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait for a response to an HTTP request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `conjur`: The path to a file containing the information needed to trust the certificate presented by the Conjur servers. When the `type` attribute is set to `vault`: The path to a file containing the information needed to trust the certificate presented by the Vault servers.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `conjur`: The path to a file containing the information needed to trust the certificate presented by the Conjur servers.\n  - `vault`: The path to a file containing the information needed to trust the certificate presented by the Vault servers.",
				Optional:            true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `conjur`: The PIN needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents. When the `type` attribute is set to `vault`: The passphrase needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `conjur`: The PIN needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.\n  - `vault`: The passphrase needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:            true,
				Sensitive:           true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\", \"PKCS12\", or \"BCFKS\".",
				Optional:    true,
				Computed:    true,
			},
			"base_url": schema.StringAttribute{
				Description: "The base URL of the external server, optionally including port number, for example \"https://externalService:9031\".",
				Optional:    true,
			},
			"hostname_verification_method": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ping-one-http`: The mechanism for checking if the hostname in the PingOne ID Token Validator's base-url value matches the name(s) stored inside the X.509 certificate presented by PingOne. When the `type` attribute is set to `http`: The mechanism for checking if the hostname of the HTTP External Server matches the name(s) stored inside the server's X.509 certificate. This is only applicable if SSL is being used for connection security.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ping-one-http`: The mechanism for checking if the hostname in the PingOne ID Token Validator's base-url value matches the name(s) stored inside the X.509 certificate presented by PingOne.\n  - `http`: The mechanism for checking if the hostname of the HTTP External Server matches the name(s) stored inside the server's X.509 certificate. This is only applicable if SSL is being used for connection security.",
				Optional:            true,
				Computed:            true,
			},
			"jdbc_driver_type": schema.StringAttribute{
				Description: "Specifies a supported database driver type. The driver class will be automatically selected based on this selection. We highly recommend using a JDBC 4 driver that is suitable for the current Java platform.",
				Optional:    true,
			},
			"jdbc_driver_url": schema.StringAttribute{
				Description: "Specify the complete JDBC URL which will be used instead of the automatic URL format. You must select type 'other' for the jdbc-driver-type.",
				Optional:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The certificate alias within the keystore to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Optional:    true,
			},
			"response_timeout": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ping-one-http`: Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to PingOne. When the `type` attribute is set to `http`: Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to the server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ping-one-http`: Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to PingOne.\n  - `http`: Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to the server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"basic_authentication_username": schema.StringAttribute{
				Description: "The username to use to authenticate to the HTTP Proxy External Server.",
				Optional:    true,
			},
			"basic_authentication_passphrase_provider": schema.StringAttribute{
				Description: "A passphrase provider that provides access to the password to use to authenticate to the HTTP Proxy External Server.",
				Optional:    true,
			},
			"transport_mechanism": schema.StringAttribute{
				Description: "The transport mechanism that should be used when communicating with the syslog server.",
				Optional:    true,
			},
			"database_name": schema.StringAttribute{
				Description: "Specifies which database to connect to. This is ignored if jdbc-driver-url is specified.",
				Optional:    true,
			},
			"verify_credentials_method": schema.StringAttribute{
				Description: "The mechanism to use to verify user credentials while ensuring that the ability to process other operations is not impacted by an alternate authorization identity.",
				Optional:    true,
				Computed:    true,
			},
			"use_administrative_operation_control": schema.BoolAttribute{
				Description: "Indicates whether to include the administrative operation request control in requests sent to this server which are intended for administrative operations (e.g., health checking) rather than requests directly from clients.",
				Optional:    true,
				Computed:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The host name or IP address of the target LDAP server. When the `type` attribute is set to `smtp`: The host name of the smtp server. When the `type` attribute is set to `jdbc`: The host name of the database server. This is ignored if jdbc-driver-url is specified. When the `type` attribute is set to `syslog`: The address of the syslog server. When the `type` attribute is set to `http-proxy`: The host name or IP address of the HTTP Proxy External Server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The host name or IP address of the target LDAP server.\n  - `smtp`: The host name of the smtp server.\n  - `jdbc`: The host name of the database server. This is ignored if jdbc-driver-url is specified.\n  - `syslog`: The address of the syslog server.\n  - `http-proxy`: The host name or IP address of the HTTP Proxy External Server.",
				Optional:            true,
			},
			"server_port": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The port number on which the server listens for requests. When the `type` attribute is set to `smtp`: The port number where the smtp server listens for requests. When the `type` attribute is set to `jdbc`: The port number where the database server listens for requests. This is ignored if jdbc-driver-url is specified When the `type` attribute is set to `syslog`: The port on which the syslog server accepts connections. When the `type` attribute is set to `http-proxy`: The port on which the HTTP Proxy External Server is listening for connections.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The port number on which the server listens for requests.\n  - `smtp`: The port number where the smtp server listens for requests.\n  - `jdbc`: The port number where the database server listens for requests. This is ignored if jdbc-driver-url is specified\n  - `syslog`: The port on which the syslog server accepts connections.\n  - `http-proxy`: The port on which the HTTP Proxy External Server is listening for connections.",
				Optional:            true,
				Computed:            true,
			},
			"location": schema.StringAttribute{
				Description: "Specifies the location for the LDAP External Server.",
				Optional:    true,
			},
			"validation_query": schema.StringAttribute{
				Description: "The SQL query that will be used to validate connections to the database before making them available to the Directory Server.",
				Optional:    true,
			},
			"validation_query_timeout": schema.StringAttribute{
				Description: "Specifies the amount of time to wait for a response from the database when executing the validation query, if one is set. If the timeout is exceeded, the Directory Server will drop the connection and obtain a new one. A value of zero indicates no timeout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jdbc_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the JDBC datasource.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"transaction_isolation_level": schema.StringAttribute{
				Description: "This property specifies the default transaction isolation level for connections to this JDBC External Server.",
				Optional:    true,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The DN to use to bind to the target LDAP server if simple authentication is required. When the `type` attribute is set to `active-directory`: The DN to use to bind to the target LDAP server if simple authentication is required. The authentication identity can also be specified in User-Principal-Name (UPN) format.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The DN to use to bind to the target LDAP server if simple authentication is required.\n  - `active-directory`: The DN to use to bind to the target LDAP server if simple authentication is required. The authentication identity can also be specified in User-Principal-Name (UPN) format.",
				Optional:            true,
			},
			"smtp_security": schema.StringAttribute{
				Description: "This property specifies type of connection security to use when connecting to the outgoing mail server.",
				Optional:    true,
				Computed:    true,
			},
			"user_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `smtp`: The name of the login account to use when connecting to the smtp server. Both username and password must be supplied if this attribute is set. When the `type` attribute is set to `jdbc`: The name of the login account to use when connecting to the database server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `smtp`: The name of the login account to use when connecting to the smtp server. Both username and password must be supplied if this attribute is set.\n  - `jdbc`: The name of the login account to use when connecting to the database server.",
				Optional:            true,
			},
			"connection_security": schema.StringAttribute{
				Description: "The mechanism to use to secure communication with the directory server.",
				Optional:    true,
				Computed:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The mechanism to use to authenticate to the target server. When the `type` attribute is set to `amazon-aws`: The mechanism to use to authenticate to AWS.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The mechanism to use to authenticate to the target server.\n  - `amazon-aws`: The mechanism to use to authenticate to AWS.",
				Optional:            true,
				Computed:            true,
			},
			"health_check_connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established for the purpose of performing a health check. If the connection cannot be established within this length of time, the server will be classified as unavailable.",
				Optional:    true,
			},
			"max_connection_age": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: Specifies the maximum length of time that connections to this server should be allowed to remain established before being closed and replaced with newly-established connections. When the `type` attribute is set to `syslog`: The maximum length of time that TCP connections should remain established. This will be ignored for UDP-based connections. A zero duration indicates that no maximum age will be imposed.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: Specifies the maximum length of time that connections to this server should be allowed to remain established before being closed and replaced with newly-established connections.\n  - `syslog`: The maximum length of time that TCP connections should remain established. This will be ignored for UDP-based connections. A zero duration indicates that no maximum age will be imposed.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"min_expired_connection_disconnect_interval": schema.StringAttribute{
				Description: "Specifies the minimum length of time that should pass between connection closures as a result of the connections being established for longer than the maximum connection age. This may help avoid cases in which a large number of connections are closed and re-established in a short period of time because of the maximum connection age.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connect_timeout": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable. When the `type` attribute is set to `syslog`: Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable. This will only be used when communicating with the syslog server over TCP (with or without TLS encryption). When the `type` attribute is set to `ping-one-http`: Specifies the maximum length of time to wait for a connection to be established before aborting a request to PingOne. When the `type` attribute is set to `http`: Specifies the maximum length of time to wait for a connection to be established before aborting a request to the server.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable.\n  - `syslog`: Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable. This will only be used when communicating with the syslog server over TCP (with or without TLS encryption).\n  - `ping-one-http`: Specifies the maximum length of time to wait for a connection to be established before aborting a request to PingOne.\n  - `http`: Specifies the maximum length of time to wait for a connection to be established before aborting a request to the server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_response_size": schema.StringAttribute{
				Description: "Specifies the maximum response size that should be supported for messages received from the LDAP external server.",
				Optional:    true,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The key manager provider to use if SSL or StartTLS is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store. When the `type` attribute is set to `http`: The key manager provider to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The key manager provider to use if SSL or StartTLS is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.\n  - `http`: The key manager provider to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Optional:            true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The trust manager provider to use if SSL or StartTLS is to be used for connection-level security. When the `type` attribute is set to `syslog`: A trust manager provider that will be used to determine whether to trust the certificate chain presented by the syslog server when communication is encrypted with TLS. This property will be ignored when not using TLS encryption. When the `type` attribute is set to `ping-one-http`: The trust manager provider to use for HTTPS connection-level security. When the `type` attribute is set to `http`: The trust manager provider to use if SSL (HTTPS) is to be used for connection-level security.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The trust manager provider to use if SSL or StartTLS is to be used for connection-level security.\n  - `syslog`: A trust manager provider that will be used to determine whether to trust the certificate chain presented by the syslog server when communication is encrypted with TLS. This property will be ignored when not using TLS encryption.\n  - `ping-one-http`: The trust manager provider to use for HTTPS connection-level security.\n  - `http`: The trust manager provider to use if SSL (HTTPS) is to be used for connection-level security.",
				Optional:            true,
				Computed:            true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "The number of connections to initially establish to the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Optional:    true,
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "The maximum number of concurrent connections to maintain for the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Optional:    true,
				Computed:    true,
			},
			"defunct_connection_result_code": schema.SetAttribute{
				Description: "Specifies the operation result code values that should cause the associated connection should be considered defunct. If an operation fails with one of these result codes, then it will be terminated and an attempt will be made to establish a new connection in its place.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"abandon_on_timeout": schema.BoolAttribute{
				Description: "Indicates whether to send an abandon request for an operation for which a response timeout is encountered. A request which has timed out on one server may be retried on another server regardless of whether an abandon request is sent, but if the initial attempt is not abandoned then a long-running operation may unnecessarily continue to consume processing resources on the initial server.",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The login password for the specified user. When the `type` attribute is set to `smtp`: The login password for the specified user name. Both username and password must be supplied if this attribute is set. When the `type` attribute is set to `jdbc`: The login password for the specified user name.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`nokia-ds`, `ping-identity-ds`, `active-directory`, `ping-identity-proxy-server`, `nokia-proxy-server`, `opendj`, `ldap`, `oracle-unified-directory`]: The login password for the specified user.\n  - `smtp`: The login password for the specified user name. Both username and password must be supplied if this attribute is set.\n  - `jdbc`: The login password for the specified user name.",
				Optional:            true,
				Sensitive:           true,
			},
			"passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the login password for the specified user.",
				Optional:    true,
			},
			"smtp_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection or attempted connection to a SMTP server may take.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the smtp server.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		serverHostNameAttr := schemaDef.Attributes["server_host_name"].(schema.StringAttribute)
		serverHostNameAttr.PlanModifiers = append(serverHostNameAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["server_host_name"] = serverHostNameAttr
		serverPortAttr := schemaDef.Attributes["server_port"].(schema.Int64Attribute)
		serverPortAttr.PlanModifiers = append(serverPortAttr.PlanModifiers, int64planmodifier.RequiresReplace())
		schemaDef.Attributes["server_port"] = serverPortAttr
		transportMechanismAttr := schemaDef.Attributes["transport_mechanism"].(schema.StringAttribute)
		transportMechanismAttr.PlanModifiers = append(transportMechanismAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["transport_mechanism"] = transportMechanismAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *externalServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanExternalServer(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_external_server")
	var planModel, configModel externalServerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for smtp type
	if resourceType == "smtp" {
		if !internaltypes.IsDefined(configModel.SmtpSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.SmtpSecurity.Equal(defaultVal) {
				planModel.SmtpSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for nokia-ds type
	if resourceType == "nokia-ds" {
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ping-identity-ds type
	if resourceType == "ping-identity-ds" {
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("retain-identity-control")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseAdministrativeOperationControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.UseAdministrativeOperationControl.Equal(defaultVal) {
				planModel.UseAdministrativeOperationControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for active-directory type
	if resourceType == "active-directory" {
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("separate-connections")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for jdbc type
	if resourceType == "jdbc" {
		if !internaltypes.IsDefined(configModel.TransactionIsolationLevel) {
			defaultVal := types.StringValue("read-committed")
			if !planModel.TransactionIsolationLevel.Equal(defaultVal) {
				planModel.TransactionIsolationLevel = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog type
	if resourceType == "syslog" {
		if !internaltypes.IsDefined(configModel.TrustManagerProvider) {
			defaultVal := types.StringValue("JVM-Default")
			if !planModel.TrustManagerProvider.Equal(defaultVal) {
				planModel.TrustManagerProvider = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ping-identity-proxy-server type
	if resourceType == "ping-identity-proxy-server" {
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("retain-identity-control")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseAdministrativeOperationControl) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseAdministrativeOperationControl.Equal(defaultVal) {
				planModel.UseAdministrativeOperationControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for nokia-proxy-server type
	if resourceType == "nokia-proxy-server" {
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("retain-identity-control")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseAdministrativeOperationControl) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseAdministrativeOperationControl.Equal(defaultVal) {
				planModel.UseAdministrativeOperationControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for opendj type
	if resourceType == "opendj" {
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("separate-connections")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ldap type
	if resourceType == "ldap" {
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("separate-connections")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ping-one-http type
	if resourceType == "ping-one-http" {
		if !internaltypes.IsDefined(configModel.HostnameVerificationMethod) {
			defaultVal := types.StringValue("strict")
			if !planModel.HostnameVerificationMethod.Equal(defaultVal) {
				planModel.HostnameVerificationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for http type
	if resourceType == "http" {
		if !internaltypes.IsDefined(configModel.HostnameVerificationMethod) {
			defaultVal := types.StringValue("strict")
			if !planModel.HostnameVerificationMethod.Equal(defaultVal) {
				planModel.HostnameVerificationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for oracle-unified-directory type
	if resourceType == "oracle-unified-directory" {
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(389)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ConnectionSecurity) {
			defaultVal := types.StringValue("none")
			if !planModel.ConnectionSecurity.Equal(defaultVal) {
				planModel.ConnectionSecurity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AuthenticationMethod) {
			defaultVal := types.StringValue("simple")
			if !planModel.AuthenticationMethod.Equal(defaultVal) {
				planModel.AuthenticationMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.VerifyCredentialsMethod) {
			defaultVal := types.StringValue("separate-connections")
			if !planModel.VerifyCredentialsMethod.Equal(defaultVal) {
				planModel.VerifyCredentialsMethod = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxResponseSize) {
			defaultVal := types.StringValue("10 mb")
			if !planModel.MaxResponseSize.Equal(defaultVal) {
				planModel.MaxResponseSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InitialConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.InitialConnections.Equal(defaultVal) {
				planModel.InitialConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxConnections) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxConnections.Equal(defaultVal) {
				planModel.MaxConnections = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefunctConnectionResultCode) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("operations-error"), types.StringValue("protocol-error"), types.StringValue("busy"), types.StringValue("unavailable"), types.StringValue("unwilling-to-perform"), types.StringValue("other"), types.StringValue("server-down"), types.StringValue("local-error"), types.StringValue("encoding-error"), types.StringValue("decoding-error"), types.StringValue("no-memory"), types.StringValue("connect-error"), types.StringValue("timeout")})
			if !planModel.DefunctConnectionResultCode.Equal(defaultVal) {
				planModel.DefunctConnectionResultCode = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AbandonOnTimeout) {
			defaultVal := types.BoolValue(true)
			if !planModel.AbandonOnTimeout.Equal(defaultVal) {
				planModel.AbandonOnTimeout = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for conjur type
	if resourceType == "conjur" {
		if !internaltypes.IsDefined(configModel.TrustStoreType) {
			defaultVal := types.StringValue("JKS")
			if !planModel.TrustStoreType.Equal(defaultVal) {
				planModel.TrustStoreType = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for vault type
	if resourceType == "vault" {
		if !internaltypes.IsDefined(configModel.TrustStoreType) {
			defaultVal := types.StringValue("JKS")
			if !planModel.TrustStoreType.Equal(defaultVal) {
				planModel.TrustStoreType = defaultVal
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

func (r *defaultExternalServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanExternalServer(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_external_server")
}

func modifyPlanExternalServer(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model externalServerResourceModel
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
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "http-proxy" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
			providerConfig.ProductVersion, resourceName+" with type \"http_proxy\"")
	}
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

func (model *externalServerResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "smtp" {
		model.MaxConnections = types.Int64Null()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "nokia-ds" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "ping-identity-ds" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "active-directory" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "jdbc" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "syslog" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "ping-identity-proxy-server" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "http-proxy" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "nokia-proxy-server" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "opendj" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "ldap" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "ping-one-http" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "http" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "oracle-unified-directory" {
		model.SmtpSecurity = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.HttpConnectTimeout = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "conjur" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.ResponseTimeout = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "amazon-aws" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpResponseTimeout = types.StringNull()
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HttpConnectTimeout = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.ResponseTimeout = types.StringNull()
		model.TrustStoreType = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
	if resourceType == "vault" {
		model.MaxConnections = types.Int64Null()
		model.SmtpSecurity = types.StringNull()
		model.VerifyCredentialsMethod = types.StringNull()
		model.AbandonOnTimeout = types.BoolNull()
		model.TrustManagerProvider = types.StringNull()
		model.MaxConnectionAge = types.StringNull()
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
		model.HostnameVerificationMethod = types.StringNull()
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AuthenticationMethod = types.StringNull()
		model.UseAdministrativeOperationControl = types.BoolNull()
		model.SmtpTimeout = types.StringNull()
		model.ValidationQueryTimeout = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.ResponseTimeout = types.StringNull()
		model.TransactionIsolationLevel = types.StringNull()
		model.MinExpiredConnectionDisconnectInterval = types.StringNull()
		model.InitialConnections = types.Int64Null()
		model.ConnectTimeout = types.StringNull()
		model.MaxResponseSize = types.StringNull()
		model.ConnectionSecurity = types.StringNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsExternalServer() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
			resourcevalidator.Conflicting(
				path.MatchRoot("password"),
				path.MatchRoot("passphrase_provider"),
			),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"amazon-aws"},
			configvalidators.Implies(
				path.MatchRoot("aws_access_key_id"),
				path.MatchRoot("aws_secret_access_key"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_host_name"),
			path.MatchRoot("type"),
			[]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "syslog", "ping-identity-proxy-server", "http-proxy", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_port"),
			path.MatchRoot("type"),
			[]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "syslog", "ping-identity-proxy-server", "http-proxy", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("smtp_security"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("user_name"),
			path.MatchRoot("type"),
			[]string{"smtp", "jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password"),
			path.MatchRoot("type"),
			[]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"smtp", "nokia-ds", "ping-identity-ds", "active-directory", "jdbc", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("smtp_timeout"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("smtp_connection_properties"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("verify_credentials_method"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_administrative_operation_control"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "ping-identity-proxy-server", "nokia-proxy-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("location"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bind_dn"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("connection_security"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("authentication_method"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory", "amazon-aws"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("health_check_connect_timeout"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_connection_age"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "syslog", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_expired_connection_disconnect_interval"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("connect_timeout"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "syslog", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "ping-one-http", "http", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_response_size"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_manager_provider"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "http", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_manager_provider"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "syslog", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "ping-one-http", "http", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("initial_connections"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_connections"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("defunct_connection_result_code"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("abandon_on_timeout"),
			path.MatchRoot("type"),
			[]string{"nokia-ds", "ping-identity-ds", "active-directory", "ping-identity-proxy-server", "nokia-proxy-server", "opendj", "ldap", "oracle-unified-directory"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jdbc_driver_type"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jdbc_driver_url"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("database_name"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("validation_query"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("validation_query_timeout"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jdbc_connection_properties"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("transaction_isolation_level"),
			path.MatchRoot("type"),
			[]string{"jdbc"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("transport_mechanism"),
			path.MatchRoot("type"),
			[]string{"syslog"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("basic_authentication_username"),
			path.MatchRoot("type"),
			[]string{"http-proxy"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("basic_authentication_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"http-proxy"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("hostname_verification_method"),
			path.MatchRoot("type"),
			[]string{"ping-one-http", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("response_timeout"),
			path.MatchRoot("type"),
			[]string{"ping-one-http", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("base_url"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ssl_cert_nickname"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_server_base_uri"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_authentication_method"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_account_name"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_connect_timeout"),
			path.MatchRoot("type"),
			[]string{"conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_response_timeout"),
			path.MatchRoot("type"),
			[]string{"conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_file"),
			path.MatchRoot("type"),
			[]string{"conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin"),
			path.MatchRoot("type"),
			[]string{"conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_type"),
			path.MatchRoot("type"),
			[]string{"conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"amazon-aws"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_access_key_id"),
			path.MatchRoot("type"),
			[]string{"amazon-aws"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_secret_access_key"),
			path.MatchRoot("type"),
			[]string{"amazon-aws"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_region_name"),
			path.MatchRoot("type"),
			[]string{"amazon-aws"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_server_base_uri"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_authentication_method"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"smtp",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"nokia-ds",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"ping-identity-ds",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"active-directory",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog",
			[]path.Expression{path.MatchRoot("server_host_name"), path.MatchRoot("transport_mechanism")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"ping-identity-proxy-server",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"http-proxy",
			[]path.Expression{path.MatchRoot("server_host_name"), path.MatchRoot("server_port")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"nokia-proxy-server",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"opendj",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"ldap",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"oracle-unified-directory",
			[]path.Expression{path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"jdbc",
			[]path.Expression{path.MatchRoot("jdbc_driver_type")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"http",
			[]path.Expression{path.MatchRoot("base_url")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"conjur",
			[]path.Expression{path.MatchRoot("conjur_server_base_uri"), path.MatchRoot("conjur_authentication_method"), path.MatchRoot("conjur_account_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"amazon-aws",
			[]path.Expression{path.MatchRoot("aws_region_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"vault",
			[]path.Expression{path.MatchRoot("vault_server_base_uri"), path.MatchRoot("vault_authentication_method")},
		),
	}
}

// Add config validators
func (r externalServerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsExternalServer()
}

// Add config validators
func (r defaultExternalServerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsExternalServer()
}

// Add optional fields to create request for smtp external-server
func addOptionalSmtpExternalServerFields(ctx context.Context, addRequest *client.AddSmtpExternalServerRequest, plan externalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SmtpSecurity) {
		smtpSecurity, err := client.NewEnumexternalServerSmtpSecurityPropFromValue(plan.SmtpSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SmtpSecurity = smtpSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UserName) {
		addRequest.UserName = plan.UserName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SmtpTimeout) {
		addRequest.SmtpTimeout = plan.SmtpTimeout.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SmtpConnectionProperties) {
		var slice []string
		plan.SmtpConnectionProperties.ElementsAs(ctx, &slice, false)
		addRequest.SmtpConnectionProperties = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for nokia-ds external-server
func addOptionalNokiaDsExternalServerFields(ctx context.Context, addRequest *client.AddNokiaDsExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	if internaltypes.IsDefined(plan.UseAdministrativeOperationControl) {
		addRequest.UseAdministrativeOperationControl = plan.UseAdministrativeOperationControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerNokiaDsConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerNokiaDsAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ping-identity-ds external-server
func addOptionalPingIdentityDsExternalServerFields(ctx context.Context, addRequest *client.AddPingIdentityDsExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	if internaltypes.IsDefined(plan.UseAdministrativeOperationControl) {
		addRequest.UseAdministrativeOperationControl = plan.UseAdministrativeOperationControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerPingIdentityDsConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerPingIdentityDsAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for active-directory external-server
func addOptionalActiveDirectoryExternalServerFields(ctx context.Context, addRequest *client.AddActiveDirectoryExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerActiveDirectoryConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerActiveDirectoryAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for jdbc external-server
func addOptionalJdbcExternalServerFields(ctx context.Context, addRequest *client.AddJdbcExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JdbcDriverURL) {
		addRequest.JdbcDriverURL = plan.JdbcDriverURL.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DatabaseName) {
		addRequest.DatabaseName = plan.DatabaseName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerHostName) {
		addRequest.ServerHostName = plan.ServerHostName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UserName) {
		addRequest.UserName = plan.UserName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidationQuery) {
		addRequest.ValidationQuery = plan.ValidationQuery.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidationQueryTimeout) {
		addRequest.ValidationQueryTimeout = plan.ValidationQueryTimeout.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.JdbcConnectionProperties) {
		var slice []string
		plan.JdbcConnectionProperties.ElementsAs(ctx, &slice, false)
		addRequest.JdbcConnectionProperties = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TransactionIsolationLevel) {
		transactionIsolationLevel, err := client.NewEnumexternalServerTransactionIsolationLevelPropFromValue(plan.TransactionIsolationLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.TransactionIsolationLevel = transactionIsolationLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for syslog external-server
func addOptionalSyslogExternalServerFields(ctx context.Context, addRequest *client.AddSyslogExternalServerRequest, plan externalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ping-identity-proxy-server external-server
func addOptionalPingIdentityProxyServerExternalServerFields(ctx context.Context, addRequest *client.AddPingIdentityProxyServerExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	if internaltypes.IsDefined(plan.UseAdministrativeOperationControl) {
		addRequest.UseAdministrativeOperationControl = plan.UseAdministrativeOperationControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerPingIdentityProxyServerConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerPingIdentityProxyServerAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for http-proxy external-server
func addOptionalHttpProxyExternalServerFields(ctx context.Context, addRequest *client.AddHttpProxyExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BasicAuthenticationUsername) {
		addRequest.BasicAuthenticationUsername = plan.BasicAuthenticationUsername.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BasicAuthenticationPassphraseProvider) {
		addRequest.BasicAuthenticationPassphraseProvider = plan.BasicAuthenticationPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for nokia-proxy-server external-server
func addOptionalNokiaProxyServerExternalServerFields(ctx context.Context, addRequest *client.AddNokiaProxyServerExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	if internaltypes.IsDefined(plan.UseAdministrativeOperationControl) {
		addRequest.UseAdministrativeOperationControl = plan.UseAdministrativeOperationControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerNokiaProxyServerConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerNokiaProxyServerAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for opendj external-server
func addOptionalOpendjExternalServerFields(ctx context.Context, addRequest *client.AddOpendjExternalServerRequest, plan externalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerOpendjConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerOpendjAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ldap external-server
func addOptionalLdapExternalServerFields(ctx context.Context, addRequest *client.AddLdapExternalServerRequest, plan externalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerLdapConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerLdapAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ping-one-http external-server
func addOptionalPingOneHttpExternalServerFields(ctx context.Context, addRequest *client.AddPingOneHttpExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HostnameVerificationMethod) {
		hostnameVerificationMethod, err := client.NewEnumexternalServerPingOneHttpHostnameVerificationMethodPropFromValue(plan.HostnameVerificationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.HostnameVerificationMethod = hostnameVerificationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResponseTimeout) {
		addRequest.ResponseTimeout = plan.ResponseTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for http external-server
func addOptionalHttpExternalServerFields(ctx context.Context, addRequest *client.AddHttpExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HostnameVerificationMethod) {
		hostnameVerificationMethod, err := client.NewEnumexternalServerHttpHostnameVerificationMethodPropFromValue(plan.HostnameVerificationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.HostnameVerificationMethod = hostnameVerificationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		addRequest.SslCertNickname = plan.SslCertNickname.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResponseTimeout) {
		addRequest.ResponseTimeout = plan.ResponseTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for oracle-unified-directory external-server
func addOptionalOracleUnifiedDirectoryExternalServerFields(ctx context.Context, addRequest *client.AddOracleUnifiedDirectoryExternalServerRequest, plan externalServerResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerOracleUnifiedDirectoryConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerOracleUnifiedDirectoryAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for conjur external-server
func addOptionalConjurExternalServerFields(ctx context.Context, addRequest *client.AddConjurExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpConnectTimeout) {
		addRequest.HttpConnectTimeout = plan.HttpConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpResponseTimeout) {
		addRequest.HttpResponseTimeout = plan.HttpResponseTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreFile) {
		addRequest.TrustStoreFile = plan.TrustStoreFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		addRequest.TrustStorePin = plan.TrustStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		addRequest.TrustStoreType = plan.TrustStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for amazon-aws external-server
func addOptionalAmazonAwsExternalServerFields(ctx context.Context, addRequest *client.AddAmazonAwsExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerAmazonAwsAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsAccessKeyID) {
		addRequest.AwsAccessKeyID = plan.AwsAccessKeyID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsSecretAccessKey) {
		addRequest.AwsSecretAccessKey = plan.AwsSecretAccessKey.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for vault external-server
func addOptionalVaultExternalServerFields(ctx context.Context, addRequest *client.AddVaultExternalServerRequest, plan externalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpConnectTimeout) {
		addRequest.HttpConnectTimeout = plan.HttpConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpResponseTimeout) {
		addRequest.HttpResponseTimeout = plan.HttpResponseTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreFile) {
		addRequest.TrustStoreFile = plan.TrustStoreFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		addRequest.TrustStorePin = plan.TrustStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		addRequest.TrustStoreType = plan.TrustStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateExternalServerUnknownValues(model *externalServerResourceModel) {
	if model.VaultServerBaseURI.IsUnknown() || model.VaultServerBaseURI.IsNull() {
		model.VaultServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ConjurServerBaseURI.IsUnknown() || model.ConjurServerBaseURI.IsNull() {
		model.ConjurServerBaseURI, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefunctConnectionResultCode.IsUnknown() || model.DefunctConnectionResultCode.IsNull() {
		model.DefunctConnectionResultCode, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SmtpConnectionProperties.IsUnknown() || model.SmtpConnectionProperties.IsNull() {
		model.SmtpConnectionProperties, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.JdbcConnectionProperties.IsUnknown() || model.JdbcConnectionProperties.IsNull() {
		model.JdbcConnectionProperties, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *externalServerResourceModel) populateAllComputedStringAttributes() {
	if model.MaxResponseSize.IsUnknown() || model.MaxResponseSize.IsNull() {
		model.MaxResponseSize = types.StringValue("")
	}
	if model.TrustManagerProvider.IsUnknown() || model.TrustManagerProvider.IsNull() {
		model.TrustManagerProvider = types.StringValue("")
	}
	if model.HostnameVerificationMethod.IsUnknown() || model.HostnameVerificationMethod.IsNull() {
		model.HostnameVerificationMethod = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.BasicAuthenticationPassphraseProvider.IsUnknown() || model.BasicAuthenticationPassphraseProvider.IsNull() {
		model.BasicAuthenticationPassphraseProvider = types.StringValue("")
	}
	if model.ValidationQueryTimeout.IsUnknown() || model.ValidationQueryTimeout.IsNull() {
		model.ValidationQueryTimeout = types.StringValue("")
	}
	if model.JdbcDriverType.IsUnknown() || model.JdbcDriverType.IsNull() {
		model.JdbcDriverType = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.SmtpSecurity.IsUnknown() || model.SmtpSecurity.IsNull() {
		model.SmtpSecurity = types.StringValue("")
	}
	if model.AuthenticationMethod.IsUnknown() || model.AuthenticationMethod.IsNull() {
		model.AuthenticationMethod = types.StringValue("")
	}
	if model.VerifyCredentialsMethod.IsUnknown() || model.VerifyCredentialsMethod.IsNull() {
		model.VerifyCredentialsMethod = types.StringValue("")
	}
	if model.TransactionIsolationLevel.IsUnknown() || model.TransactionIsolationLevel.IsNull() {
		model.TransactionIsolationLevel = types.StringValue("")
	}
	if model.ResponseTimeout.IsUnknown() || model.ResponseTimeout.IsNull() {
		model.ResponseTimeout = types.StringValue("")
	}
	if model.SmtpTimeout.IsUnknown() || model.SmtpTimeout.IsNull() {
		model.SmtpTimeout = types.StringValue("")
	}
	if model.BasicAuthenticationUsername.IsUnknown() || model.BasicAuthenticationUsername.IsNull() {
		model.BasicAuthenticationUsername = types.StringValue("")
	}
	if model.ServerHostName.IsUnknown() || model.ServerHostName.IsNull() {
		model.ServerHostName = types.StringValue("")
	}
	if model.MinExpiredConnectionDisconnectInterval.IsUnknown() || model.MinExpiredConnectionDisconnectInterval.IsNull() {
		model.MinExpiredConnectionDisconnectInterval = types.StringValue("")
	}
	if model.BaseURL.IsUnknown() || model.BaseURL.IsNull() {
		model.BaseURL = types.StringValue("")
	}
	if model.MaxConnectionAge.IsUnknown() || model.MaxConnectionAge.IsNull() {
		model.MaxConnectionAge = types.StringValue("")
	}
	if model.AwsAccessKeyID.IsUnknown() || model.AwsAccessKeyID.IsNull() {
		model.AwsAccessKeyID = types.StringValue("")
	}
	if model.Password.IsUnknown() || model.Password.IsNull() {
		model.Password = types.StringValue("")
	}
	if model.TrustStorePin.IsUnknown() || model.TrustStorePin.IsNull() {
		model.TrustStorePin = types.StringValue("")
	}
	if model.ConnectTimeout.IsUnknown() || model.ConnectTimeout.IsNull() {
		model.ConnectTimeout = types.StringValue("")
	}
	if model.UserName.IsUnknown() || model.UserName.IsNull() {
		model.UserName = types.StringValue("")
	}
	if model.AwsRegionName.IsUnknown() || model.AwsRegionName.IsNull() {
		model.AwsRegionName = types.StringValue("")
	}
	if model.VaultAuthenticationMethod.IsUnknown() || model.VaultAuthenticationMethod.IsNull() {
		model.VaultAuthenticationMethod = types.StringValue("")
	}
	if model.ConjurAccountName.IsUnknown() || model.ConjurAccountName.IsNull() {
		model.ConjurAccountName = types.StringValue("")
	}
	if model.HealthCheckConnectTimeout.IsUnknown() || model.HealthCheckConnectTimeout.IsNull() {
		model.HealthCheckConnectTimeout = types.StringValue("")
	}
	if model.DatabaseName.IsUnknown() || model.DatabaseName.IsNull() {
		model.DatabaseName = types.StringValue("")
	}
	if model.ConjurAuthenticationMethod.IsUnknown() || model.ConjurAuthenticationMethod.IsNull() {
		model.ConjurAuthenticationMethod = types.StringValue("")
	}
	if model.ConnectionSecurity.IsUnknown() || model.ConnectionSecurity.IsNull() {
		model.ConnectionSecurity = types.StringValue("")
	}
	if model.ValidationQuery.IsUnknown() || model.ValidationQuery.IsNull() {
		model.ValidationQuery = types.StringValue("")
	}
	if model.TransportMechanism.IsUnknown() || model.TransportMechanism.IsNull() {
		model.TransportMechanism = types.StringValue("")
	}
	if model.KeyManagerProvider.IsUnknown() || model.KeyManagerProvider.IsNull() {
		model.KeyManagerProvider = types.StringValue("")
	}
	if model.SslCertNickname.IsUnknown() || model.SslCertNickname.IsNull() {
		model.SslCertNickname = types.StringValue("")
	}
	if model.HttpConnectTimeout.IsUnknown() || model.HttpConnectTimeout.IsNull() {
		model.HttpConnectTimeout = types.StringValue("")
	}
	if model.JdbcDriverURL.IsUnknown() || model.JdbcDriverURL.IsNull() {
		model.JdbcDriverURL = types.StringValue("")
	}
	if model.PassphraseProvider.IsUnknown() || model.PassphraseProvider.IsNull() {
		model.PassphraseProvider = types.StringValue("")
	}
	if model.HttpResponseTimeout.IsUnknown() || model.HttpResponseTimeout.IsNull() {
		model.HttpResponseTimeout = types.StringValue("")
	}
	if model.AwsSecretAccessKey.IsUnknown() || model.AwsSecretAccessKey.IsNull() {
		model.AwsSecretAccessKey = types.StringValue("")
	}
	if model.BindDN.IsUnknown() || model.BindDN.IsNull() {
		model.BindDN = types.StringValue("")
	}
	if model.TrustStoreFile.IsUnknown() || model.TrustStoreFile.IsNull() {
		model.TrustStoreFile = types.StringValue("")
	}
	if model.TrustStoreType.IsUnknown() || model.TrustStoreType.IsNull() {
		model.TrustStoreType = types.StringValue("")
	}
	if model.Location.IsUnknown() || model.Location.IsNull() {
		model.Location = types.StringValue("")
	}
}

// Read a SmtpExternalServerResponse object into the model struct
func readSmtpExternalServerResponse(ctx context.Context, r *client.SmtpExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.SmtpSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerSmtpSecurityProp(r.SmtpSecurity), true)
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, internaltypes.IsEmptyString(expectedValues.UserName))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.SmtpTimeout = internaltypes.StringTypeOrNil(r.SmtpTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("smtp_timeout",
		expectedValues.SmtpTimeout, state.SmtpTimeout, diagnostics)
	state.SmtpConnectionProperties = internaltypes.GetStringSet(r.SmtpConnectionProperties)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a NokiaDsExternalServerResponse object into the model struct
func readNokiaDsExternalServerResponse(ctx context.Context, r *client.NokiaDsExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nokia-ds")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a PingIdentityDsExternalServerResponse object into the model struct
func readPingIdentityDsExternalServerResponse(ctx context.Context, r *client.PingIdentityDsExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-identity-ds")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a ActiveDirectoryExternalServerResponse object into the model struct
func readActiveDirectoryExternalServerResponse(ctx context.Context, r *client.ActiveDirectoryExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("active-directory")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a JdbcExternalServerResponse object into the model struct
func readJdbcExternalServerResponse(ctx context.Context, r *client.JdbcExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JdbcDriverType = types.StringValue(r.JdbcDriverType.String())
	state.JdbcDriverURL = internaltypes.StringTypeOrNil(r.JdbcDriverURL, internaltypes.IsEmptyString(expectedValues.JdbcDriverURL))
	state.DatabaseName = internaltypes.StringTypeOrNil(r.DatabaseName, internaltypes.IsEmptyString(expectedValues.DatabaseName))
	state.ServerHostName = internaltypes.StringTypeOrNil(r.ServerHostName, internaltypes.IsEmptyString(expectedValues.ServerHostName))
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, internaltypes.IsEmptyString(expectedValues.UserName))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ValidationQuery = internaltypes.StringTypeOrNil(r.ValidationQuery, internaltypes.IsEmptyString(expectedValues.ValidationQuery))
	state.ValidationQueryTimeout = internaltypes.StringTypeOrNil(r.ValidationQueryTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("validation_query_timeout",
		expectedValues.ValidationQueryTimeout, state.ValidationQueryTimeout, diagnostics)
	state.JdbcConnectionProperties = internaltypes.GetStringSet(r.JdbcConnectionProperties)
	state.TransactionIsolationLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerTransactionIsolationLevelProp(r.TransactionIsolationLevel), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a SyslogExternalServerResponse object into the model struct
func readSyslogExternalServerResponse(ctx context.Context, r *client.SyslogExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.TransportMechanism = types.StringValue(r.TransportMechanism.String())
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.TrustManagerProvider = types.StringValue(r.TrustManagerProvider)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a PingIdentityProxyServerExternalServerResponse object into the model struct
func readPingIdentityProxyServerExternalServerResponse(ctx context.Context, r *client.PingIdentityProxyServerExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-identity-proxy-server")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a HttpProxyExternalServerResponse object into the model struct
func readHttpProxyExternalServerResponse(ctx context.Context, r *client.HttpProxyExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http-proxy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.BasicAuthenticationUsername = internaltypes.StringTypeOrNil(r.BasicAuthenticationUsername, internaltypes.IsEmptyString(expectedValues.BasicAuthenticationUsername))
	state.BasicAuthenticationPassphraseProvider = internaltypes.StringTypeOrNil(r.BasicAuthenticationPassphraseProvider, internaltypes.IsEmptyString(expectedValues.BasicAuthenticationPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a NokiaProxyServerExternalServerResponse object into the model struct
func readNokiaProxyServerExternalServerResponse(ctx context.Context, r *client.NokiaProxyServerExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nokia-proxy-server")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a OpendjExternalServerResponse object into the model struct
func readOpendjExternalServerResponse(ctx context.Context, r *client.OpendjExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("opendj")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a LdapExternalServerResponse object into the model struct
func readLdapExternalServerResponse(ctx context.Context, r *client.LdapExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a PingOneHttpExternalServerResponse object into the model struct
func readPingOneHttpExternalServerResponse(ctx context.Context, r *client.PingOneHttpExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one-http")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerPingOneHttpHostnameVerificationMethodProp(r.HostnameVerificationMethod), true)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.ConnectTimeout = internaltypes.StringTypeOrNil(r.ConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.ResponseTimeout = internaltypes.StringTypeOrNil(r.ResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("response_timeout",
		expectedValues.ResponseTimeout, state.ResponseTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a HttpExternalServerResponse object into the model struct
func readHttpExternalServerResponse(ctx context.Context, r *client.HttpExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseURL = types.StringValue(r.BaseURL)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerHttpHostnameVerificationMethodProp(r.HostnameVerificationMethod), true)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.ConnectTimeout = internaltypes.StringTypeOrNil(r.ConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.ResponseTimeout = internaltypes.StringTypeOrNil(r.ResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("response_timeout",
		expectedValues.ResponseTimeout, state.ResponseTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a OracleUnifiedDirectoryExternalServerResponse object into the model struct
func readOracleUnifiedDirectoryExternalServerResponse(ctx context.Context, r *client.OracleUnifiedDirectoryExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oracle-unified-directory")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, true)
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a ConjurExternalServerResponse object into the model struct
func readConjurExternalServerResponse(ctx context.Context, r *client.ConjurExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurServerBaseURI = internaltypes.GetStringSet(r.ConjurServerBaseURI)
	state.ConjurAuthenticationMethod = types.StringValue(r.ConjurAuthenticationMethod)
	state.ConjurAccountName = types.StringValue(r.ConjurAccountName)
	state.HttpConnectTimeout = internaltypes.StringTypeOrNil(r.HttpConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_connect_timeout",
		expectedValues.HttpConnectTimeout, state.HttpConnectTimeout, diagnostics)
	state.HttpResponseTimeout = internaltypes.StringTypeOrNil(r.HttpResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_response_timeout",
		expectedValues.HttpResponseTimeout, state.HttpResponseTimeout, diagnostics)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a AmazonAwsExternalServerResponse object into the model struct
func readAmazonAwsExternalServerResponse(ctx context.Context, r *client.AmazonAwsExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-aws")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.AuthenticationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerAmazonAwsAuthenticationMethodProp(r.AuthenticationMethod), true)
	state.AwsAccessKeyID = internaltypes.StringTypeOrNil(r.AwsAccessKeyID, internaltypes.IsEmptyString(expectedValues.AwsAccessKeyID))
	state.AwsRegionName = types.StringValue(r.AwsRegionName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Read a VaultExternalServerResponse object into the model struct
func readVaultExternalServerResponse(ctx context.Context, r *client.VaultExternalServerResponse, state *externalServerResourceModel, expectedValues *externalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultServerBaseURI = internaltypes.GetStringSet(r.VaultServerBaseURI)
	state.VaultAuthenticationMethod = types.StringValue(r.VaultAuthenticationMethod)
	state.HttpConnectTimeout = internaltypes.StringTypeOrNil(r.HttpConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_connect_timeout",
		expectedValues.HttpConnectTimeout, state.HttpConnectTimeout, diagnostics)
	state.HttpResponseTimeout = internaltypes.StringTypeOrNil(r.HttpResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_response_timeout",
		expectedValues.HttpResponseTimeout, state.HttpResponseTimeout, diagnostics)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExternalServerUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *externalServerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *externalServerResourceModel) {
	if !expectedValues.Password.IsUnknown() {
		state.Password = expectedValues.Password
	}
	if !expectedValues.TrustStorePin.IsUnknown() {
		state.TrustStorePin = expectedValues.TrustStorePin
	}
	if !expectedValues.AwsSecretAccessKey.IsUnknown() {
		state.AwsSecretAccessKey = expectedValues.AwsSecretAccessKey
	}
}

// Create any update operations necessary to make the state match the plan
func createExternalServerOperations(plan externalServerResourceModel, state externalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.VaultServerBaseURI, state.VaultServerBaseURI, "vault-server-base-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultAuthenticationMethod, state.VaultAuthenticationMethod, "vault-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ConjurServerBaseURI, state.ConjurServerBaseURI, "conjur-server-base-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsAccessKeyID, state.AwsAccessKeyID, "aws-access-key-id")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsSecretAccessKey, state.AwsSecretAccessKey, "aws-secret-access-key")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsRegionName, state.AwsRegionName, "aws-region-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurAuthenticationMethod, state.ConjurAuthenticationMethod, "conjur-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurAccountName, state.ConjurAccountName, "conjur-account-name")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpConnectTimeout, state.HttpConnectTimeout, "http-connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpResponseTimeout, state.HttpResponseTimeout, "http-response-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseURL, state.BaseURL, "base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.HostnameVerificationMethod, state.HostnameVerificationMethod, "hostname-verification-method")
	operations.AddStringOperationIfNecessary(&ops, plan.JdbcDriverType, state.JdbcDriverType, "jdbc-driver-type")
	operations.AddStringOperationIfNecessary(&ops, plan.JdbcDriverURL, state.JdbcDriverURL, "jdbc-driver-url")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.ResponseTimeout, state.ResponseTimeout, "response-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.BasicAuthenticationUsername, state.BasicAuthenticationUsername, "basic-authentication-username")
	operations.AddStringOperationIfNecessary(&ops, plan.BasicAuthenticationPassphraseProvider, state.BasicAuthenticationPassphraseProvider, "basic-authentication-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TransportMechanism, state.TransportMechanism, "transport-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.DatabaseName, state.DatabaseName, "database-name")
	operations.AddStringOperationIfNecessary(&ops, plan.VerifyCredentialsMethod, state.VerifyCredentialsMethod, "verify-credentials-method")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseAdministrativeOperationControl, state.UseAdministrativeOperationControl, "use-administrative-operation-control")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.Location, state.Location, "location")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidationQuery, state.ValidationQuery, "validation-query")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidationQueryTimeout, state.ValidationQueryTimeout, "validation-query-timeout")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JdbcConnectionProperties, state.JdbcConnectionProperties, "jdbc-connection-properties")
	operations.AddStringOperationIfNecessary(&ops, plan.TransactionIsolationLevel, state.TransactionIsolationLevel, "transaction-isolation-level")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDN, state.BindDN, "bind-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpSecurity, state.SmtpSecurity, "smtp-security")
	operations.AddStringOperationIfNecessary(&ops, plan.UserName, state.UserName, "user-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionSecurity, state.ConnectionSecurity, "connection-security")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthenticationMethod, state.AuthenticationMethod, "authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, "health-check-connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxConnectionAge, state.MaxConnectionAge, "max-connection-age")
	operations.AddStringOperationIfNecessary(&ops, plan.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, "min-expired-connection-disconnect-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectTimeout, state.ConnectTimeout, "connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxResponseSize, state.MaxResponseSize, "max-response-size")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefunctConnectionResultCode, state.DefunctConnectionResultCode, "defunct-connection-result-code")
	operations.AddBoolOperationIfNecessary(&ops, plan.AbandonOnTimeout, state.AbandonOnTimeout, "abandon-on-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.PassphraseProvider, state.PassphraseProvider, "passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.SmtpTimeout, state.SmtpTimeout, "smtp-timeout")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SmtpConnectionProperties, state.SmtpConnectionProperties, "smtp-connection-properties")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a smtp external-server
func (r *externalServerResource) CreateSmtpExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddSmtpExternalServerRequest([]client.EnumsmtpExternalServerSchemaUrn{client.ENUMSMTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERSMTP},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalSmtpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddSmtpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readSmtpExternalServerResponse(ctx, addResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a nokia-ds external-server
func (r *externalServerResource) CreateNokiaDsExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddNokiaDsExternalServerRequest([]client.EnumnokiaDsExternalServerSchemaUrn{client.ENUMNOKIADSEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERNOKIA_DS},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalNokiaDsExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddNokiaDsExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readNokiaDsExternalServerResponse(ctx, addResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ping-identity-ds external-server
func (r *externalServerResource) CreatePingIdentityDsExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddPingIdentityDsExternalServerRequest([]client.EnumpingIdentityDsExternalServerSchemaUrn{client.ENUMPINGIDENTITYDSEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERPING_IDENTITY_DS},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalPingIdentityDsExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddPingIdentityDsExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readPingIdentityDsExternalServerResponse(ctx, addResponse.PingIdentityDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a active-directory external-server
func (r *externalServerResource) CreateActiveDirectoryExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddActiveDirectoryExternalServerRequest([]client.EnumactiveDirectoryExternalServerSchemaUrn{client.ENUMACTIVEDIRECTORYEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERACTIVE_DIRECTORY},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalActiveDirectoryExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddActiveDirectoryExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readActiveDirectoryExternalServerResponse(ctx, addResponse.ActiveDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jdbc external-server
func (r *externalServerResource) CreateJdbcExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	jdbcDriverType, err := client.NewEnumexternalServerJdbcDriverTypePropFromValue(plan.JdbcDriverType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for JdbcDriverType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddJdbcExternalServerRequest([]client.EnumjdbcExternalServerSchemaUrn{client.ENUMJDBCEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERJDBC},
		*jdbcDriverType,
		plan.Name.ValueString())
	err = addOptionalJdbcExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddJdbcExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readJdbcExternalServerResponse(ctx, addResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog external-server
func (r *externalServerResource) CreateSyslogExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	transportMechanism, err := client.NewEnumexternalServerTransportMechanismPropFromValue(plan.TransportMechanism.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for TransportMechanism", err.Error())
		return nil, err
	}
	addRequest := client.NewAddSyslogExternalServerRequest([]client.EnumsyslogExternalServerSchemaUrn{client.ENUMSYSLOGEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERSYSLOG},
		plan.ServerHostName.ValueString(),
		*transportMechanism,
		plan.Name.ValueString())
	err = addOptionalSyslogExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddSyslogExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readSyslogExternalServerResponse(ctx, addResponse.SyslogExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ping-identity-proxy-server external-server
func (r *externalServerResource) CreatePingIdentityProxyServerExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddPingIdentityProxyServerExternalServerRequest([]client.EnumpingIdentityProxyServerExternalServerSchemaUrn{client.ENUMPINGIDENTITYPROXYSERVEREXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERPING_IDENTITY_PROXY_SERVER},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalPingIdentityProxyServerExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddPingIdentityProxyServerExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readPingIdentityProxyServerExternalServerResponse(ctx, addResponse.PingIdentityProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a http-proxy external-server
func (r *externalServerResource) CreateHttpProxyExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddHttpProxyExternalServerRequest([]client.EnumhttpProxyExternalServerSchemaUrn{client.ENUMHTTPPROXYEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERHTTP_PROXY},
		plan.ServerHostName.ValueString(),
		plan.ServerPort.ValueInt64(),
		plan.Name.ValueString())
	err := addOptionalHttpProxyExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddHttpProxyExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readHttpProxyExternalServerResponse(ctx, addResponse.HttpProxyExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a nokia-proxy-server external-server
func (r *externalServerResource) CreateNokiaProxyServerExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddNokiaProxyServerExternalServerRequest([]client.EnumnokiaProxyServerExternalServerSchemaUrn{client.ENUMNOKIAPROXYSERVEREXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERNOKIA_PROXY_SERVER},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalNokiaProxyServerExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddNokiaProxyServerExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readNokiaProxyServerExternalServerResponse(ctx, addResponse.NokiaProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a opendj external-server
func (r *externalServerResource) CreateOpendjExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddOpendjExternalServerRequest([]client.EnumopendjExternalServerSchemaUrn{client.ENUMOPENDJEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVEROPENDJ},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalOpendjExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddOpendjExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readOpendjExternalServerResponse(ctx, addResponse.OpendjExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldap external-server
func (r *externalServerResource) CreateLdapExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddLdapExternalServerRequest([]client.EnumldapExternalServerSchemaUrn{client.ENUMLDAPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERLDAP},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalLdapExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddLdapExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readLdapExternalServerResponse(ctx, addResponse.LdapExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ping-one-http external-server
func (r *externalServerResource) CreatePingOneHttpExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddPingOneHttpExternalServerRequest([]client.EnumpingOneHttpExternalServerSchemaUrn{client.ENUMPINGONEHTTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERPING_ONE_HTTP},
		plan.Name.ValueString())
	err := addOptionalPingOneHttpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddPingOneHttpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readPingOneHttpExternalServerResponse(ctx, addResponse.PingOneHttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a http external-server
func (r *externalServerResource) CreateHttpExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddHttpExternalServerRequest([]client.EnumhttpExternalServerSchemaUrn{client.ENUMHTTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERHTTP},
		plan.BaseURL.ValueString(),
		plan.Name.ValueString())
	err := addOptionalHttpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddHttpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readHttpExternalServerResponse(ctx, addResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a oracle-unified-directory external-server
func (r *externalServerResource) CreateOracleUnifiedDirectoryExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddOracleUnifiedDirectoryExternalServerRequest([]client.EnumoracleUnifiedDirectoryExternalServerSchemaUrn{client.ENUMORACLEUNIFIEDDIRECTORYEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERORACLE_UNIFIED_DIRECTORY},
		plan.ServerHostName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalOracleUnifiedDirectoryExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddOracleUnifiedDirectoryExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readOracleUnifiedDirectoryExternalServerResponse(ctx, addResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a conjur external-server
func (r *externalServerResource) CreateConjurExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	var ConjurServerBaseURISlice []string
	plan.ConjurServerBaseURI.ElementsAs(ctx, &ConjurServerBaseURISlice, false)
	addRequest := client.NewAddConjurExternalServerRequest([]client.EnumconjurExternalServerSchemaUrn{client.ENUMCONJUREXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERCONJUR},
		ConjurServerBaseURISlice,
		plan.ConjurAuthenticationMethod.ValueString(),
		plan.ConjurAccountName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalConjurExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddConjurExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readConjurExternalServerResponse(ctx, addResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a amazon-aws external-server
func (r *externalServerResource) CreateAmazonAwsExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	addRequest := client.NewAddAmazonAwsExternalServerRequest([]client.EnumamazonAwsExternalServerSchemaUrn{client.ENUMAMAZONAWSEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERAMAZON_AWS},
		plan.AwsRegionName.ValueString(),
		plan.Name.ValueString())
	err := addOptionalAmazonAwsExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddAmazonAwsExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readAmazonAwsExternalServerResponse(ctx, addResponse.AmazonAwsExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a vault external-server
func (r *externalServerResource) CreateVaultExternalServer(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan externalServerResourceModel) (*externalServerResourceModel, error) {
	var VaultServerBaseURISlice []string
	plan.VaultServerBaseURI.ElementsAs(ctx, &VaultServerBaseURISlice, false)
	addRequest := client.NewAddVaultExternalServerRequest([]client.EnumvaultExternalServerSchemaUrn{client.ENUMVAULTEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERVAULT},
		VaultServerBaseURISlice,
		plan.VaultAuthenticationMethod.ValueString(),
		plan.Name.ValueString())
	err := addOptionalVaultExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for External Server", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerAPI.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddVaultExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerAPI.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the External Server", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state externalServerResourceModel
	readVaultExternalServerResponse(ctx, addResponse.VaultExternalServerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *externalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan externalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *externalServerResourceModel
	var err error
	if plan.Type.ValueString() == "smtp" {
		state, err = r.CreateSmtpExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "nokia-ds" {
		state, err = r.CreateNokiaDsExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ping-identity-ds" {
		state, err = r.CreatePingIdentityDsExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "active-directory" {
		state, err = r.CreateActiveDirectoryExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jdbc" {
		state, err = r.CreateJdbcExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog" {
		state, err = r.CreateSyslogExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ping-identity-proxy-server" {
		state, err = r.CreatePingIdentityProxyServerExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "http-proxy" {
		state, err = r.CreateHttpProxyExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "nokia-proxy-server" {
		state, err = r.CreateNokiaProxyServerExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "opendj" {
		state, err = r.CreateOpendjExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldap" {
		state, err = r.CreateLdapExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ping-one-http" {
		state, err = r.CreatePingOneHttpExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "http" {
		state, err = r.CreateHttpExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "oracle-unified-directory" {
		state, err = r.CreateOracleUnifiedDirectoryExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "conjur" {
		state, err = r.CreateConjurExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "amazon-aws" {
		state, err = r.CreateAmazonAwsExternalServer(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "vault" {
		state, err = r.CreateVaultExternalServer(ctx, req, resp, plan)
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
func (r *defaultExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan externalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerAPI.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state externalServerResourceModel
	if readResponse.SmtpExternalServerResponse != nil {
		readSmtpExternalServerResponse(ctx, readResponse.SmtpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaDsExternalServerResponse != nil {
		readNokiaDsExternalServerResponse(ctx, readResponse.NokiaDsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityDsExternalServerResponse != nil {
		readPingIdentityDsExternalServerResponse(ctx, readResponse.PingIdentityDsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveDirectoryExternalServerResponse != nil {
		readActiveDirectoryExternalServerResponse(ctx, readResponse.ActiveDirectoryExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcExternalServerResponse != nil {
		readJdbcExternalServerResponse(ctx, readResponse.JdbcExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogExternalServerResponse != nil {
		readSyslogExternalServerResponse(ctx, readResponse.SyslogExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityProxyServerExternalServerResponse != nil {
		readPingIdentityProxyServerExternalServerResponse(ctx, readResponse.PingIdentityProxyServerExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpProxyExternalServerResponse != nil {
		readHttpProxyExternalServerResponse(ctx, readResponse.HttpProxyExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaProxyServerExternalServerResponse != nil {
		readNokiaProxyServerExternalServerResponse(ctx, readResponse.NokiaProxyServerExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OpendjExternalServerResponse != nil {
		readOpendjExternalServerResponse(ctx, readResponse.OpendjExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapExternalServerResponse != nil {
		readLdapExternalServerResponse(ctx, readResponse.LdapExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingOneHttpExternalServerResponse != nil {
		readPingOneHttpExternalServerResponse(ctx, readResponse.PingOneHttpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpExternalServerResponse != nil {
		readHttpExternalServerResponse(ctx, readResponse.HttpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OracleUnifiedDirectoryExternalServerResponse != nil {
		readOracleUnifiedDirectoryExternalServerResponse(ctx, readResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurExternalServerResponse != nil {
		readConjurExternalServerResponse(ctx, readResponse.ConjurExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonAwsExternalServerResponse != nil {
		readAmazonAwsExternalServerResponse(ctx, readResponse.AmazonAwsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultExternalServerResponse != nil {
		readVaultExternalServerResponse(ctx, readResponse.VaultExternalServerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerAPI.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerAPI.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SmtpExternalServerResponse != nil {
			readSmtpExternalServerResponse(ctx, updateResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NokiaDsExternalServerResponse != nil {
			readNokiaDsExternalServerResponse(ctx, updateResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingIdentityDsExternalServerResponse != nil {
			readPingIdentityDsExternalServerResponse(ctx, updateResponse.PingIdentityDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ActiveDirectoryExternalServerResponse != nil {
			readActiveDirectoryExternalServerResponse(ctx, updateResponse.ActiveDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcExternalServerResponse != nil {
			readJdbcExternalServerResponse(ctx, updateResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogExternalServerResponse != nil {
			readSyslogExternalServerResponse(ctx, updateResponse.SyslogExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingIdentityProxyServerExternalServerResponse != nil {
			readPingIdentityProxyServerExternalServerResponse(ctx, updateResponse.PingIdentityProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpProxyExternalServerResponse != nil {
			readHttpProxyExternalServerResponse(ctx, updateResponse.HttpProxyExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NokiaProxyServerExternalServerResponse != nil {
			readNokiaProxyServerExternalServerResponse(ctx, updateResponse.NokiaProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OpendjExternalServerResponse != nil {
			readOpendjExternalServerResponse(ctx, updateResponse.OpendjExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapExternalServerResponse != nil {
			readLdapExternalServerResponse(ctx, updateResponse.LdapExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingOneHttpExternalServerResponse != nil {
			readPingOneHttpExternalServerResponse(ctx, updateResponse.PingOneHttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpExternalServerResponse != nil {
			readHttpExternalServerResponse(ctx, updateResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OracleUnifiedDirectoryExternalServerResponse != nil {
			readOracleUnifiedDirectoryExternalServerResponse(ctx, updateResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurExternalServerResponse != nil {
			readConjurExternalServerResponse(ctx, updateResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonAwsExternalServerResponse != nil {
			readAmazonAwsExternalServerResponse(ctx, updateResponse.AmazonAwsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultExternalServerResponse != nil {
			readVaultExternalServerResponse(ctx, updateResponse.VaultExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *externalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readExternalServer(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readExternalServer(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state externalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerAPI.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the External Server", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the External Server", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SmtpExternalServerResponse != nil {
		readSmtpExternalServerResponse(ctx, readResponse.SmtpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaDsExternalServerResponse != nil {
		readNokiaDsExternalServerResponse(ctx, readResponse.NokiaDsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityDsExternalServerResponse != nil {
		readPingIdentityDsExternalServerResponse(ctx, readResponse.PingIdentityDsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveDirectoryExternalServerResponse != nil {
		readActiveDirectoryExternalServerResponse(ctx, readResponse.ActiveDirectoryExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcExternalServerResponse != nil {
		readJdbcExternalServerResponse(ctx, readResponse.JdbcExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogExternalServerResponse != nil {
		readSyslogExternalServerResponse(ctx, readResponse.SyslogExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityProxyServerExternalServerResponse != nil {
		readPingIdentityProxyServerExternalServerResponse(ctx, readResponse.PingIdentityProxyServerExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpProxyExternalServerResponse != nil {
		readHttpProxyExternalServerResponse(ctx, readResponse.HttpProxyExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaProxyServerExternalServerResponse != nil {
		readNokiaProxyServerExternalServerResponse(ctx, readResponse.NokiaProxyServerExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OpendjExternalServerResponse != nil {
		readOpendjExternalServerResponse(ctx, readResponse.OpendjExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapExternalServerResponse != nil {
		readLdapExternalServerResponse(ctx, readResponse.LdapExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingOneHttpExternalServerResponse != nil {
		readPingOneHttpExternalServerResponse(ctx, readResponse.PingOneHttpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpExternalServerResponse != nil {
		readHttpExternalServerResponse(ctx, readResponse.HttpExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OracleUnifiedDirectoryExternalServerResponse != nil {
		readOracleUnifiedDirectoryExternalServerResponse(ctx, readResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurExternalServerResponse != nil {
		readConjurExternalServerResponse(ctx, readResponse.ConjurExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonAwsExternalServerResponse != nil {
		readAmazonAwsExternalServerResponse(ctx, readResponse.AmazonAwsExternalServerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultExternalServerResponse != nil {
		readVaultExternalServerResponse(ctx, readResponse.VaultExternalServerResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *externalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan externalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state externalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerAPI.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerAPI.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SmtpExternalServerResponse != nil {
			readSmtpExternalServerResponse(ctx, updateResponse.SmtpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NokiaDsExternalServerResponse != nil {
			readNokiaDsExternalServerResponse(ctx, updateResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingIdentityDsExternalServerResponse != nil {
			readPingIdentityDsExternalServerResponse(ctx, updateResponse.PingIdentityDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ActiveDirectoryExternalServerResponse != nil {
			readActiveDirectoryExternalServerResponse(ctx, updateResponse.ActiveDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcExternalServerResponse != nil {
			readJdbcExternalServerResponse(ctx, updateResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogExternalServerResponse != nil {
			readSyslogExternalServerResponse(ctx, updateResponse.SyslogExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingIdentityProxyServerExternalServerResponse != nil {
			readPingIdentityProxyServerExternalServerResponse(ctx, updateResponse.PingIdentityProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpProxyExternalServerResponse != nil {
			readHttpProxyExternalServerResponse(ctx, updateResponse.HttpProxyExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NokiaProxyServerExternalServerResponse != nil {
			readNokiaProxyServerExternalServerResponse(ctx, updateResponse.NokiaProxyServerExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OpendjExternalServerResponse != nil {
			readOpendjExternalServerResponse(ctx, updateResponse.OpendjExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapExternalServerResponse != nil {
			readLdapExternalServerResponse(ctx, updateResponse.LdapExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PingOneHttpExternalServerResponse != nil {
			readPingOneHttpExternalServerResponse(ctx, updateResponse.PingOneHttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpExternalServerResponse != nil {
			readHttpExternalServerResponse(ctx, updateResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OracleUnifiedDirectoryExternalServerResponse != nil {
			readOracleUnifiedDirectoryExternalServerResponse(ctx, updateResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurExternalServerResponse != nil {
			readConjurExternalServerResponse(ctx, updateResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonAwsExternalServerResponse != nil {
			readAmazonAwsExternalServerResponse(ctx, updateResponse.AmazonAwsExternalServerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultExternalServerResponse != nil {
			readVaultExternalServerResponse(ctx, updateResponse.VaultExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *externalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state externalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerAPI.DeleteExternalServerExecute(r.apiClient.ExternalServerAPI.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the External Server", err, httpResp)
		return
	}
}

func (r *externalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExternalServer(ctx, req, resp)
}

func (r *defaultExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExternalServer(ctx, req, resp)
}

func importExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
