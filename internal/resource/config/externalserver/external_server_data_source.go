package externalserver

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
	_ datasource.DataSource              = &externalServerDataSource{}
	_ datasource.DataSourceWithConfigure = &externalServerDataSource{}
)

// Create a External Server data source
func NewExternalServerDataSource() datasource.DataSource {
	return &externalServerDataSource{}
}

// externalServerDataSource is the datasource implementation.
type externalServerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *externalServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_server"
}

// Configure adds the provider configured client to the data source.
func (r *externalServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type externalServerDataSourceModel struct {
	Id                                     types.String `tfsdk:"id"`
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

// GetSchema defines the schema for the datasource.
func (r *externalServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a External Server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of External Server resource. Options are ['smtp', 'nokia-ds', 'ping-identity-ds', 'active-directory', 'jdbc', 'syslog', 'ping-identity-proxy-server', 'http-proxy', 'nokia-proxy-server', 'opendj', 'ldap', 'ping-one-http', 'http', 'oracle-unified-directory', 'conjur', 'amazon-aws', 'vault']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the Vault server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://vault.example.com:8200/\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"vault_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Vault server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the AWS service. Supported in PingDirectory product version 9.2.0.0+.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the CyberArk Conjur server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://conjur.example.com:8443/\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"aws_access_key_id": schema.StringAttribute{
				Description: "The access key ID that will be used if authentication should use an access key. If this is provided, then an aws-secret-access-key must also be provided.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Description: "The secret access key that will be used if authentication should use an access key. If this is provided, then an aws-access-key-id must also be provided.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"aws_region_name": schema.StringAttribute{
				Description: "The name of the AWS region containing the resources that will be accessed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Conjur server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_account_name": schema.StringAttribute{
				Description: "The name of the account with which the desired secrets are associated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_file": schema.StringAttribute{
				Description: "The path to a file containing the information needed to trust the certificate presented by the Conjur servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "The PIN needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\", \"PKCS12\", or \"BCFKS\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_url": schema.StringAttribute{
				Description: "The base URL of the external server, optionally including port number, for example \"https://externalService:9031\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"hostname_verification_method": schema.StringAttribute{
				Description: "The mechanism for checking if the hostname in the PingOne ID Token Validator's base-url value matches the name(s) stored inside the X.509 certificate presented by PingOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jdbc_driver_type": schema.StringAttribute{
				Description: "Specifies a supported database driver type. The driver class will be automatically selected based on this selection. We highly recommend using a JDBC 4 driver that is suitable for the current Java platform.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jdbc_driver_url": schema.StringAttribute{
				Description: "Specify the complete JDBC URL which will be used instead of the automatic URL format. You must select type 'other' for the jdbc-driver-type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The certificate alias within the keystore to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"response_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to PingOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"basic_authentication_username": schema.StringAttribute{
				Description: "The username to use to authenticate to the HTTP Proxy External Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"basic_authentication_passphrase_provider": schema.StringAttribute{
				Description: "A passphrase provider that provides access to the password to use to authenticate to the HTTP Proxy External Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"transport_mechanism": schema.StringAttribute{
				Description: "The transport mechanism that should be used when communicating with the syslog server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"database_name": schema.StringAttribute{
				Description: "Specifies which database to connect to. This is ignored if jdbc-driver-url is specified.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"verify_credentials_method": schema.StringAttribute{
				Description: "The mechanism to use to verify user credentials while ensuring that the ability to process other operations is not impacted by an alternate authorization identity.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_administrative_operation_control": schema.BoolAttribute{
				Description: "Indicates whether to include the administrative operation request control in requests sent to this server which are intended for administrative operations (e.g., health checking) rather than requests directly from clients.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description: "The host name of the smtp server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "The port number where the smtp server listens for requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "Specifies the location for the LDAP External Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"validation_query": schema.StringAttribute{
				Description: "The SQL query that will be used to validate connections to the database before making them available to the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"validation_query_timeout": schema.StringAttribute{
				Description: "Specifies the amount of time to wait for a response from the database when executing the validation query, if one is set. If the timeout is exceeded, the Directory Server will drop the connection and obtain a new one. A value of zero indicates no timeout.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jdbc_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the JDBC datasource.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"transaction_isolation_level": schema.StringAttribute{
				Description: "This property specifies the default transaction isolation level for connections to this JDBC External Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "The DN to use to bind to the target LDAP server if simple authentication is required.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"smtp_security": schema.StringAttribute{
				Description: "This property specifies type of connection security to use when connecting to the outgoing mail server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_name": schema.StringAttribute{
				Description: "The name of the login account to use when connecting to the smtp server. Both username and password must be supplied if this attribute is set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_security": schema.StringAttribute{
				Description: "The mechanism to use to secure communication with the directory server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: "The mechanism to use to authenticate to the target server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"health_check_connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established for the purpose of performing a health check. If the connection cannot be established within this length of time, the server will be classified as unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_connection_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that connections to this server should be allowed to remain established before being closed and replaced with newly-established connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_expired_connection_disconnect_interval": schema.StringAttribute{
				Description: "Specifies the minimum length of time that should pass between connection closures as a result of the connections being established for longer than the maximum connection age. This may help avoid cases in which a large number of connections are closed and re-established in a short period of time because of the maximum connection age.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_response_size": schema.StringAttribute{
				Description: "Specifies the maximum response size that should be supported for messages received from the LDAP external server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "The key manager provider to use if SSL or StartTLS is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL or StartTLS is to be used for connection-level security.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "The number of connections to initially establish to the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "The maximum number of concurrent connections to maintain for the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"defunct_connection_result_code": schema.SetAttribute{
				Description: "Specifies the operation result code values that should cause the associated connection should be considered defunct. If an operation fails with one of these result codes, then it will be terminated and an attempt will be made to establish a new connection in its place.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"abandon_on_timeout": schema.BoolAttribute{
				Description: "Indicates whether to send an abandon request for an operation for which a response timeout is encountered. A request which has timed out on one server may be retried on another server regardless of whether an abandon request is sent, but if the initial attempt is not abandoned then a long-running operation may unnecessarily continue to consume processing resources on the initial server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The login password for the specified user name. Both username and password must be supplied if this attribute is set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the login password for the specified user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"smtp_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a connection or attempted connection to a SMTP server may take.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"smtp_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the smtp server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a SmtpExternalServerResponse object into the model struct
func readSmtpExternalServerResponseDataSource(ctx context.Context, r *client.SmtpExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.SmtpSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerSmtpSecurityProp(r.SmtpSecurity), false)
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.SmtpTimeout = internaltypes.StringTypeOrNil(r.SmtpTimeout, false)
	state.SmtpConnectionProperties = internaltypes.GetStringSet(r.SmtpConnectionProperties)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a NokiaDsExternalServerResponse object into the model struct
func readNokiaDsExternalServerResponseDataSource(ctx context.Context, r *client.NokiaDsExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nokia-ds")
	state.Id = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a PingIdentityDsExternalServerResponse object into the model struct
func readPingIdentityDsExternalServerResponseDataSource(ctx context.Context, r *client.PingIdentityDsExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-identity-ds")
	state.Id = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ActiveDirectoryExternalServerResponse object into the model struct
func readActiveDirectoryExternalServerResponseDataSource(ctx context.Context, r *client.ActiveDirectoryExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("active-directory")
	state.Id = types.StringValue(r.Id)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a JdbcExternalServerResponse object into the model struct
func readJdbcExternalServerResponseDataSource(ctx context.Context, r *client.JdbcExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc")
	state.Id = types.StringValue(r.Id)
	state.JdbcDriverType = types.StringValue(r.JdbcDriverType.String())
	state.JdbcDriverURL = internaltypes.StringTypeOrNil(r.JdbcDriverURL, false)
	state.DatabaseName = internaltypes.StringTypeOrNil(r.DatabaseName, false)
	state.ServerHostName = internaltypes.StringTypeOrNil(r.ServerHostName, false)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ValidationQuery = internaltypes.StringTypeOrNil(r.ValidationQuery, false)
	state.ValidationQueryTimeout = internaltypes.StringTypeOrNil(r.ValidationQueryTimeout, false)
	state.JdbcConnectionProperties = internaltypes.GetStringSet(r.JdbcConnectionProperties)
	state.TransactionIsolationLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerTransactionIsolationLevelProp(r.TransactionIsolationLevel), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SyslogExternalServerResponse object into the model struct
func readSyslogExternalServerResponseDataSource(ctx context.Context, r *client.SyslogExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.TransportMechanism = types.StringValue(r.TransportMechanism.String())
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.TrustManagerProvider = types.StringValue(r.TrustManagerProvider)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a PingIdentityProxyServerExternalServerResponse object into the model struct
func readPingIdentityProxyServerExternalServerResponseDataSource(ctx context.Context, r *client.PingIdentityProxyServerExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-identity-proxy-server")
	state.Id = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a HttpProxyExternalServerResponse object into the model struct
func readHttpProxyExternalServerResponseDataSource(ctx context.Context, r *client.HttpProxyExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http-proxy")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.BasicAuthenticationUsername = internaltypes.StringTypeOrNil(r.BasicAuthenticationUsername, false)
	state.BasicAuthenticationPassphraseProvider = internaltypes.StringTypeOrNil(r.BasicAuthenticationPassphraseProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a NokiaProxyServerExternalServerResponse object into the model struct
func readNokiaProxyServerExternalServerResponseDataSource(ctx context.Context, r *client.NokiaProxyServerExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nokia-proxy-server")
	state.Id = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a OpendjExternalServerResponse object into the model struct
func readOpendjExternalServerResponseDataSource(ctx context.Context, r *client.OpendjExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("opendj")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a LdapExternalServerResponse object into the model struct
func readLdapExternalServerResponseDataSource(ctx context.Context, r *client.LdapExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a PingOneHttpExternalServerResponse object into the model struct
func readPingOneHttpExternalServerResponseDataSource(ctx context.Context, r *client.PingOneHttpExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one-http")
	state.Id = types.StringValue(r.Id)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerPingOneHttpHostnameVerificationMethodProp(r.HostnameVerificationMethod), false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.ConnectTimeout = internaltypes.StringTypeOrNil(r.ConnectTimeout, false)
	state.ResponseTimeout = internaltypes.StringTypeOrNil(r.ResponseTimeout, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a HttpExternalServerResponse object into the model struct
func readHttpExternalServerResponseDataSource(ctx context.Context, r *client.HttpExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	state.BaseURL = types.StringValue(r.BaseURL)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerHttpHostnameVerificationMethodProp(r.HostnameVerificationMethod), false)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, false)
	state.ConnectTimeout = internaltypes.StringTypeOrNil(r.ConnectTimeout, false)
	state.ResponseTimeout = internaltypes.StringTypeOrNil(r.ResponseTimeout, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a OracleUnifiedDirectoryExternalServerResponse object into the model struct
func readOracleUnifiedDirectoryExternalServerResponseDataSource(ctx context.Context, r *client.OracleUnifiedDirectoryExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("oracle-unified-directory")
	state.Id = types.StringValue(r.Id)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, false)
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, false)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, false)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ConjurExternalServerResponse object into the model struct
func readConjurExternalServerResponseDataSource(ctx context.Context, r *client.ConjurExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.ConjurServerBaseURI = internaltypes.GetStringSet(r.ConjurServerBaseURI)
	state.ConjurAuthenticationMethod = types.StringValue(r.ConjurAuthenticationMethod)
	state.ConjurAccountName = types.StringValue(r.ConjurAccountName)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, false)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AmazonAwsExternalServerResponse object into the model struct
func readAmazonAwsExternalServerResponseDataSource(ctx context.Context, r *client.AmazonAwsExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-aws")
	state.Id = types.StringValue(r.Id)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.AuthenticationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerAmazonAwsAuthenticationMethodProp(r.AuthenticationMethod), false)
	state.AwsAccessKeyID = internaltypes.StringTypeOrNil(r.AwsAccessKeyID, false)
	state.AwsRegionName = types.StringValue(r.AwsRegionName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a VaultExternalServerResponse object into the model struct
func readVaultExternalServerResponseDataSource(ctx context.Context, r *client.VaultExternalServerResponse, state *externalServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.VaultServerBaseURI = internaltypes.GetStringSet(r.VaultServerBaseURI)
	state.VaultAuthenticationMethod = types.StringValue(r.VaultAuthenticationMethod)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, false)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *externalServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state externalServerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SmtpExternalServerResponse != nil {
		readSmtpExternalServerResponseDataSource(ctx, readResponse.SmtpExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaDsExternalServerResponse != nil {
		readNokiaDsExternalServerResponseDataSource(ctx, readResponse.NokiaDsExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityDsExternalServerResponse != nil {
		readPingIdentityDsExternalServerResponseDataSource(ctx, readResponse.PingIdentityDsExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveDirectoryExternalServerResponse != nil {
		readActiveDirectoryExternalServerResponseDataSource(ctx, readResponse.ActiveDirectoryExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcExternalServerResponse != nil {
		readJdbcExternalServerResponseDataSource(ctx, readResponse.JdbcExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogExternalServerResponse != nil {
		readSyslogExternalServerResponseDataSource(ctx, readResponse.SyslogExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PingIdentityProxyServerExternalServerResponse != nil {
		readPingIdentityProxyServerExternalServerResponseDataSource(ctx, readResponse.PingIdentityProxyServerExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HttpProxyExternalServerResponse != nil {
		readHttpProxyExternalServerResponseDataSource(ctx, readResponse.HttpProxyExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NokiaProxyServerExternalServerResponse != nil {
		readNokiaProxyServerExternalServerResponseDataSource(ctx, readResponse.NokiaProxyServerExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OpendjExternalServerResponse != nil {
		readOpendjExternalServerResponseDataSource(ctx, readResponse.OpendjExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapExternalServerResponse != nil {
		readLdapExternalServerResponseDataSource(ctx, readResponse.LdapExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PingOneHttpExternalServerResponse != nil {
		readPingOneHttpExternalServerResponseDataSource(ctx, readResponse.PingOneHttpExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HttpExternalServerResponse != nil {
		readHttpExternalServerResponseDataSource(ctx, readResponse.HttpExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OracleUnifiedDirectoryExternalServerResponse != nil {
		readOracleUnifiedDirectoryExternalServerResponseDataSource(ctx, readResponse.OracleUnifiedDirectoryExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurExternalServerResponse != nil {
		readConjurExternalServerResponseDataSource(ctx, readResponse.ConjurExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonAwsExternalServerResponse != nil {
		readAmazonAwsExternalServerResponseDataSource(ctx, readResponse.AmazonAwsExternalServerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VaultExternalServerResponse != nil {
		readVaultExternalServerResponseDataSource(ctx, readResponse.VaultExternalServerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
